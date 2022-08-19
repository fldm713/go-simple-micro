package main

import (
	"context"
	"log"
	"logger/data"
	"logger/logs"
	"net"

	"google.golang.org/grpc"
)


type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// write a log
	logEntry := data.LogEntry{
		Name: input.GetName(),
		Data: input.GetData(),
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{
			Result: "failed",
		}
		return res, err
	}

	// return a response
	res := &logs.LogResponse{
		Result: "logged",
	}
	return res, nil
}

func (app *Config) grpcListen() {
	listen, err := net.Listen("tcp", ":" + grpcPort)
	if err != nil {
		log.Fatalf("failed to listen to grpc: %v", err)
	}

	server := grpc.NewServer()

	logs.RegisterLogServiceServer(server, &LogServer{Models: app.Models})

	log.Println("grpc server listening on port", grpcPort)

	if err = server.Serve(listen); err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	}
}