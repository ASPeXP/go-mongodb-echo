package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName = "consumer"
)

func main(){

	e := echo.New()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI")).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	defer func(){
		if err = client.Disconnect(context.TODO()); err != nil {
			panic  (err )
		}
	}()

		if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
			panic (err )
		}
		fmt.Println("Pinged your deployments. You successfully connected to mongoDB!")

		serverPort := ":" + os.Getenv("PORT")

		e.GET("/api/v1/product", )

	go func() {
		e.Logger.Fatal(e.Start(serverPort))
	}()
	log.Println("server started.")

	<-stop
	log.Println("shutting down the server")
}