package dataBase

import (
	"context"
	"webServer/dataBaseLib/schema"

	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
	"webServer/common"
)

var (
	//globalCtx context.Context
	dataBase=Init()
	docDB = dataBase.Collection("doc")
	UserDB=dataBase.Collection("Users")
	RoutineDB=dataBase.Collection("Routines")
)

func Init() *mongo.Database {
	var uri = common.LocalGetEnv("mongoDBSecretURL")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	//defer func() {
	//	if err = client.Disconnect(ctx); err != nil {
	//		panic(err)
	//	}
	//}()
	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	//globalCtx = ctx
	fmt.Println("Successfully connected to dataBase.")
	quickStartDatabase := client.Database("quickStart")
	return quickStartDatabase
}

func CreateArray() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := docDB.InsertOne(ctx, schema.TestDataType{
		Id: primitive.NewObjectID(),
		Title: "Testing Purpose",
		Array: []string{"first","second"},
	})
	//result, err := docDB.InsertOne(ctx, bson.D{
	//	{"title", "The Polyglot Developer Podcast"},
	//	{"array", bson.A{"development", "programming", "coding"}},
	//})
	if err != nil {
		log.Fatal(err)
	}
	//result, err = docDB.InsertOne(ctx, )
	fmt.Println(result)
}


func ReadArray() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id, _ := primitive.ObjectIDFromHex("6131878a5da0f4ef05e8ea48")
	result := docDB.FindOne(ctx, bson.D{
		{"_id", id},
	})
	if result.Err() != nil {
		log.Fatal(result.Err())
	}

	var dataJson schema.TestDataType
	err := result.Decode(&dataJson)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dataJson)
}

func UpdateArray() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id, _ := primitive.ObjectIDFromHex("6131878a5da0f4ef05e8ea48")
	result := docDB.FindOne(ctx, bson.M{"_id": id})


	if result.Err() != nil {
		log.Fatal(result.Err())
	}

	var dataJson schema.TestDataType
	err := result.Decode(&dataJson)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dataJson)
	fmt.Println(dataJson.Array)
	dataJson.Array[1]="Masbah"

	//updateData,err := docDB.UpdateOne(ctx, bson.M{"_id": id},bson.D{{"$set",bson.D{{"array", dataJson.Array}}}})
	updateData,err := docDB.UpdateOne(ctx, bson.M{"_id": id},bson.M{"$set":bson.M{"array": dataJson.Array}})
	fmt.Println(updateData)
}
