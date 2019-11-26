package main

import (
	bl10 "bl10server/bl10comms"
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type bl10LockServer struct {
}

func (s *bl10LockServer) Unlock(ctx context.Context, lock *bl10.Lock) (*bl10.LockStatus, error) {
	err := Unlock(lock.GetImei())
	if err != nil {
		return nil, err
	}

	// No feature was found, return an unnamed feature
	return &bl10.LockStatus{IsCharching: true, IsLocked: true, Imei: lock.GetImei()}, nil
}

func startGrpcServer() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9090))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	bl10.RegisterBL10LockServer(grpcServer, &bl10LockServer{})
	grpcServer.Serve(lis)
}
