package service

import (
	"google.golang.org/grpc"
	req "atlant/requestor"
	srv "atlant/generated/interface"
	"context"
	"os"
	"net"
	"log"
	"os/signal"
	"syscall"
)

type server struct {
	srv.UnimplementedProductServiceServer
}
/*
struct RequestHandler {
	Requestor r
	Processor p
	Storage s
}

func (p *RequestHandler) DoFetch(path string) error {
	r.GetFile(url)
	p.Process()
	s.Save()
}

func (p *RequestHandler) DoList() error {
}
*/
func (s *server) Fetch(ctx context.Context, request *srv.FetchRequest) (*srv.FetchReply, error) {
	doneCh := make(chan srv.FetchReply_Status)

	go func() {
		file, err := req.GetCsvFile(request.GetUrl(), 30)
		//storage.Save(proc.ProcessFile(file))
		if err != nil {
			log.Printf("Error during getting csv file: %v", err)
			doneCh <-srv.FetchReply_RESOURCE_UNAVAILABLE
		}
		doneCh <- srv.FetchReply_OK
	}()
/*
	go func() {
		s.handler.DoFetch(url)
	}()*/
	status := srv.FetchReply_OK
	select {
		case <-ctx.Done():
			log.Printf("Operation timed out")
			status := srv.FetchReply_RESOURCE_UNAVAILABLE
		case status := <-doneCh:
	}

	return &srv.FetchReply {
		Status : status,
	}, nil
}

func (s* server) List(ctx context.Context, request *srv.ListRequest) (*srv.ListReply, error) {
	return &srv.ListReply {
		Product : "p",
		Price : 12.1234,
		Timestamp : 123456,
		Changed : 4,
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

