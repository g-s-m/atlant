package service

import (
	aerr "atlant/errors"
	srv "atlant/generated/interface"
	dto "atlant/service/dto"
	"context"
	"errors"
	"testing"
	"time"
)

type HandlerMock struct {
}

func (h HandlerMock) DoFetch(file string) error {
	if file != "ok" {
		if file == "unavailable" {
			return aerr.NewServiceError(aerr.ResourceUnavailable)
		}
		if file == "wrong" {
			return aerr.NewServiceError(aerr.WrongFile)
		}
		if file == "int" {
			return aerr.NewServiceError(aerr.InternalError)
		}
		if file == "custom" {
			return errors.New("Fake Error")
		}
	}
	return nil
}

func (h HandlerMock) DoList(page dto.Page, sort dto.SortParams) ([]*dto.Product, error) {
	return []*dto.Product{
		&dto.Product{
			Name:        "product",
			Price:       1.234,
			ChangeCount: 0,
			ChangeDate:  time.Date(2020, 11, 1, 12, 0, 0, 0, time.UTC),
		},
		&dto.Product{
			Name:        "product2",
			Price:       5.678,
			ChangeCount: 1,
			ChangeDate:  time.Date(2021, 1, 1, 13, 0, 0, 0, time.UTC),
		},
	}, nil
}

type HandlerLongOperationMock struct {
}

func (h HandlerLongOperationMock) DoFetch(file string) error {
	time.Sleep(10 * time.Second)
	return nil
}

func (h HandlerLongOperationMock) DoList(page dto.Page, sort dto.SortParams) ([]*dto.Product, error) {
	time.Sleep(10 * time.Second)
	return []*dto.Product{}, nil
}

func TestServiceFetchOk(t *testing.T) {
	s := NewServer(HandlerMock{})
	req := &srv.FetchRequest{
		Url: "ok",
	}
	r, err := s.Fetch(context.Background(), req)
	if err != nil {
		t.Errorf("Not nil error: %v", err)
	}
	if r.Status != srv.FetchReply_OK {
		t.Errorf("Status error. Expected(%v), Actual(%v)", srv.FetchReply_OK, r.Status)
	}
}

func TestServiceFetchErrors(t *testing.T) {
	s := NewServer(HandlerMock{})
	errMap := map[string]srv.FetchReply_Status{
		"unavailable": srv.FetchReply_RESOURCE_UNAVAILABLE,
		"wrong":       srv.FetchReply_WRONG_FILE_FORMAT,
		"int":         srv.FetchReply_INTERNAL_ERROR,
		"custom":      srv.FetchReply_INTERNAL_ERROR,
	}
	for k, v := range errMap {
		req := &srv.FetchRequest{
			Url: k,
		}
		r, err := s.Fetch(context.Background(), req)
		if err != nil {
			t.Errorf("Not nil error: %v", err)
		}
		if r.Status != v {
			t.Errorf("Status error. Expected(%v), Actual(%v)", v, r.Status)
		}
	}
}

func TestServiceFetchTimeout(t *testing.T) {
	s := NewServer(HandlerLongOperationMock{})
	req := &srv.FetchRequest{
		Url: "ok",
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	r, err := s.Fetch(ctx, req)
	if err != nil {
		t.Errorf("Not nil error: %v", err)
	}
	if r.Status != srv.FetchReply_RESOURCE_UNAVAILABLE {
		t.Errorf("Status error. Expected(%v), Actual(%v)", srv.FetchReply_RESOURCE_UNAVAILABLE, r.Status)
	}
}

func TestServiceListOk(t *testing.T) {
	s := NewServer(HandlerMock{})
	req := &srv.ListRequest{
		Page: &srv.Page{
			Size:  1,
			Start: 0,
		},
		Sort: &srv.Sort{
			SortType: srv.Sort_BY_PRODUCT,
			SortUp:   true,
		},
	}
	r, err := s.List(context.Background(), req)
	if err != nil {
		t.Errorf("Not nil error: %v", err)
	}
	if len(r.ProductList) != 2 {
		t.Errorf("Wrong response list length. Expected(%v), Actual(%v)", 0, len(r.ProductList))
	}
	if r.ProductList[0].GetProduct() != "product" {
		t.Errorf("Wrong Product. Expected(%v), Actual(%v)", "product", r.ProductList[0].GetProduct())
	}
	if r.ProductList[0].GetPrice() != 1.234 {
		t.Errorf("Wrong Price. Expected(%v), Actual(%v)", 1.234, r.ProductList[0].GetPrice())
	}
	if r.ProductList[0].GetChanged() != 0 {
		t.Errorf("Wrong Change Count. Expected(%v), Actual(%v)", 0, r.ProductList[0].GetChanged())
	}
	if r.ProductList[0].GetTimestamp() != time.Date(2020, 11, 1, 12, 0, 0, 0, time.UTC).Unix() {
		t.Errorf("Wrong Timestamp. Expected(%v), Actual(%v)",
			time.Date(2020, 11, 1, 12, 0, 0, 0, time.UTC).Unix(), r.ProductList[0].GetTimestamp())
	}

	if r.ProductList[1].GetProduct() != "product2" {
		t.Errorf("Wrong Product. Expected(%v), Actual(%v)", "product2", r.ProductList[1].GetProduct())
	}
	if r.ProductList[1].GetPrice() != 5.678 {
		t.Errorf("Wrong Price. Expected(%v), Actual(%v)", 5.678, r.ProductList[1].GetPrice())
	}
	if r.ProductList[1].GetChanged() != 1 {
		t.Errorf("Wrong Change Count. Expected(%v), Actual(%v)", 1, r.ProductList[1].GetChanged())
	}
	if r.ProductList[1].GetTimestamp() != time.Date(2021, 1, 1, 13, 0, 0, 0, time.UTC).Unix() {
		t.Errorf("Wrong Timestamp. Expected(%v), Actual(%v)",
			time.Date(2021, 1, 1, 13, 0, 0, 0, time.UTC).Unix(), r.ProductList[1].GetTimestamp())
	}
}

func TestServiceListTimeout(t *testing.T) {
	s := NewServer(HandlerLongOperationMock{})
	req := &srv.ListRequest{
		Page: &srv.Page{
			Size:  1,
			Start: 0,
		},
		Sort: &srv.Sort{
			SortType: srv.Sort_BY_PRODUCT,
			SortUp:   true,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	r, err := s.List(ctx, req)
	if err != nil {
		t.Errorf("Not nil error: %v", err)
	}
	if len(r.ProductList) != 0 {
		t.Errorf("Wrong response list length. Expected(%v), Actual(%v)", 0, len(r.ProductList))
	}
}
