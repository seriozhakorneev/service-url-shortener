package integration_test

import (
	"context"
	"net/http"
	"testing"

	. "github.com/Eun/go-hit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "service-url-shortener/internal/entrypoint/grpc/shortener_proto"
)

const (
	container = "app"

	// HTTP
	basePath = "http://" + container + ":8080"

	// GRPC
	serverAddr = container + ":50051"
)

// HTTP GET: /.
func TestHTTPRedirectGet(t *testing.T) {
	Test(t,
		Description("Redirect Get"),
		Get(basePath+"/testlink"),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().JSON().JQ(".error").Equal(
			"provided short URL too long or "+
				"impossible with current configurations",
		),
	)
}

// GRPC: Shortener.
func TestGRPCShortener(t *testing.T) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("GRPC client - grpc.Dial: %v", err)
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			t.Fatal("GRPC client - conn.Close()", err)
		}
	}()

	data := &pb.ShortenerData{URL: "not_valid_url"}
	expectedStatus := codes.InvalidArgument

	client := pb.NewShortenerClient(conn)

	result, err := client.Create(context.Background(), data)
	if result != nil {
		t.Fatalf("Expected nil in Get result, Got: %v", result)
	}

	code := status.Code(err)
	if code != expectedStatus {
		t.Fatalf("Expected status code in Create: %s, Got: %s", expectedStatus, code)
	}

	result, err = client.Get(context.Background(), data)
	if result != nil {
		t.Fatalf("Expected nil in Get result, Got: %v", result)
	}

	code = status.Code(err)

	if code != expectedStatus {
		t.Fatalf("Expected status code in Get: %s, Got: %s", expectedStatus, code)
	}
}
