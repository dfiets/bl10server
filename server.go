package main

import (
	bl10 "bl10server/bl10comms"
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

type bl10LockServer struct {
}

func (s *bl10LockServer) Unlock(ctx context.Context, lock *bl10.Lock) (*bl10.LockStatus, error) {
	err := SendCommandToLock(lock.GetImei(), "UNLOCK#")
	if err != nil {
		return nil, err
	}

	// No feature was found, return an unnamed feature
	return &bl10.LockStatus{IsCharching: true, IsLocked: true, Imei: lock.GetImei()}, nil
}

func (s *bl10LockServer) Alarm(ctx context.Context, lock *bl10.Lock) (*bl10.LockStatus, error) {
	err := SendCommandToLock(lock.GetImei(), "SDFIND,ON,5,15,1#")
	if err != nil {
		return nil, err
	}

	// No feature was found, return an unnamed feature
	return &bl10.LockStatus{IsCharching: true, IsLocked: true, Imei: lock.GetImei()}, nil
}

func (s *bl10LockServer) Locate(ctx context.Context, lock *bl10.Lock) (*bl10.LockStatus, error) {
	err := SendCommandToLock(lock.GetImei(), "LJDW#")
	if err != nil {
		return nil, err
	}

	// No feature was found, return an unnamed feature
	return &bl10.LockStatus{IsCharching: true, IsLocked: true, Imei: lock.GetImei()}, nil
}

func (s *bl10LockServer) StatusUpdates(empty *empty.Empty, stream bl10.BL10Lock_StatusUpdatesServer) error {
	log.Println("new connection")
	addConsumer(stream)
	return nil

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
