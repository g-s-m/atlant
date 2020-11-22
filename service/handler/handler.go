package handler

import (
	"atlant/data"
	aerr "atlant/errors"
	"atlant/requestor"
	"atlant/service/dto"
	"encoding/csv"
	"log"
	"strconv"
	"strings"
)

type RequestHandler struct {
	repo data.IProductsActions
}

func (p RequestHandler) DoFetch(path string) error {
	file, err := requestor.GetCsvFile(path, 30)
	if err != nil {
		return err
	}
	r := csv.NewReader(strings.NewReader(string(file)))
	r.Comma = ';'
	records, err := r.ReadAll()
	if err != nil {
		log.Printf("Error in csv file: %v", err)
		return aerr.NewServiceError(aerr.WrongFile)
	}

	for _, line := range records {
		if len(line) != 2 {
			log.Printf("Csv file contains more than 2 columns")
			return aerr.NewServiceError(aerr.WrongFile)
		}
		price, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			log.Printf("Wrong format of price in csv file")
			return aerr.NewServiceError(aerr.WrongFile)
		}
		p.repo.Save(line[0], price)
	}

	return nil
}

func (p RequestHandler) DoList(page dto.Page, sort dto.SortParams) ([]*dto.Product, error) {
	type SortFunc func(uint64, int64, bool) ([]*dto.Product, error)
	var s = map[dto.SortType]SortFunc{
		dto.ByProduct: func(start uint64, leng int64, upSort bool) ([]*dto.Product, error) {
			return p.repo.LoadByProduct(start, leng, upSort)
		},
		dto.ByPrice: func(start uint64, leng int64, upSort bool) ([]*dto.Product, error) {
			return p.repo.LoadByPrice(start, leng, upSort)
		},
		dto.ByPriceChange: func(start uint64, leng int64, upSort bool) ([]*dto.Product, error) {
			return p.repo.LoadByChangeCount(start, leng, upSort)
		},
		dto.ByLastChange: func(start uint64, leng int64, upSort bool) ([]*dto.Product, error) {
			return p.repo.LoadByDate(start, leng, upSort)
		},
	}
	return s[sort.Type](page.Start, page.Size, sort.SortUp)
}

func NewRequestHandler(actions data.IProductsActions) RequestHandler {
	return RequestHandler{
		repo: actions,
	}
}
