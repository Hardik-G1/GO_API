package readroutes

import (
	"encoding/json"
	"net/http"
	. "v1/Helper/SendResponse"
	. "v1/Models"
	. "v1/MongoConnection"

	"go.mongodb.org/mongo-driver/bson"
)

func SearchDataAll(w http.ResponseWriter, r *http.Request) {
	var form map[string]string
	json.NewDecoder(r.Body).Decode(&form)
	value, exists := form["search"]
	if exists == false {
		SendResponse(w, true, "010")
		return
	}
	var results []*Formula
	var filter interface{}
	filter = bson.M{"private": false, "name": bson.M{"$regex": value}}
	data, err := Client.Database("formulas").Collection("data").Find(Ctx, filter)

	if err != nil {
		SendResponse(w, true, "100")
		return
	}
	for data.Next(Ctx) {
		var elem Formula
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
