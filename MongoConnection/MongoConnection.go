package MongoConnection

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client     *mongo.Client
	Ctx        context.Context
	Cancel     context.CancelFunc
	Database   = "formulas"
	Col        = "data"
	Collection *mongo.Collection
)

func Connect() {
	uri := "{MONGO_URI}" \\ replace with your mongo uri
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	fmt.Println("connected")
	if err != nil {
		panic(err)
	}
	collection := client.Database(Database).Collection(Col)
	Client = client
	Ctx = ctx
	Collection = collection
	return
}

func Close(client *mongo.Client, ctx context.Context) {
	defer func() { ctx.Done() }()
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}
