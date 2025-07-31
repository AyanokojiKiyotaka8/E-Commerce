package main

import (
	"context"
	"log"
	"net/http"

	"github.com/AyanokojiKiyotaka8/E-Commerce/product_service/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := proto.RegisterProductServiceHandlerFromEndpoint(context.Background(), mux, ":50051", opts); err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal(err)
	}
}
