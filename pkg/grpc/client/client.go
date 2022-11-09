package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	pb "service-url-shortener/internal/entrypoint/grpc/shortener_proto"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.example.com", "The server name used to verify the hostname returned by the TLS handshake")
)

// TODO client construction for tests

func createShort(client pb.ShortenerClient, data *pb.ShortenerData) {
	log.Printf("Creating short url for (%s)", data.GetURL())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	data, err := client.Create(ctx, data)
	if err != nil {
		log.Fatalf("client.Create failed: %v", err)
	}
	log.Println(data.URL)
}

func getOriginal(client pb.ShortenerClient, data *pb.ShortenerData) {
	log.Printf("Getting url for (%s)", data.GetURL())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	data, err := client.Get(ctx, data)
	if err != nil {
		log.Fatalf("client.Get failed: %v", err)
	}
	log.Println(data.URL)
}

func main() {
	var opts []grpc.DialOption
	if *tls {
		//if *caFile == "" {
		//	*caFile = data.Path("x509/ca_cert.pem")
		//}
		// TODO: credentials
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewShortenerClient(conn)

	// Looking for a URL
	getOriginal(client, &pb.ShortenerData{URL: "http://127.0.0.1:8080/Bффывфыв"})
	//createShort(client, &pb.ShortenerData{URL: "https://apps.apple.com/ru/app/football-manager-2023-touch/id1626267810"})
}
