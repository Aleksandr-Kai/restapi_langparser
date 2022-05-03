package langfinder

import "time"

const lifeTime = time.Hour * 72

type requestResult struct {
	total      int
	position   int
	DomainIDs  []int64
	expiration time.Time
}

func NewResult(total int) *requestResult {
	return &requestResult{
		total:      total,
		DomainIDs:  make([]int64, total),
		expiration: time.Now().Add(lifeTime),
	}
}

func (r *requestResult) Ready() bool {
	return r.position == r.total
}

func (r *requestResult) Add(id int64) {
	if r.position == r.total {
		return
	}
	r.DomainIDs[r.position] = id
	r.position++
}
