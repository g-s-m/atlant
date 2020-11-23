package data

import (
	"atlant/service/dto"
)

type IProductsActions interface {
	Save(product string, price float64) error

	LoadByProduct(start uint64, leng int64, upSort bool) ([]*dto.Product, error)
	LoadByPrice(start uint64, leng int64, upSort bool) ([]*dto.Product, error)
	LoadByChangeCount(start uint64, leng int64, upSort bool) ([]*dto.Product, error)
	LoadByDate(start uint64, leng int64, upSort bool) ([]*dto.Product, error)
}
