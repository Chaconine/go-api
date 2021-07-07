package main

import (
	"encoding/json"	
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

type Article struct {
	Id string `json:"Id"`
	Title string `json:"Title"`
	Desc string `json:"desc"`
	Content string `json:"content"`
}

// let's declare a global Articles array
// that we can then populate in our main function
// to simulate a database
var Articles []Article

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Article struct
	// append this to our Articles array instead of just 
	// returning string expression
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article)
	// update our global Articles array to include 
	// our new article
	Articles = append(Articles, article)

	json.NewEncoder(w).Encode(article)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	// we need to parse the path parameters
	vars := mux.Vars(r)
	// we will need to extract the `id` of the article
	// we wish to delete
	id := vars["id"]

	// we then need to loop through all of our articles
	for index, article := range Articles{
		if article.Id == id {
			// updates our Articles array to remove
			// the article
			Articles = append(Articles[:index], Articles[index+1:]...)
			fmt.Println("testing...")
		}
	}
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	// Parse JSON from request
	// we need to unmarshal the json into the new
	// updated Articles array replacing the current
	// Article
	reqBody, _ := ioutil.ReadAll(r.Body)
	var updated Article
	json.Unmarshal(reqBody, &updated)

	// we need to parse the path parameters
	vars := mux.Vars(r)
	// we will need to extract the `id` of the article
	// we wish to delete
	id := vars["id"]

	// we then need to loop through all of our articles
	for index, article := range Articles{
		if article.Id == id {
			// update our Article array to replace 
			// the article
			Articles[index] = updated
		}
	}
}


func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/articles", returnAllArticles)
	// NOTE: ordering is important here! This has to be defined before
	// the other `article` endpoint
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	// add our new DELETE endpoint here
	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	// add our new UPDATE endpoint here
	myRouter.HandleFunc("/article/{id}", updateArticle).Methods("PUT")
	myRouter.HandleFunc("/article/{id}", returnSingleArticle)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func returnAllArticles(w http.ResponseWriter, r *http.Request){
	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(Articles)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	key := vars["id"]

	// Loop over all of the articles
	// if the article.Id equals the key we pass in
	// return the article encoded as JSON
	for _, article := range Articles{
		if article.Id == key {
			json.NewEncoder(w).Encode(article)
		}
	}
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	Articles=[]Article{
		Article{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
		Article{Id: "2", Title: "Hello2", Desc: "Article Description", Content: "Article Content"},
	}
	handleRequests()
}
