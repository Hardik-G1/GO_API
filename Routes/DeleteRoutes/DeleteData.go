package deleteroutes

import (
	"encoding/json"
	"net/http"
	. "v1/Helper/SendResponse"
	. "v1/Models"
	. "v1/MongoConnection"
	. "v1/Redis"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func DeleteData(w http.ResponseWriter, r *http.Request) {

	LoggedInUser := r.Header.Get("user")
	//get the parameter
	var filter interface{}
	var authUser User
	filter = bson.M{"username": LoggedInUser}
	dbErr := Client.Database("formulas").Collection("User").FindOne(Ctx, filter).Decode(&authUser)
	if dbErr != nil {
		if dbErr == mongo.ErrNoDocuments {
			SendResponse(w, true, "010")
			return
		}
	}
	if LoggedInUser != "" {
		url := mux.Vars(r)["url"]
		var result *Formula
		var filter interface{}
		filter = bson.M{"url": url}
		data, err := Client.Database("formulas").Collection("data").Find(Ctx, filter)
		if err != nil {
			SendResponse(w, true, "100")
			return
		}
		for data.Next(Ctx) {
			err := data.Decode(&result)
			if err != nil {
				SendResponse(w, true, "040")
				return
			}
		}
		if result != nil {
			if LoggedInUser == authUser.ID.String() {
				DeleteDataRedis(url)
				deleted, err := Client.Database("formulas").Collection("data").DeleteOne(Ctx, filter)
				if err != nil {
					SendResponse(w, true, "100")
					return
				}
				json.NewEncoder(w).Encode(deleted)
				return
			} else {
				SendResponse(w, true, "003")
				return
			}
		} else {
			SendResponse(w, true, "1020")
			return
		}
	} else {
		SendResponse(w, true, "008")
		return
	}
}
