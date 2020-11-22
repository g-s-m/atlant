package handler

import (
	"atlant/service/dto"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ActionsMock struct {
}

func (p ActionsMock) Save(product string, price float64) error {
	return nil
}

func (p ActionsMock) LoadByProduct(start uint64, leng int64, upSort bool) ([]*dto.Product, error) {
	return []*dto.Product{
		&dto.Product{
			Name: "ByProduct",
		},
	}, nil
}

func (p ActionsMock) LoadByPrice(start uint64, leng int64, upSort bool) ([]*dto.Product, error) {
	return []*dto.Product{
		&dto.Product{
			Name: "ByPrice",
		},
	}, nil
}

func (p ActionsMock) LoadByChangeCount(start uint64, leng int64, upSort bool) ([]*dto.Product, error) {
	return []*dto.Product{
		&dto.Product{
			Name: "ByPriceChange",
		},
	}, nil
}
func (p ActionsMock) LoadByDate(start uint64, leng int64, upSort bool) ([]*dto.Product, error) {
	return []*dto.Product{
		&dto.Product{
			Name: "ByLastChange",
		},
	}, nil
}

func TestHandlerOk(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintln(w, `product;1.00
product1;12.23
super_product2;10.02
OneMoreProduct;0.05`)
	}))
	defer ts.Close()

	handler := NewRequestHandler(ActionsMock{})
	err := handler.DoFetch(ts.URL)
	if err != nil {
		t.Errorf("DoFetch returns a error %v", err)
	}
}

func TestHandlerWrongFormat1(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintln(w, `product;1.00
product1;qwe
super_product2;10.02
OneMoreProduct;0.05`)
	}))
	defer ts.Close()

	handler := NewRequestHandler(ActionsMock{})
	err := handler.DoFetch(ts.URL)
	if err == nil {
		t.Errorf("DoFetch must return a error")
	}
}

func TestHandlerWrongFormat2(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintln(w, `product;1.00
product1;1.23
super_product2;
10.02;
OneMoreProduct;0.05`)
	}))
	defer ts.Close()

	handler := NewRequestHandler(ActionsMock{})
	err := handler.DoFetch(ts.URL)
	if err == nil {
		t.Errorf("DoFetch rmust return a error")
	}
}

func TestHandlerDoListByName(t *testing.T) {
	handler := NewRequestHandler(ActionsMock{})
	result, err := handler.DoList(
		dto.Page{
			Size:  1,
			Start: 0,
		},
		dto.SortParams{
			Type:   dto.ByProduct,
			SortUp: true,
		})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if result[0].Name != "ByProduct" {
		t.Errorf("SortType ByProduct is expected, actual: %v", result[0].Name)
	}
}

func TestHandlerDoListByPrice(t *testing.T) {
	handler := NewRequestHandler(ActionsMock{})
	result, err := handler.DoList(
		dto.Page{
			Size:  1,
			Start: 0,
		},
		dto.SortParams{
			Type:   dto.ByPrice,
			SortUp: true,
		})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if result[0].Name != "ByPrice" {
		t.Errorf("SortType ByPrice is expected, actual: %v", result[0].Name)
	}
}

func TestHandlerDoListByChangeCount(t *testing.T) {
	handler := NewRequestHandler(ActionsMock{})
	result, err := handler.DoList(
		dto.Page{
			Size:  1,
			Start: 0,
		},
		dto.SortParams{
			Type:   dto.ByPriceChange,
			SortUp: true,
		})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if result[0].Name != "ByPriceChange" {
		t.Errorf("SortType ByPriceChange is expected, actual: %v", result[0].Name)
	}
}

func TestHandlerDoListByDate(t *testing.T) {
	handler := NewRequestHandler(ActionsMock{})
	result, err := handler.DoList(
		dto.Page{
			Size:  1,
			Start: 0,
		},
		dto.SortParams{
			Type:   dto.ByLastChange,
			SortUp: true,
		})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if result[0].Name != "ByLastChange" {
		t.Errorf("SortType ByLastChange is expected, actual: %v", result[0].Name)
	}
}
