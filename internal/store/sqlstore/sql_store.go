package sqlstore

import (
	"crypto/md5"
	"errors"
	"fmt"
	_ "github.com/lib/pq" //nolint:goimports
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"restapi_langparser/internal/model"
	"restapi_langparser/internal/store"
	"strings"
	"sync"
	"time"
)

const (
	BadDomain = `select d.* from domains d
					left join queues q on d.id=q.domain_id
					left join requests r on d.id=r.domain_id
					where not q.domain_id isnull
					and d.response_code != ''
					and q.update_at<now()
					and q.deleted_at isnull
					order by q.update_at
					limit 1`
	UserRequest = `select d.* from domains d
					left join queues q on d.id=q.domain_id
					left join requests r on d.id=r.domain_id
					where not q.domain_id isnull
					and not r.code isnull
					and d.response_code=''
					and q.deleted_at isnull
					order by q.update_at
					limit 1`
	Queue = `select d.* from domains d
					left join queues q on d.id=q.domain_id
					left join requests r on d.id=r.domain_id
					where not q.domain_id isnull
					and d.response_code=''
					and q.update_at<now()
					and r.code isnull
					and q.deleted_at isnull
					order by q.update_at
					limit 1`
)

type Store struct {
	db               *gorm.DB
	ProxyRepository  store.IProxyRepository
	DomainRepository store.IDomainRepository
	callbacks        map[string]string
	m                sync.Mutex
}

func New(db *gorm.DB) store.IStore {
	return &Store{
		db:               db,
		DomainRepository: NewDomainRepository(db),
		ProxyRepository:  NewProxyRepository(db),
		callbacks:        make(map[string]string),
	}
}

func (s *Store) Migrate() error {
	m := s.db.Migrator()
	return m.AutoMigrate(&model.Proxy{}, &model.Domain{}, &model.Request{}, &model.Queue{})
}

func (s *Store) Proxy() store.IProxyRepository {
	return s.ProxyRepository
}

func (s *Store) Domain() store.IDomainRepository {
	return s.DomainRepository
}

// AddDomains create new domains and add it to queue
func (s *Store) AddDomains(list *[]model.Domain) error {
	s.m.Lock()
	defer s.m.Unlock()

	err := s.db.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "host"}},
			//DoUpdates: clause.AssignmentColumns([]string{"host"}),
			DoNothing: true,
		}).Clauses(clause.Returning{}).Create(&list).Error
	if err != nil {
		return err
	}

	// add to queue
	queue := make([]model.Queue, 0)
	for _, td := range *list { // exclude processed domains
		if td.ResponseCode == model.ResponseOk {
			continue
		}
		queue = append(queue, model.Queue{
			DomainID: td.ID,
			UpdateAt: time.Now(),
		})
	}
	if len(queue) == 0 {
		return nil
	}
	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "domain_id"}},
		UpdateAll: true,
	}).Create(&queue).Error
}

func (s *Store) GetDomains(hosts []string) ([]model.Domain, error) {
	var domains []model.Domain
	err := s.db.Joins("left join queues on queues.domain_id=domains.id").Where("host in ?", hosts).Where("queues.domain_id isnull").Find(&domains).Error
	if err != nil {
		return nil, err
	}
	if len(domains) != len(hosts) {
		return nil, nil
	}
	return domains, nil
}

// SaveDomain
func (s *Store) SaveDomain(domain model.Domain) error {
	return s.db.Save(&domain).Error
}

func (s *Store) GetCompletedRequests(domain model.Domain) ([]string, error) {
	var codes []string
	err := s.db.Raw(`select r.code from requests r
		left join queues q on r.domain_id=q.domain_id
		left join domains d on d.id=r.domain_id
		left join (select code from requests where domain_id=?) c on c.code=r.code
		group by r.code
		having count(q.domain_id)=0`, domain.ID).Scan(&codes).Error
	return codes, err
}

func (s *Store) GetCallbacks(codes []string) map[string]string {
	res := make(map[string]string)
	for _, code := range codes {
		if callback, exist := s.callbacks[code]; exist {
			res[code] = callback
		}
	}
	return res
}

//todo remove it
func (s *Store) Test() {
	s.callbacks["fddc1482c9d29f47961855ec9bacc9aa"] = "localhost:8080/echo"
	testDomain := model.Domain{}
	testDomain.ID = 1
	callbacks, err := s.GetCompletedRequests(testDomain)
	if err != nil {
		logrus.Error(err)
		return
	}
	logrus.Infof("found callbacks: %v", callbacks)
}

func (s *Store) GetFromQueue(priority string) *model.Domain {
	var domain model.Domain

	err := s.db.Raw(priority).Scan(&domain).Error
	if err != nil || domain.ID == 0 {
		return nil
	}

	err = s.db.Where("domain_id=?", domain.ID).Delete(&model.Queue{}).Error
	if err != nil {
		return nil
	}

	return &domain
}

func (s *Store) AddToQueue(updateAt time.Time, list ...model.Domain) error {
	q := make([]model.Queue, len(list))
	for i, domain := range list {
		q[i].DomainID = domain.ID
		q[i].UpdateAt = updateAt
	}
	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "domain_id"}},
		DoNothing: true,
	}).Create(&q).Error
}

func (s *Store) ReturnToQueue(updateAt time.Time, domain model.Domain) error {
	return s.db.
		Unscoped().
		Model(&model.Queue{}).
		Where("domain_id", domain.ID).
		Updates(map[string]interface{}{"deleted_at": nil, "update_at": updateAt}).
		Error
}

func (s *Store) RemoveFromQueue(domain model.Domain) error {
	return s.db.
		Unscoped().
		Where("domain_id=?", domain.ID).
		Delete(&model.Queue{}).
		Error
}

func createRequestCode(list []model.Domain) string {
	sb := strings.Builder{}
	for _, domain := range list {
		sb.WriteString(domain.Host)
	}
	return fmt.Sprintf("%x", md5.Sum([]byte(sb.String())))
}

func (s *Store) CreateRequest(list []model.Domain, callback *string) (requestCode string, err error) {
	requestCode = createRequestCode(list)
	if callback != nil {
		s.callbacks[requestCode] = *callback
	}
	request := make([]model.Request, len(list))
	for i, domain := range list {
		request[i] = model.Request{
			DomainID: domain.ID,
			Code:     requestCode,
		}
	}

	return requestCode, s.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "domain_id"},
			{Name: "code"},
		},
		DoNothing: true,
	}).Create(&request).Error
}

func (s *Store) GetRequest(requestCode string) ([]model.Domain, error) {
	var cnt int64
	s.db.Model(&model.Request{}).Joins("left join queues on queues.domain_id=requests.domain_id where not queues.domain_id isnull and requests.code=?", requestCode).Count(&cnt)
	if cnt > 0 {
		return nil, errors.New("not ready")
	}

	var domains []model.Domain
	s.db.Joins("left join requests on domains.id=requests.domain_id").Where("requests.code=?", requestCode).Find(&domains)

	return domains, nil
}
