package updateroutes

import (
	"net/http"
	"time"
	. "v1/Helper"
	. "v1/Helper/SendResponse"
	. "v1/MongoConnection"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ForkData(w http.ResponseWriter, r *http.Request) {
	LoggedInUser := r.Header.Get("user")
	//get the parameter
	url := mux.Vars(r)["url"]
	if LoggedInUser == "" {
		SendResponse(w, true, "008")
		return
	}
	a, _, result := GetPrivate(url)
	if a != true {
		result.Forked = true
		result.ForkFrom = result.User
		result.ID = primitive.NewObjectID()
		result.User = LoggedInUser
		result.Private = false
		result.Url = result.Url + "-" + LoggedInUser
		result.Used = 0
		result.CreatedAt = time.Now()
		result.LatestUse = time.Now()
		result.Stars = 0
		result.Report = []string{}
		_, err := Client.Database("formulas").Collection("data").InsertOne(Ctx, result)
		if err != nil {
			SendResponse(w, true, "100")
			return
		}
		var Userfilter interface{}
		Userfilter = bson.M{"username": LoggedInUser}
		updateData := bson.D{
			{Key: "$push",
				Value: bson.D{
					{Key: "createdData", Value: result.ID},
				},
			},
		}
		_, errUser := Client.Database("formulas").Collection("User").UpdateOne(Ctx, Userfilter, updateData)
		if errUser != nil {
			SendResponse(w, true, "100")
			return
		}
		SendResponse(w, false, "1100")
		return
	} else {
		SendResponse(w, true, "003")

		return
	}
}
