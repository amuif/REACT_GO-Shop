package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Watches struct {
	ID               primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Sold             bool               `json:"sold"`
	Description      string             `json:"description"`
	Image            string             `json:"image"`
	Availible_Colors []string           `json:"colors"`
	Price            int                `json:"price"`
}

var collection *mongo.Collection

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file", err)
	}
	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal("error while connecting to DB", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	fmt.Println("Connected TO MONGODB Atlas")
	collection = client.Database("watchShop_DB").Collection("watches")

	app := fiber.New()

	app.Get("/api/watches", getWatches)
	app.Post("/api/watches", addWatches)
	app.Patch("/api/watches/:id", UpdateWatches)
	app.Delete("/api/watches/:id", DeleteWatches)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "5000"
	}
	log.Fatal(app.Listen("0.0.0.0:" + PORT))
}

func getWatches(c *fiber.Ctx) error {
	var watches []Watches

	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var watch Watches
		if err := cursor.Decode(&watch); err != nil {
			fmt.Println("error")
			return err
		}
		fmt.Println("no error")
		watches = append(watches, watch)

	}
	return c.JSON(watches)

}
func addWatches(c *fiber.Ctx) error {
	watch := new(Watches)
	if err := c.BodyParser(watch); err != nil {
		return err
	}

	if watch.Description == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Watch Should Have a description and an Image"})
	}
	insertResult, err := collection.InsertOne(context.Background(), watch)
	if err != nil {
		return err
	}

	watch.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(200).JSON(fiber.Map{"status": "watch have successfully been added"})

}

func UpdateWatches(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"msg": "invalid ID"})

	}
	filter := bson.M{"_id": objectID}

	update := bson.M{"$set": bson.M{"sold": true}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"success": "Item Have Been Updated"})

}

func DeleteWatches(c *fiber.Ctx)error{
	id := c.Params("id")
	objectID , err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error":"There is no such Watch"}) 
	}
	filter := bson.M{"_id": objectID}
	_ , err = collection.DeleteOne(context.Background(),filter)

	if err != nil {
		return err 
	}

	return c.Status(200).JSON(fiber.Map{"success":true})
}
