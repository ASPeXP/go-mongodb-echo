package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aspexp/go-mongodb-echo/product"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName = "consumer"
)

type RequestBody struct {
	Id string `json:"id"`
	Product_Id float64     `json:"product_id"`
	Product_Name         string     `json:"product_name"`
	Retail_Price  float64 `json:"retail_price"`
		
}
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
		// if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		// 	panic(err )
		// }
		fmt.Println("Pinged your deployments. You successfully connected to mongoDB!")

		serverPort := ":" + os.Getenv("PORT")

		e.GET("/api/v1/listdb", func (c echo.Context) error {
			database, err := client.ListDatabaseNames(context.TODO(), bson.M{})
			if err != nil {
				panic(err)
			}
			return c.JSON(http.StatusOK, database)
		})
		e.GET("/api/v1/products", func (c echo.Context) error {
			
			result := product.GetProduct(client, c)
			return result
		})
		e.POST("/api/v1/products", func (c echo.Context) error {
			
			var body RequestBody
			err := c.Bind(&body)
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
			p := product.Product{
				Product_Id: int64(body.Product_Id),
				Product_Name: body.Product_Name,
				Retail_Price: body.Retail_Price,
			}

			result, err  := product.PostProduct(p, client, c)
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
			return c.JSON(http.StatusOK, result )
		})
		e.PATCH("/api/v1/products", func (c echo.Context) error {
			var body RequestBody
			err := c.Bind(&body)
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
			id,_  := primitive.ObjectIDFromHex(body.Id)
			filter := bson.D{{"_id", id}}
			update := bson.D{{"$set", bson.D{{"Retail_Price", 29.10}, {"avg_rating", 4.4}}}}
				
			result := product.UpdateProduct(filter, update, client, c)
			// if err != nil {
			// 	return c.String(http.StatusBadRequest, err.Error())
			// }
			return c.JSON(http.StatusOK, result )
			
		})

	go func() {
		e.Logger.Fatal(e.Start(serverPort))
	}()
	log.Println("server started.")

	<-stop
	log.Println("shutting down the server")
}

