package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Article structure
type Article struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title string             `json:"title" bson:"title,omitempty"`
	Body  string             `json:"body" bson:"body,omitempty"`
	Tags  string             `json:"tags" bson:"tags,omitempty"`
}

//Creating collection var as global variable for easy access across Functions
var collection = ConnecttoDB()

func main() {
	//Init Router
	router := httprouter.New()

	//Routing for different HTTP methods
	router.GET("/article", getArticles)
	router.GET("/article/:id", getArticle)
	router.GET("/articles/search?q=title", searchArticle)
	router.POST("/articles", createArticle)
	// set our port address as 8081
	log.Fatal(http.ListenAndServe(":8081", router))
}

// ConnecttoDB : unction to connect to mongoDB locally
func ConnecttoDB() *mongo.Collection {

	// Set client options
	//change the URI according to your database
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	//Error Handling
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	//DB collection address which we are going to use
	//available to functions of all scope
	collection := client.Database("Appointy").Collection("NewsArticles")

	return collection
}

//Function to get all Articles in DataBase
func getArticles(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// Article array
	var Articles []Article

	// bson.M{},  we passed empty filter of unordered map.
	cur, err := collection.Find(context.TODO(), bson.M{})

	//Error Handling
	if err != nil {
		log.Fatal(err)
	}

	// Close the cursor once finished
	defer cur.Close(context.TODO())

	//Loops over the cursor stream and appends to []Article array
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var article Article
		// decode similar to deserialize process.
		err := cur.Decode(&article)

		//Error Handling
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		Articles = append(Articles, article)
	}
	//Error Handling
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	//Encoding the data in Array to JSON format
	json.NewEncoder(w).Encode(Articles)
}

//Function to create a new Article in Database
func createArticle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	var article Article

	//Decoding the Data from JSON format to Article variable
	_ = json.NewDecoder(r.Body).Decode(&article)
	//inserts the data from decoded var to MongoDB in BSON format
	result, err := collection.InsertOne(context.TODO(), article)
	//Error Handling
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(result)
}

//Function to search Article by ID
func getArticle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	var article Article

	// string to primitive.ObjectID (typeCasting)
	id, _ := primitive.ObjectIDFromHex(ps.ByName("id"))

	// creating filter of unordered map with ID as input
	filter := bson.M{"_id": id}

	//Searching in DB with given ID as keyword
	err := collection.FindOne(context.TODO(), filter).Decode(&article)
	//Error Handling
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(article)
}

//Function to search Article by title
func searchArticle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var article Article

	//recovers the argument of search query present in URL after "q"
	title := string(r.URL.Query().Get("q"))

	//makes an unordered map filter of title
	filter := bson.M{"title": title}

	//Searching in DB with given title as keyword
	err := collection.FindOne(context.TODO(), filter).Decode(&article)
	//Error Handling
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(article)
}
