package updateroutes

import (
	"encoding/json"
	"net/http"
	. "v1/Helper/SendResponse"
	. "v1/MongoConnection"

	"go.mongodb.org/mongo-driver/bson"
)

func EditUser(w http.ResponseWriter, r *http.Request) {
	LoggedInUser := r.Header.Get("user")
	//get the parameter
	if LoggedInUser == "" {
		SendResponse(w, true, "008")
		return
	}
	var form map[string]string
	json.NewDecoder(r.Body).Decode(&form)
	value, exists := form["mail"]
	if exists == false {
		SendResponse(w, true, "010")
		return
	}
	var filter interface{}
	filter = bson.M{"username": LoggedInUser}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "email", Value: value},
			},
		},
	}

	_ = Client.Database("formulas").Collection("User").FindOneAndUpdate(Ctx, filter, update)
	SendResponse(w, false, "1100")
	return
}
