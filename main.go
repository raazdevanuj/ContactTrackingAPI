package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func viewuser(w http.ResponseWriter, r *http.Request) {
	// articles := Articles{
	// 	Article{Title: "Hello", Desc: "Article Description", Content: "Article Content"},
	// }
	// fmt.Println("Endpoint Hit: all articles end")
	// json.NewEncoder(w).Encode(articles)

	fmt.Printf("Req: %s %s", r.URL, r.URL.Path)

}
func createusers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test Post created")
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
	fmt.Fprintf(w, "test contact get")
}
func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint Hit")
}
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homepage)
	myRouter.HandleFunc("/users", viewuser).Methods("GET")
	myRouter.HandleFunc("/users", createusers).Methods("POST")
	myRouter.HandleFunc("/contacts", contact).Methods("GET")
	myRouter.HandleFunc("/contacts", contact).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
func main() {
	handleRequests()
}
