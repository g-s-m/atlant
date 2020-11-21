package service

import (
	srv "atlant/generated/interface"
	dto "atlant/service/dto"
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

type IRequestHandler interface {
	DoFetch(file string) error
	DoList(sortBy dto.SortType, sortUp bool) ([]dto.Product, error)
}

type server struct {
	srv.UnimplementedProductServiceServer
	worker IRequestHandler
}

func (s *server) Fetch(ctx context.Context, request *srv.FetchRequest) (*srv.FetchReply, error) {
	doneCh := make(chan srv.FetchReply_Status)

	go func() {
		doneCh <- dto.ErrorToStatus(s.worker.DoFetch(request.GetUrl()))
	}()

	status := srv.FetchReply_OK
	select {
	case <-ctx.Done():
		log.Printf("Operation was canceled")
		status = srv.FetchReply_RESOURCE_UNAVAILABLE
	case status = <-doneCh:
	}

	return &srv.FetchReply{
		Status: status,
	}, nil
}

func (s *server) List(ctx context.Context, request *srv.ListRequest) (*srv.ListReply, error) {
	prod := dto.ProductDto(dto.Product{
		Name:        "q",
		Price:       1.23,
		ChangeCount: 2,
		ChangeDate:  time.Now(),
	})
	return &srv.ListReply{
		ProductList: []*srv.Product{
			&prod,
		},
	}, nil
}

func RunService(port string) {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Can't start a server: %v", err)
	}
	s := grpc.NewServer()
	srv.RegisterProductServiceServer(s, &server{})

	term := make(chan os.Signal)
	signal.Notify(term, syscall.SIGTERM, syscall.SIGINT)

	errCh := make(chan error)

	go func() {
		if err := s.Serve(listen); err != nil {
			errCh <- err
		}
	}()
	defer s.GracefulStop()
	select {
	case err := <-errCh:
		log.Fatalf("failed to serve: %v", err)
	case <-term:
	}
}
