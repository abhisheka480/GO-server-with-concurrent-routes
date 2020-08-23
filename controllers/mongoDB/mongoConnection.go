package mongoDB

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/GO-server-with-concurrent-routes/config"
	"github.com/GO-server-with-concurrent-routes/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Collection {

	// Set client options
	clientOptions := options.Client().ApplyURI(config.MONGO_URL)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database(config.DB_NAME).Collection(config.COLLECTION_NAME)

	return collection
}

func GetError(err error, w http.ResponseWriter) {

	//log.Fatal(err.Error())
	var response = models.ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusBadRequest,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}
