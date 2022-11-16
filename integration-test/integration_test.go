package integration_test

import (
	"net/http"
	"testing"

	. "github.com/Eun/go-hit"
)

const (
	//host       = "app:8080"
	host = "localhost:8080"

	// HTTP
	basePath = "http://" + host

	// RabbitMQ RPC
	rmqURL            = "amqp://guest:guest@rabbitmq:5672/"
	rpcServerExchange = "rpc_server"
	rpcClientExchange = "rpc_client"
	requests          = 10
)

// HTTP GET: /.
func TestHTTPRedirectGet(t *testing.T) {
	Test(t,
		Description("Redirect Get"),
		Get(basePath+"/<testlink>"),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().JSON().JQ(".error").Equal(
			"provided short URL too long or "+
				"impossible with current configurations",
		),
	)
}

// RabbitMQ RPC Client: getHistory.
//func TestRMQClientRPC(t *testing.T) {
//	rmqClient, err := client.New(rmqURL, rpcServerExchange, rpcClientExchange)
//	if err != nil {
//		t.Fatal("RabbitMQ RPC Client - init error - client.New")
//	}
//
//	defer func() {
//		err = rmqClient.Shutdown()
//		if err != nil {
//			t.Fatal("RabbitMQ RPC Client - shutdown error - rmqClient.RemoteCall", err)
//		}
//	}()
//
//	type Translation struct {
//		Source      string `json:"source"`
//		Destination string `json:"destination"`
//		Original    string `json:"original"`
//		Translation string `json:"translation"`
//	}
//
//	type historyResponse struct {
//		History []Translation `json:"history"`
//	}
//
//	for i := 0; i < requests; i++ {
//		var history historyResponse
//
//		err = rmqClient.RemoteCall("getHistory", nil, &history)
//		if err != nil {
//			t.Fatal("RabbitMQ RPC Client - remote call error - rmqClient.RemoteCall", err)
//		}
//
//		if history.History[0].Original != "текст для перевода" {
//			t.Fatal("Original != текст для перевода")
//		}
//	}
//}
