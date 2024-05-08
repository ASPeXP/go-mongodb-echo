package product

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"

	"github.com/aspexp/go-mongodb-echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func readProduct(c echo.Context){
	response := map[string]interface{}{}
	client, err := mongo.Connect(c, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		fmt.Println(err.Error())
	}

	collection := client.Database(dbName).Collection("products")

}