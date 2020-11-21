package service

import (
	srv "atlant/generated/interface"
	dto "atlant/service/dto"
	"atlant/service/handler"
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

type IRequestHandler interface {
	DoFetch(file string) error
	DoList(page dto.Page, sort dto.SortParams) ([]dto.Product, error)
}

type Server struct {
	srv.UnimplementedProductServiceServer
	worker IRequestHandler
}

func (s *Server) Fetch(ctx context.Context, request *srv.FetchRequest) (*srv.FetchReply, error) {
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

func (s *Server) List(ctx context.Context, request *srv.ListRequest) (*srv.ListReply, error) {
	doneCh := make(chan struct{})
	products := []dto.Product{}
	go func() {
		var err error
		products, err = s.worker.DoList(dto.PageDto(request.GetPage()), dto.SortDto(request.GetSort()))
		if err != nil {
			log.Printf("List products failed: %v", err)
		}
		doneCh <- struct{}{}
	}()

	var result []*srv.Product
	select {
	case <-ctx.Done():
		log.Printf("Operation was canceled")
	case <-doneCh:
		for _, p := range products {
			product := dto.ProductDto(p)
			result = append(result, &product)
		}
	}

	return &srv.ListReply{
		ProductList: result,
	}, nil
}

func (s *Server) Run(port string) {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Can't start a server: %v", err)
	}
	grpcServer := grpc.NewServer()
	srv.RegisterProductServiceServer(grpcServer, s)

	term := make(chan os.Signal)
	signal.Notify(term, syscall.SIGTERM, syscall.SIGINT)

	errCh := make(chan error)

	go func() {
		if err := grpcServer.Serve(listen); err != nil {
			errCh <- err
		}
	}()
	defer grpcServer.GracefulStop()
	select {
	case err := <-errCh:
		log.Fatalf("failed to serve: %v", err)
	case <-term:
	}
}

func NewServer(h IRequestHandler) *Server {
	return &Server{
		worker: h,
	}
}

func RunService(port string) {
	s := NewServer(handler.RequestHandler{})
	s.Run(port)
}
