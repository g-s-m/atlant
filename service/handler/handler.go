package handler

import (
	"atlant/service"
)

struct RequestHandler {
}

func (p *RequestHandler) DoFetch(path string) error {
	return nil
}

func (p *RequestHandler) DoList() ([]Product, error) {
	return Product{}, nil
}
