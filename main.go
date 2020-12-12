package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strconv"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Users ...
type Users struct {
	ID          int64  `json:"_id" bson:"_id,omitempty"`
	Name        string `json:"name" bson:"name,omitempty"`
	DateOfBirth string `json:"date_of_birth" bson:"date_of_birth,omitempty"`
	PhoneNumber string `json:"phone_number" bson:"phone_number,omitempty"`
	Email       string `json:"email" bson:"email,omitempty"`
	Timestamp   string `json:"timestamp" bson:"timestamp,omitempty"`
}

// Contacts ...
type Contacts struct {
	UserOne   int64  `json:"user_one"`
	UserTwo   int64  `json:"user_two"`
	Timestamp string `json:"timestamp"`
}

var clientInstance *mongo.Client
var clientInstanceError error
var mongoOnce sync.Once

const (
	CONNECTIONSTRING = "mongodb://localhost:27017"
	DB               = "appointy"
)

func GetMongoClient() (*mongo.Client, error) {
	mongoOnce.Do(func() {

		clientOptions := options.Client().ApplyURI(CONNECTIONSTRING)

		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
		}

		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
		}
		clientInstance = client
	})
	return clientInstance, clientInstanceError
}

/* fmt.Println("Endpoint Hit: all articles end")
json.NewEncoder(w).Encode(articles)*/

func getbyid(code string) (Users, error) {
	result := Users{}
	i, err := strconv.Atoi(code)
	filter := bson.D{primitive.E{Key: "_id", Value: i}}
	client, err := GetMongoClient()
	if err != nil {
		return result, err
	}
	collection := client.Database(DB).Collection("users")
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func viewuserutil(w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("Req: %s %s", r.URL, r.URL.Path)
	var a string = r.URL.Path
	var user, err = getbyid(path.Base(a))
	if err != nil {
		fmt.Fprintf(w, "Error")
		return
	}
	json.NewEncoder(w).Encode(user)
}

func viewuser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		viewuserutil(w, r)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}
}
func createuser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		w.Header().Set("Content-Type", "application/json")

		var user Users

		// we decode our body request params
		_ = json.NewDecoder(r.Body).Decode(&user)
		client, err := GetMongoClient()
		if err != nil {
			return
		}

		collection := client.Database(DB).Collection("users")

		var lastUser Users
		findOptions := options.FindOne()
		findOptions.SetSort(bson.D{{"_id", -1}})

		err2 := collection.FindOne(context.TODO(), bson.D{}, findOptions).Decode(&lastUser)
		if err2 != nil {
			user.ID = 1
		}

		user.ID = lastUser.ID + 1
		// insert our user model.
		result, err := collection.InsertOne(context.TODO(), user)

		if err != nil {
			log.Fatal(result)
			log.Fatal(err)
			return

		}

		json.NewEncoder(w).Encode(user)
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}
}
func contact(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		for k, v := range r.URL.Query() {
			fmt.Printf("%s: %s\n", k, v)
		}
		w.Write([]byte("Received a GET request\n"))
	case "POST":
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", reqBody)
		w.Write([]byte("Received a POST request\n"))
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}
}
func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint Hit")
}
func handleRequests() {
	http.HandleFunc("/", homepage)
	http.HandleFunc("/users", createuser)
	http.HandleFunc("/users/", viewuser)
	http.HandleFunc("/contacts", contact)

	fmt.Println("Serving API on port 8080")
	fmt.Println("URL for local testing http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
func main() {
	handleRequests()
}
