package server

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

var Logs = func(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	l := grpclog.NewLoggerV2(defaultWriter, nil, nil)

	resp, err := handler(ctx, req)
	if err != nil {
		l.Infof("method %q failed: %s", info.FullMethod, err)
	} else {
		l.Infof("method %q success: %s", info.FullMethod, status.Code(err))
	}

	return resp, err
}
