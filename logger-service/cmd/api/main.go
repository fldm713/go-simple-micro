package main

import (
	"context"
	"log"
	"logger/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort = "80"
	rpcPort = "5001"
	mongoURL = "mongodb://mongo:27017"
	grpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	// create a context inorder to disconnect from mongo
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// disconnect from mongo
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// register the RPC server
	err = rpc.Register(new(RPCServer))
	if err != nil {
		log.Panic(err)
	}
	go app.rpcListen()

	go app.grpcListen()
	// start web server
	log.Println("Web server started on port", webPort)
	srv := &http.Server{
		Addr: ":" + webPort,
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (app *Config) rpcListen() error {
	log.Println("RPC server started on port", rpcPort)
	listen, err := net.Listen("tcp", "0.0.0.0:" + rpcPort)
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			log.Println("Error accepting connection", err)
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

func connectToMongo() (*mongo.Client, error) {
	// create connect options
	clientOpitons := options.Client().ApplyURI(mongoURL)
	clientOpitons.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect to mongo
	c, err := mongo.Connect(context.TODO(), clientOpitons)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	log.Println("Connected to MongoDB")
	return c, nil
}