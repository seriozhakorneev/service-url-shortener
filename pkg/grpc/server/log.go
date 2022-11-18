package server

import (
	"context"
	"fmt"

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

	fmt.Println(info.Server)

	resp, err := handler(ctx, req)
	if err != nil {

		//[GIN] 2022/11/18 - 15:10:18 | 400 |    2.999625ms |      172.27.0.4 | GET      "/testlink"

		l.Infof("method %q failed: %s", info.FullMethod, err)
	} else {
		l.Infof("method %q success: %s", info.FullMethod, status.Code(err))
	}

	return resp, err
}
