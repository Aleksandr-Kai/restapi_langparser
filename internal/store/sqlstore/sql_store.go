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
)

type Store struct {
	db               *gorm.DB
	ProxyRepository  store.IProxyRepository
	DomainRepository store.IDomainRepository
	callbacks        map[string]string
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

func (s *Store) AddDomains(list *[]model.Domain) error {
	err := s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "host"}},
		DoUpdates: clause.AssignmentColumns([]string{"host"}),
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

func (s *Store) SaveDomain(domain model.Domain) error {
	if err := s.db.Save(&domain).Error; err != nil {
		return err
	}

	if domain.ResponseCode == model.ResponseOk { // remove from queue if no errors
		if err := s.RemoveFromQueue(domain); err != nil {
			return fmt.Errorf("failed to delete from queue: %w", err)
		}
	} else if err := s.ReturnToQueue(domain); err != nil { // return to queue if errors
		return fmt.Errorf("failed to add to queue: %w", err)
	}
	return nil
}

func (s *Store) GetFromQueue() *model.Domain {
	var domain model.Domain

	err := s.db.Transaction(func(tx *gorm.DB) error {
		var task model.Queue

		// get USER request queue
		err := tx.Joins("left join requests on requests.domain_id=queues.domain_id").Find(&task).Error
		if err != nil {
			return err
		}

		if task.DomainID == 0 {
			// get LIST queue
			err = tx.Joins("left join requests on requests.domain_id=queues.domain_id where requests.domain_id is null order by queues.domain_id").Find(&task).Error
			if err != nil {
				return err
			}
		}

		if task.DomainID == 0 {
			return errors.New("queue empty")
		}

		err = tx.Where("domain_id=?", task.DomainID).Delete(&task).Error
		if err != nil {
			return err
		}

		err = tx.Find(&domain, task.DomainID).Error
		if err != nil {
			return err
		}

		if domain.ID == 0 {
			return errors.New("queue empty")
		}
		return nil
	})
	if err != nil {
		logrus.Warning(err)
		return nil
	}

	return &domain
}

func (s *Store) AddToQueue(list ...model.Domain) error {
	q := make([]model.Queue, len(list))
	for i, domain := range list {
		q[i].DomainID = domain.ID
	}
	return s.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "domain_id"}},
		DoNothing: true,
	}).Create(&q).Error
}

func (s *Store) ReturnToQueue(domain model.Domain) error {
	return s.db.Unscoped().Model(&model.Queue{}).Where("domain_id", domain.ID).Update("deleted_at", nil).Error
}

func (s *Store) RemoveFromQueue(domain model.Domain) error {
	return s.db.Unscoped().Delete(&model.Queue{
		DomainID: domain.ID,
	}).Error
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
