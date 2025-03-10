package readroutes

import (
	"encoding/json"
	"net/http"
	. "v1/Helper"
	. "v1/Helper/SendResponse"
	. "v1/Models"
	. "v1/MongoConnection"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserPage(w http.ResponseWriter, r *http.Request) {
	//get the parameter
	user := mux.Vars(r)["user"]
	LoggedInUser := r.Header.Get("user")
	var filter interface{}
	var authUser User
	filter = bson.M{"username": user}
	dbErr := Client.Database("formulas").Collection("User").FindOne(Ctx, filter).Decode(&authUser)
	if dbErr != nil {
		if dbErr == mongo.ErrNoDocuments {
			SendResponse(w, true, "010")
			return
		}
	}
	getPrivate := false
	if authUser.ID.String() == LoggedInUser {
		getPrivate = true
	}
	pip1 := bson.D{
		{
			Key: "$match",
			Value: bson.M{
				"username": user,
			},
		},
	}
	pip2 := bson.D{
		{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "data"},
				{Key: "localField", Value: "createdData"},
				{Key: "foreignField", Value: "id"},
				{Key: "as", Value: "dataOwned"},
			},
		},
	}
	pip3 := bson.D{
		{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "data"},
				{Key: "localField", Value: "starred"},
				{Key: "foreignField", Value: "id"},
				{Key: "as", Value: "dataStarred"},
			},
		},
	}
	pip4 := bson.D{
		{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "data"},
				{Key: "localField", Value: "forked"},
				{Key: "foreignField", Value: "id"},
				{Key: "as", Value: "dataForked"},
			},
		},
	}
	pip5 := bson.D{
		{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "User"},
				{Key: "localField", Value: "follows"},
				{Key: "foreignField", Value: "id"},
				{Key: "as", Value: "dataFollows"},
			},
		},
	}
	pip6 := bson.D{
		{
			Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "User"},
				{Key: "localField", Value: "followedBy"},
				{Key: "foreignField", Value: "id"},
				{Key: "as", Value: "dataFollowedBy"},
			},
		},
	}
	pip7 := bson.D{
		{
			Key: "$project",
			Value: bson.D{
				{Key: "createdData", Value: 1},
				{Key: "starred", Value: 1},
				{Key: "forked", Value: 1},
				{Key: "username", Value: 1},
				{Key: "email", Value: 1},
				{Key: "dataStarred", Value: 1},
				{Key: "dataForked", Value: 1},
				{Key: "dataFollows", Value: 1},
				{Key: "dataFollowedBy", Value: 1},
				{Key: "dataOwned",
					Value: bson.D{
						{
							Key: "$filter",
							Value: bson.D{
								{Key: "input", Value: "$dataOwned"},
								{Key: "as", Value: "data"},
								{Key: "cond", Value: bson.D{
									{Key: "$eq",
										Value: bson.A{
											"$$data.private",
											getPrivate,
										},
									},
								},
								},
							},
						},
					},
				},
			},
		},
	}

	data, err := Client.Database("formulas").Collection("User").Aggregate(Ctx, mongo.Pipeline{pip1, pip2, pip3, pip4, pip5, pip6, pip7})

	var showsLoaded []bson.M
	if err = data.All(Ctx, &showsLoaded); err != nil {
		SendResponse(w, true, "100")
		return
	}

	UserId := showsLoaded[0]["id"].(primitive.ObjectID).Hex()

	if LoggedInUser != "" && LoggedInUser != authUser.ID.String() {
		var result *User
		var filter interface{}
		filter = bson.M{"username": LoggedInUser}
		dataUser, err := Client.Database("formulas").Collection("User").Find(Ctx, filter)
		if err != nil {
			SendResponse(w, true, "100")
			return
		}
		for dataUser.Next(Ctx) {
			err := dataUser.Decode(&result)
			if err != nil {
				SendResponse(w, true, "100")
				return
			}
		}
		LoggedInUserFollowList := result.Follows
		showsLoaded[0]["OwnProfile"] = false
		showsLoaded[0]["LoginNeeded"] = false
		if IdinArray(UserId, LoggedInUserFollowList) {
			showsLoaded[0]["isFollowedByCurrentUser"] = true
		} else {
			showsLoaded[0]["isFollowedByCurrentUser"] = false
		}
		showsLoaded[0]["LoggedInuserData"] = result

		jsonResponse, err := json.Marshal(showsLoaded)
		if err != nil {
			SendResponse(w, true, "100")
			return
		}
		w.Write(jsonResponse)
	} else if LoggedInUser == "" {
		showsLoaded[0]["OwnProfile"] = false
		showsLoaded[0]["isFollowedByCurrentUser"] = false
		showsLoaded[0]["LoginNeeded"] = true
		jsonResponse, err := json.Marshal(showsLoaded)
		if err != nil {
			SendResponse(w, true, "100")
			return
		}
		w.Write(jsonResponse)
	} else if LoggedInUser == user {
		showsLoaded[0]["OwnProfile"] = true
		showsLoaded[0]["LoginNeeded"] = false
		showsLoaded[0]["isFollowedByCurrentUser"] = false
		jsonResponse, err := json.Marshal(showsLoaded)
		if err != nil {
			SendResponse(w, true, "100")
			return
		}
		w.Write(jsonResponse)
	}
}
