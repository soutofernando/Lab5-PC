package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "grpc-filesharing/fileSearch"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

type server struct {
	pb.UnimplementedFileSearchServer
	mu       sync.Mutex
	fileMap  map[string][]string
	reqChan  chan request
	respChan chan response
}

type request struct {
	hash   string
	client string
}

type response struct {
	hash string
	ips  []string
}

func newServer() *server {
	return &server{
		fileMap:  make(map[string][]string),
		reqChan:  make(chan request),
		respChan: make(chan response),
	}
}

func (s *server) SendFileHashes(ctx context.Context, req *pb.FileHashes) (*pb.Response, error) {
	clientIP, _ := getClientIP(ctx)
	
	for _, hash := range req.Hashes {
		s.mu.Lock()
		s.fileMap[hash] = append(s.fileMap[hash], clientIP)
		s.mu.Unlock()
	}
	
	return &pb.Response{Message: "Hashes received"}, nil
}

func (s *server) GetMachinesWithFile(ctx context.Context, req *pb.FileHash) (*pb.FileLocations, error) {
	s.mu.Lock()
	ips := s.fileMap[req.Hash]
	s.mu.Unlock()
	return &pb.FileLocations{Ips: ips}, nil
}

func getClientIP(ctx context.Context) (string, error) {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("erro ao obter IP do cliente")
	}
	return pr.Addr.String(), nil
}

func (s *server) run() {
	for {
		select {
		case req := <-s.reqChan:
			s.mu.Lock()
			ips := s.fileMap[req.hash]
			s.mu.Unlock()
			s.respChan <- response{hash: req.hash, ips: ips}
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	s := newServer()
	go s.run()

	pb.RegisterFileSearchServer(grpcServer, s)
	log.Printf("Servidor rodando na porta :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
