package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//List todos.
func List(responseWriter http.ResponseWriter, request *http.Request) {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:32768"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}

	collection := client.Database("gotodo").Collection("todo")

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}

	var todos []Todo

	cur, err := collection.Find(ctx, bson.D{})

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var todo Todo
		err := cur.Decode(&todo)

		if err != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
		}
		todos = append(todos, todo)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	bytes, err := json.Marshal(todos)

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write(bytes)
}

//Delete todo item by id.
func Delete(responseWriter http.ResponseWriter, request *http.Request) {
	var pathparameters = mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(pathparameters["todoID"])

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:32768"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}

	db := client.Database("gotodo")
	col := db.Collection("todo")

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}

	if err != nil {
		log.Fatal("primitive.ObjectIDFromHex ERROR:", err)
	} else {
		res, err := col.DeleteOne(ctx, bson.M{"_id": id})
		if err != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
		}

		if res.DeletedCount == 0 {
			responseWriter.WriteHeader(http.StatusInternalServerError)
		} else {
			responseWriter.WriteHeader(http.StatusOK)
		}
	}
}

//Toggle todo item by id.
func Toggle(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")

	var pathparameters = mux.Vars(request)

	//get id prom parameters
	id, _ := primitive.ObjectIDFromHex(pathparameters["todoID"])

	var todo Todo

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:32768"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}
	collection := client.Database("gotodo").Collection("todo")

	filter := bson.M{"_id": id}
	err = collection.FindOne(context.TODO(), filter).Decode(&todo)

	update := bson.D{
		{"$set", bson.D{
			{"isactive", !todo.IsActive},
		}},
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}

	todo.IsActive = !todo.IsActive

	json.NewEncoder(responseWriter).Encode(todo)
}

//Create todo item.
func Create(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")

	var todo Todo

	err := json.NewDecoder(request.Body).Decode(&todo)

	todo.IsActive = true

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:32768"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}

	collection := client.Database("gotodo").Collection("todo")

	result, err := collection.InsertOne(context.TODO(), todo)

	todo.ID = result.InsertedID.(primitive.ObjectID)

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(responseWriter).Encode(todo)
}

//GetOneTodo by id.
func GetOneTodo(responseWriter http.ResponseWriter, request *http.Request) {
	var result Todo
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:32768"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}

	collection := client.Database("gotodo").Collection("todo")

	err = collection.FindOne(context.TODO(), bson.D{}).Decode(&result)

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	responseWriter.Write([]byte(fmt.Sprintf(`{ "id": %d , "title" : %s  , "isactive" : %v}`, result.ID, result.Title, result.IsActive)))
}

func main() {
	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/todo/list", List).Methods(http.MethodGet)
	api.HandleFunc("/todo/create", Create).Methods(http.MethodPost)
	api.HandleFunc("/todo/toggle/{todoID}", Toggle).Methods(http.MethodPost)
	api.HandleFunc("/todo/delete/{todoID}", Delete).Methods(http.MethodDelete)
	api.HandleFunc("/todo/getonebyid/{todoID}", GetOneTodo).Methods(http.MethodGet)

	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5001", "http://localhost:3000"},
		AllowedMethods: []string{
			http.MethodGet, //http methods for your app
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},

		AllowedHeaders: []string{
			"*",
		},
	})

	log.Fatal(http.ListenAndServe(":5001", corsOpts.Handler(router)))
}
