package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	handlers "github.com/paisit04/recipes-api/handlers"
)

var recipesHandler *handlers.RecipesHandler

func init() {

	ctx := context.Background()
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	// insertData(ctx, collection)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping(ctx)
	fmt.Println(status)

	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
}

type RecipeData struct {
	Name         string    `json:"name" bson:"name"`
	Tags         []string  `json:"tags" bson:"tags"`
	Ingredients  []string  `json:"ingredients" bson:"ingredients"`
	Instructions []string  `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time `json:"publishedAt" bson:"publishedAt"`
}

func insertData(ctx context.Context, collection *mongo.Collection) {
	recipes := make([]RecipeData, 0)
	file, _ := os.ReadFile("recipes.json")
	_ = json.Unmarshal(file, &recipes)

	var listOfRecipes []interface{}
	for _, recipe := range recipes {
		listOfRecipes = append(listOfRecipes, recipe)
	}
	insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted recipes: ", len(insertManyResult.InsertedIDs))
}

func main() {
	router := gin.Default()

	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.GET("/recipes/:id", recipesHandler.GetOneRecipeHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)

	// router.GET("/recipes/search", SearchRecipesHandler)

	router.Run() // listen and serve on 0.0.0.0:8080
}
