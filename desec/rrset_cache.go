package desec

import (
	"context"
	"sync"

	dsc "github.com/nrdcg/desec"
)

type DesecCache struct {
	mutex sync.Mutex
	data  map[string]map[string]dsc.RRSet
}

func NewDesecCache() DesecCache {
	return DesecCache{sync.Mutex{}, nil}
}

func (r *DesecCache) GetRRSetById(ctx context.Context, c *dsc.Client, id string) (*dsc.RRSet, error) {
	domainName, _, _, err := namesFromId(id)
	if err != nil {
		return nil, err
	}

	r.mutex.Lock()

	defer func() {
		r.mutex.Unlock()
	}()

	if r.data == nil {
		r.data = make(map[string]map[string]dsc.RRSet)
	}
	if r.data[domainName] == nil {
		recs, err := c.Records.GetAll(ctx, domainName, nil)
		if err != nil {
			if isNotFoundError(err) {
				return nil, nil
			}
			return nil, err
		}

		d := make(map[string]dsc.RRSet)
		for _, rec := range recs {
			id := idFromNames(rec.Domain, rec.SubName, rec.Type)
			d[id] = rec
		}
		r.data[domainName] = d
	}

	result, ok := r.data[domainName][id]
	if ok {
		return &result, nil
	} else {
		return nil, nil
	}
}

func (r *DesecCache) Clear() {
	r.mutex.Lock()
	r.data = nil
	r.mutex.Unlock()
}
