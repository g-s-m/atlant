package dto

import (
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

type Product struct {
	Name        string
	Price       float64
	ChangeCount uint64
	ChangeDate  time.Time
}

func ErrorToStatus(err error) pb.FetchReply_Status {
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
