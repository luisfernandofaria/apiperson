package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

const DATABASE = "senai"
const COLLECTION = "people"

type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Contact   *Contact           `json:"contact,omitempty"`
}

type Contact struct {
	Address *Address `json:"address,omitempty"`
	Phone   *Phone   `json:"phone,omitempty"`
}

type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

type Phone struct {
	Ddd    string `json:"ddd,omitempty"`
	Number string `json:"number,omitempty"`
}

func createPerson(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var person Person
	_ = json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database(DATABASE).Collection(COLLECTION)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}

func readPerson(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var people []Person
	collection := client.Database(DATABASE).Collection(COLLECTION)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	personID := mux.Vars(request)["_id"]
	if len(personID) == 0 {
		retrivePerson(ctx, collection, response, request)
	} else {
		retriveOnePerson(personID, response, request)
	}

	json.NewEncoder(response).Encode(people)
}

func retriveOnePerson(personID string, response http.ResponseWriter, request *http.Request) {

	id, _ := primitive.ObjectIDFromHex(personID)
	var person Person
	collection := client.Database(DATABASE).Collection(COLLECTION)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Person{ID: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
}

func retrivePerson(ctx context.Context, collection *mongo.Collection,
	response http.ResponseWriter, request *http.Request) {
	var people []Person
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
}

func main() {
	fmt.Println("Iniciando a aplicação..")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/person", createPerson).Methods("POST")
	router.HandleFunc("/person", readPerson).Methods("GET")
	router.HandleFunc("/person/{id}", readPerson).Methods("GET")

	router.HandleFunc("/person/{id}/contact", buscarContato).Methods("GET")
	router.HandleFunc("/person/{id}/contact", criarContato).Methods("POST")
	router.HandleFunc("/person/{id}/contact", alterarContato).Methods("PUT")
	router.HandleFunc("/person/{id}/contact", atualizarContado).Methods("PATCH")
	router.HandleFunc("/person/{id}/contact", deletarContato).Methods("DELETE")

	http.ListenAndServe(":5000", router)
}

func criarContato(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var person Person
	_ = json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database(DATABASE).Collection(COLLECTION)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}

func buscarContato(personID string, response http.ResponseWriter, request *http.Request) {

	id, _ := primitive.ObjectIDFromHex(personID)
	var person Person
	collection := client.Database(DATABASE).Collection(COLLECTION)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Person{ID: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
}

func retrivePerson(ctx context.Context, collection *mongo.Collection,
	response http.ResponseWriter, request *http.Request) {
	var people []Person
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
}
