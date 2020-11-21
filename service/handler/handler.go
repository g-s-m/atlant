package handler

import "atlant/service/dto"

type RequestHandler struct {
}

func (p RequestHandler) DoFetch(path string) error {
	return nil
}

func (p RequestHandler) DoList(page dto.Page, sort dto.SortParams) ([]dto.Product, error) {
	return []dto.Product{}, nil
}
