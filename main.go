package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

//users .....
type Users struct {
	ID          int    `json:"id"`
	NAME        string `json:"name"`
	DateOfBirth string `json:"dateofbirth"`
	PhoneNumber int32  `json:"phonenumber"`
	EMAIL       string `json:"email"`
	Timestamp   string `json:"timestamp"`
}

//contacts...
type Contacts struct {
	USERIDONE int    `json:"useridone"`
	USERIDTWO int    `json:"useridtwo"`
	Timestamp string `json:"timestamp"`
}
type userss []Users

// fmt.Println("Endpoint Hit: all articles end")
// json.NewEncoder(w).Encode(articles)

func viewuserutil(w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("Req: %s %s", r.URL, r.URL.Path)
	var a string = r.URL.Path

	fmt.Println(path.Base(a))
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
