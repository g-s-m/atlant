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

func (h HandlerMock) DoList(page dto.Page, sort dto.SortParams) ([]dto.Product, error) {
	return []dto.Product{}, nil
}

type HandlerLongOperationMock struct {
}

func (h HandlerLongOperationMock) DoFetch(file string) error {
	time.Sleep(10 * time.Second)
	return nil
}

func (h HandlerLongOperationMock) DoList(page dto.Page, sort dto.SortParams) ([]dto.Product, error) {
	time.Sleep(10 * time.Second)
	return []dto.Product{}, nil
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
		Url: "k",
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

func TestServiceDoListTimeout(t *testing.T) {
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
