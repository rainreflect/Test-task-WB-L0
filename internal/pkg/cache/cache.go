package cache

import (
	"log"
	"rec/internal/pkg/postgresql/model"
)

type CacheStore map[string]model.OrderData

type CacheServ struct {
	cS CacheStore
}

func NewCacheServ() *CacheServ {
	CStore := make(CacheStore)
	CacheService := CacheServ{
		cS: CStore,
	}
	return &CacheService
}

func (CS *CacheServ) WriteToCache(data model.OrderData) {

	CS.cS[data.OrderUid] = data
	log.Println("Data in cache stored: ", data)

}

func (CS *CacheServ) ReadCache(orderUID string) model.OrderData {

	return CS.cS[orderUID]
}
