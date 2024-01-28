package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/kleczynski/go-microservices-k8s/database"
	"github.com/kleczynski/go-microservices-k8s/details"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

type DetailsResponse struct {
	Hostname         string `json:"hostname"`
	IP               string `json:"ip"`
	OperationSystem  string `json:"operationaSystem"`
	CPUCount         string `json:"cpuCount"`
	MemoryAlloc      string `json:"memoryAlloc,omitempty"`
	MemoryTotalAlloc string `json:"memoryTotalAlloc,omitempty"`
}

type Bio struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var BioData = make([]Bio, 0)

func HealthHanler(w http.ResponseWriter, r *http.Request) {
	log.Println("Checking application health")
	response := HealthResponse{
		Status:    "UP",
		Timestamp: time.Now().Format(time.RFC3339),
	}
	w.WriteHeader(http.StatusOK)
	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s", responseBytes)

}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving the homepage")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Application is up and running")
}

func DetailsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching details")
	hostname, err := details.GetHostName()
	if err != nil {
		log.Fatal(err)
	}
	ip, err := details.GetIP()
	if err != nil {
		log.Fatal(err)
	}
	operationSystem := details.GetOperatingSystem()
	CPUCount := details.GetCPUCount()
	strCPUCount := strconv.Itoa(CPUCount)
	response := DetailsResponse{
		Hostname:        hostname,
		IP:              ip.String(),
		OperationSystem: operationSystem,
		CPUCount:        strCPUCount,
	}
	w.WriteHeader(http.StatusOK)
	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s", responseBytes)
}

func APIListHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving API handler list")
	fmt.Fprintf(w, "List of API endpoints\n/api/v1/create\n/api/v1/read\n/api/v1/update\n/api/v1/delete")
}

func Create(w http.ResponseWriter, r *http.Request) {
	conn := database.Client.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_COLLECTION"))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var human Bio
	if err := json.NewDecoder(r.Body).Decode(&human); err != nil {
		log.Fatalf("There was an error decoding request: err - %s\n", err)
	}
	result, err := conn.InsertOne(context.TODO(), human)
	if err != nil {
		log.Fatalf("There was an error InsertingOne into MongoDB in Create method: %v", err)
	}
	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)

}

func Read(w http.ResponseWriter, r *http.Request) {
	conn := database.Client.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_COLLECTION"))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	cursor, err := conn.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatalf("There was an error in Finding documents in MongoDB for Read method: %v", err)
	}

	var results []Bio
	if err := cursor.All(context.TODO(), &results); err != nil {
		log.Fatalf("There was an error in access all of the documents from MongoDB: %v", err)
	}

	for _, result := range results {
		res, _ := json.Marshal(result)
		fmt.Println(string(res))
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	conn := database.Client.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_COLLECTION"))
	w.Header().Set("Content-Type", "application/json")
	var human Bio
	if err := json.NewDecoder(r.Body).Decode(&human); err != nil {
		log.Fatalf("There was an error decoding struct: err - %s\n", err)
	}
	filter := bson.D{{Key: "name", Value: human.Name}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: human.Name}, {Key: "age", Value: human.Age}}}}

	result, err := conn.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatalf("There was an error Updating document: %v", err)
	}
	fmt.Printf("Updated document: %d\n", result.ModifiedCount)

}

func Delete(w http.ResponseWriter, r *http.Request) {
	conn := database.Client.Database(os.Getenv("MONGO_DATABASE")).Collection(os.Getenv("MONGO_COLLECTION"))
	w.Header().Set("Content-Type", "application/json")
	name := mux.Vars(r)["name"]
	opts := options.Delete().SetCollation(&options.Collation{
		CaseLevel: false,
	})
	log.Println(name)
	result, err := conn.DeleteOne(context.TODO(), bson.D{{Key: "name", Value: name}}, opts)
	if err != nil {
		log.Fatalf("There was an error Deleting document: %v", err)
	}
	fmt.Printf("Deleted document: %d\n", result.DeletedCount)
}
