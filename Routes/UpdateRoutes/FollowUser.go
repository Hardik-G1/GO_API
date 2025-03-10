package updateroutes

import (
	"encoding/json"
	"net/http"
	. "v1/Helper/SendResponse"
	. "v1/MongoConnection"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FollowUser(w http.ResponseWriter, r *http.Request) {
	LoggedInUser := r.Header.Get("user")
	//get the parameter
	user := mux.Vars(r)["usertofollow"]
	set := mux.Vars(r)["set"]
	if LoggedInUser == "" {
		SendResponse(w, true, "008")
		return
	}
	pip1 := bson.D{
		{
			Key: "$match",
			Value: bson.M{
				"username": user,
			},
		},
	}

	data, err1 := Client.Database("formulas").Collection("User").Aggregate(Ctx, mongo.Pipeline{pip1})
	var showsLoaded []bson.M
	if err1 = data.All(Ctx, &showsLoaded); err1 != nil {
		SendResponse(w, true, "100")
		return
	}

	UserId := showsLoaded[0]["id"].(primitive.ObjectID).Hex()
	var filter interface{}
	filter = bson.M{"user": LoggedInUser}
	update1 := bson.D{
		{Key: "$addToSet",
			Value: bson.D{
				{Key: "follows", Value: UserId},
			},
		},
	}
	update2 := bson.D{
		{Key: "$pull",
			Value: bson.D{
				{Key: "follows", Value: UserId},
			},
		},
	}
	update := update1
	if set == "0" {
		update = update2
	}
	ok, err2 := Client.Database("formulas").Collection("User").UpdateOne(Ctx, filter, update)
	if err2 != nil {
		SendResponse(w, true, "100")
		return
	}
	json.NewEncoder(w).Encode(ok)
	return
}
