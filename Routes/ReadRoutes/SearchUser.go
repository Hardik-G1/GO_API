package readroutes

import (
	"encoding/json"
	"net/http"
	. "v1/Helper/SendResponse"
	. "v1/Models"
	. "v1/MongoConnection"

	"go.mongodb.org/mongo-driver/bson"
)

func SearchUserAll(w http.ResponseWriter, r *http.Request) {
	var form map[string]string
	json.NewDecoder(r.Body).Decode(&form)
	value, exists := form["searchUser"]
	if exists == false {
		SendResponse(w, true, "010")
		return
	}
	var results []*User
	var filter interface{}
	filter = bson.M{"username": bson.M{"$regex": value}}
	data, err := Client.Database("formulas").Collection("User").Find(Ctx, filter)

	if err != nil {
		SendResponse(w, true, "100")
		return
	}
	for data.Next(Ctx) {
		var elem User
		err := data.Decode(&elem)
		if err != nil {
			SendResponse(w, true, "100")

			return
		}
		results = append(results, &elem)
	}
	jsonResponse, err := json.Marshal(results)
	if err != nil {
		SendResponse(w, true, "100")

		return
	}
	w.Write(jsonResponse)
}
