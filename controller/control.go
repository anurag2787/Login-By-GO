package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"login/model"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://anurag2787:ayush2005@cluster0.bxvpni2.mongodb.net/userdata?retryWrites=true&w=majority&appName=Cluster0"
const dbName = "userdata"
const collectionName = "signinup"

var collection *mongo.Collection

func init() {
	clientOptions := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connected successfully")

	collection = client.Database(dbName).Collection(collectionName)

}

func loginuser(login model.Login) string {
	username := login.Username
	filterdata := bson.M{"username": username}
	var result model.Login
	err := collection.FindOne(context.Background(), filterdata).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// User not found or wrong credentials
			return "Username does not exist"
		}
		log.Fatal(err)
		return "Internal Error Occured"
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(login.Password))
	if err != nil {
		return "Incorrect Password"
	}

	fmt.Println("Login successful")
	return "Login Sucessfully"

}

func registeruser(register model.Register) string {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return "Internal Error Occured, User Not Registered"
	}
	register.Password = string(hashedPassword)

	inserted, err := collection.InsertOne(context.Background(), register)
	if err != nil {
		if mongoErr, ok := err.(mongo.WriteException); ok && len(mongoErr.WriteErrors) > 0 && mongoErr.WriteErrors[0].Code == 11000 {
			return "Username already exists"
		}
		log.Println("❌ Error inserting user:", err)
		return "Internal Error Occured, User Not Registered"
	}
	fmt.Println("User Registered Succesfully: ", inserted.InsertedID)
	id := inserted.InsertedID.(primitive.ObjectID).Hex()
	return id
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Origin", "POST")

	var register model.Register
	err := json.NewDecoder(r.Body).Decode(&register)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	id := registeruser(register)
	json.NewEncoder(w).Encode(map[string]string{
		"Message/User ID": id,
	})

}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Origin", "Get")

	var login model.Login
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	data := loginuser(login)

	// json.NewEncoder(w).Encode(data)
	json.NewEncoder(w).Encode(map[string]string{
		"message": data,
	})
}
