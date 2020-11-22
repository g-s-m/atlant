package dto

import (
	"atlant/errors"
	pb "atlant/generated/interface"
	"time"
)

type SortType uint32

const (
	ByProduct SortType = iota
	ByPrice
	ByPriceChange
	ByLastChange
)

var SortTypeTranslate = map[pb.Sort_Type]SortType{
	pb.Sort_BY_PRODUCT:      ByProduct,
	pb.Sort_BY_PRICE:        ByPrice,
	pb.Sort_BY_PRICE_CHANGE: ByPriceChange,
	pb.Sort_BY_LAST_CHANGE:  ByLastChange,
}

var ErrorToStatusTranslate = map[errors.ErrorType]pb.FetchReply_Status{
	errors.Ok:                  pb.FetchReply_OK,
	errors.WrongFile:           pb.FetchReply_WRONG_FILE_FORMAT,
	errors.ResourceUnavailable: pb.FetchReply_RESOURCE_UNAVAILABLE,
	errors.InternalError:       pb.FetchReply_INTERNAL_ERROR,
}

type Product struct {
	Name        string    `bson:"product"`
	Price       float64   `bson:"price"`
	ChangeCount uint64    `bson:"change_count"`
	ChangeDate  time.Time `bson:"change_date"`
}

type Page struct {
	Start uint64
	Size  int64
}

type SortParams struct {
	Type   SortType
	SortUp bool
}

func SortDto(s *pb.Sort) SortParams {
	return SortParams{
		Type:   SortTypeTranslate[s.GetSortType()],
		SortUp: s.GetSortUp(),
	}
}

func PageDto(page *pb.Page) Page {
	return Page{
		Start: page.GetStart(),
		Size:  page.GetSize(),
	}
}

func ErrorToStatus(err error) pb.FetchReply_Status {
	if err != nil {
		if v, ok := err.(errors.ServiceError); ok {
			return ErrorToStatusTranslate[v.Type()]
		}
		return pb.FetchReply_INTERNAL_ERROR
	}
	return pb.FetchReply_OK
}

func ProductDto(p Product) pb.Product {
	return pb.Product{
		Product:   p.Name,
		Price:     p.Price,
		Timestamp: p.ChangeDate.Unix(),
		Changed:   p.ChangeCount,
	}
}
