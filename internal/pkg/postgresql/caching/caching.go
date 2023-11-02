package caching

import (
	"log"
	"rec/internal/pkg/cache"
	"rec/internal/pkg/postgresql"
	"rec/internal/pkg/postgresql/model"

	"github.com/go-playground/validator"
)

type CacheService struct {
	cache cache.CacheServ
	db    *postgresql.DBService
}

func NewCacheService(cache cache.CacheServ, db *postgresql.DBService) *CacheService {
	cS := CacheService{cache: cache, db: db}
	return &cS
}

func (cS *CacheService) Caching(data []byte) error {
	order := new(model.OrderData)
	item := new(model.OrderItem)

	err := order.Scan(data)
	if err != nil {
		log.Println("can't scan file")
		return err
	}
	validate := validator.New()
	err = validate.Struct(order)
	if err != nil {
		log.Println(err)
		return err
	}
	item.Data = *order
	item.ID = order.OrderUid

	cS.cache.WriteToCache(*order)
	return err
}

func (cS *CacheService) UncachingById(id string) model.OrderData {
	return cS.cache.ReadCache(id)
}

func (cS *CacheService) UncachingAll() ([]model.OrderItem, error) {
	orders, err := cS.db.Orders()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return orders, err
}

func (cS *CacheService) RestoreCache() error {
	items, err := cS.UncachingAll()
	if err == nil {
		if items == nil {
			log.Println("empty caching")
			return err
		} else {

			for _, item := range items {
				cS.cache.WriteToCache(item.Data)
			}
			log.Println("Cache restored")
		}
	} else {
		return err
	}
	return err
}
