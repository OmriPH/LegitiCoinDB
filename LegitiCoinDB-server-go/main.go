package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Player struct {
	Uuid   string `json:"uuid"`
	Lcoins int    `json:"lcoins"`
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	log.Println("Got request for /")
}

func top10(w http.ResponseWriter, players *mongo.Collection) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find top 10 players sorted by score in descending order
	findOptions := options.Find().SetLimit(10).SetSort(bson.D{{"Lcoins", -1}})
	cursor, err := players.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		http.Error(w, "Error fetching top players", http.StatusInternalServerError)
		log.Printf("Error fetching top players: %v", err)
		return
	}
	defer cursor.Close(ctx)

	var results []Player
	if err = cursor.All(ctx, &results); err != nil {
		http.Error(w, "Error processing players", http.StatusInternalServerError)
		log.Printf("Error processing players: %v", err)
		return
	}

	json.NewEncoder(w).Encode(map[string][]Player{"result": results})
	log.Println("Got request for /list-top/10")
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get configuration from environment
	port := os.Getenv("PORT")
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB")

	// Set up MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(clientOptions)

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Could not ping MongoDB: %v", err)
	}

	// Get database collections
	players := client.Database(dbName).Collection("players")

	log.Println("MongoDB Connected")

	// Set up HTTP routes
	http.HandleFunc("/", root)
	http.HandleFunc("/list-top/10", func(w http.ResponseWriter, r *http.Request) {
		top10(w, players)
	})

	// Start server
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
