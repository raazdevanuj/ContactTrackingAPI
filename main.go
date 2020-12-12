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
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//users .....
type Users struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	DateOfBirth string    `json:"date_of_birth"`
	PhoneNumber int32     `json:"phone_number"`
	Email       string    `json:"email"`
	Timestamp   time.Time `json:"timestamp"`
}

//contacts...
type Contacts struct {
	UserIdOne int       `json:"user_id_one"`
	UserIdTwo int       `json:"user_id_two"`
	Timestamp time.Time `json:"timestamp"`
}

var clientInstance *mongo.Client
var clientInstanceError error
var mongoOnce sync.Once

const (
	CONNECTIONSTRING = "mongodb://localhost:27017"
	DB               = "db_issue_manager"
	ISSUES           = "col_issues"
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
	filter := bson.D{primitive.E{Key: "ID", Value: i}}
	client, err := GetMongoClient()
	if err != nil {
		return result, err
	}
	collection := client.Database(DB).Collection(ISSUES)
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
func createuserutil(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test Post created")
}
func viewuser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		viewuserutil(w, r)
		w.Write([]byte("Received a GET request\n"))
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}
}
func createuser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		createuserutil(w, r)
		w.Write([]byte("Received a POST request\n"))
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
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func main() {
	handleRequests()
}
