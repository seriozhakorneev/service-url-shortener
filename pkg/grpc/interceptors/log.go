package interceptors

import (
	"context"
	"io"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

var defaultWriter io.Writer = os.Stdout

var Logs = func(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	l := grpclog.NewLoggerV2(defaultWriter, nil, nil)

	resp, err := handler(ctx, req)
	status := status.Code(err)
	if err != nil {
		l.Infof("| %d - %s |   %s |   %q\n%s\n",
			status, status, time.Since(start), info.FullMethod, err)
	} else {
		l.Infof("| %d - %s |   %s |   %q",
			status, status, time.Since(start), info.FullMethod)
	}

	return resp, err
}
