package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	pb "grpc-filesharing/fileSearch"

	"google.golang.org/grpc"
)

func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil { 
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("Não foi possível obter um endereço IP")
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Uso: %s <command> [args...]", os.Args[0])
	}

	command := os.Args[1]
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Não conseguiu se conectar: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileSearchClient(conn)

	switch command {
	case "send":
		if len(os.Args) < 3 {
			log.Fatalf("Uso: %s send <file_hash_1> <file_hash_2> ...", os.Args[0])
		}
		hashes := os.Args[2:]
		sendFileHashes(client, hashes)
	case "search":
		if len(os.Args) != 3 {
			log.Fatalf("Uso: %s search <file_hash>", os.Args[0])
		}
		hash := os.Args[2]
		getMachinesWithFile(client, hash)
	}
}

func sendFileHashes(client pb.FileSearchClient, hashes []string) {
	ip, err := GetLocalIP()
	if err != nil {
		log.Fatalf("Erro ao obter IP: %v", err)
	}

	log.Printf("Enviando hashes com IP: %s\n", ip)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = client.SendFileHashes(ctx, &pb.FileHashes{Hashes: hashes})
	if err != nil {
		log.Fatalf("Erro ao enviar hashes: %v", err)
	}
	log.Println("Hashes enviados com sucesso")
}

func getMachinesWithFile(client pb.FileSearchClient, hash string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetMachinesWithFile(ctx, &pb.FileHash{Hash: hash})
	if err != nil {
		log.Fatalf("Erro ao buscar IPs: %v", err)
	}

	log.Printf("Máquinas com o arquivo %s: %v", hash, res.Ips)
}
