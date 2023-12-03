package ninjasApi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	limit  = 1
	apiURL = "https://api.api-ninjas.com/v1/facts?limit=%d"
	apiKey = "your api key"
)

type NinjaFact struct {
	client *mongo.Client
}

func NewNinjaFact(client *mongo.Client) *NinjaFact {
	return &NinjaFact{
		client: client,
	}
}

func (nf *NinjaFact) start() error {
	url := fmt.Sprintf(apiURL, limit)

	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("X-Api-Key", apiKey)

	collection := nf.client.Database("database").Collection("name")
	ticker := time.NewTicker(2 * time.Second)

	for {
		response, err := client.Do(request)
		if err != nil {
			log.Println("Error making API request:", err)
			continue
		}

		if response.StatusCode != http.StatusOK {
			log.Println("API request failed. Status code:", response.StatusCode)
			body, _ := io.ReadAll(response.Body)
			log.Println("Error response:", string(body))
			continue
		}

		var dataNinja []bson.M
		if err = json.NewDecoder(response.Body).Decode(&dataNinja); err != nil {
			log.Println("Error decoding JSON:", err)
			continue
		}

		for _, doc := range dataNinja {
			_, err = collection.InsertOne(context.TODO(), doc)
			if err != nil {
				log.Println("Error inserting into MongoDB:", err)
				continue
			}
		}

		log.Println("Data inserted into MongoDB successfully.")

		<-ticker.C
	}
}

type Server struct {
	client *mongo.Client
}

func NewServer(client *mongo.Client) *Server {
	return &Server{
		client: client,
	}
}

func (s *Server) GetNinjaFact(w http.ResponseWriter, r *http.Request) {
	collection := s.client.Database("database").Collection("name")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Println("Error finding documents in MongoDB:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var results []bson.M

	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Println("Error decoding MongoDB documents:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(results); err != nil {
		log.Println("Error encoding JSON response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
		return
	}

	ninjaFact := NewNinjaFact(client)
	go ninjaFact.start()

	server := NewServer(client)
	http.HandleFunc("/facts", server.GetNinjaFact)
	http.ListenAndServe(":3000", nil)
}
