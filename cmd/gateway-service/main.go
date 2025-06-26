package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sammyqtran/url-shortener/internal/gateway"
	pb "github.com/sammyqtran/url-shortener/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.NewClient("url-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	grpcClient := pb.NewURLServiceClient(conn)

	server := &gateway.GatewayServer{
		GrpcClient: grpcClient,
	}

	r := mux.NewRouter()

	r.HandleFunc("/create", server.HandleCreateShortURL).Methods("POST")
	r.HandleFunc("/healthz", server.HandleHealthCheck).Methods("GET")
	r.HandleFunc("/{shortCode}", server.HandleGetOriginalURL).Methods("GET")

	log.Println("Gateway service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
