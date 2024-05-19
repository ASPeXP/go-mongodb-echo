package product

import (
	"context"
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Product struct {
	ID primitive.ObjectID `bson:"_id"`
	Product_Id int64 `bson:"product_id,omitempty"`
	Product_Name string `bson:"product_name,omitempty"`
	Retail_Price float64 `bson:"retail_price,omitempty"`			
}
func GetProduct(client *mongo.Client, c echo.Context) error{

			// var result []bson.D
			var result []Product 
			consumerDatabase := client.Database("consumer")
			productCollection := consumerDatabase.Collection("products")
			cursor, err := productCollection.Find(context.TODO(), bson.D{})
			if err != nil {
				panic(err)
			}
			for cursor.Next(context.TODO()){
				var p Product
				if err := cursor.Decode(&p); err != nil {
					log.Fatal(err)
				}
				result = append(result, p)

			}
			return c.JSON(http.StatusOK, result )

}

func PostProduct(product Product, client *mongo.Client, c echo.Context) (interface{}, error ) {

			consumerDatabase := client.Database("consumer")
			productCollection := consumerDatabase.Collection("products")
			productResult, err := productCollection.InsertOne(context.TODO(), product)
			if err != nil {
				panic(err)
			}
			return productResult.InsertedID, nil 
}