package updateroutes

import (
	"encoding/json"
	"net/http"
	. "v1/Helper"
	. "v1/Helper/SendResponse"
	. "v1/MongoConnection"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func ReportData(w http.ResponseWriter, r *http.Request) {
	LoggedInUser := r.Header.Get("user")
	//get the parameter
	url := mux.Vars(r)["url"]
	if LoggedInUser == "" {
		SendResponse(w, true, "008")
		return
	}
	a, b, _ := GetPrivate(url)
	if a != true {
		var filter interface{}
		filter = bson.M{"id": b}
		update := bson.D{
			{Key: "$addToSet",
				Value: bson.D{
					{Key: "report", Value: b},
				},
			},
		}
		ok, err1 := Client.Database("formulas").Collection("data").UpdateOne(Ctx, filter, update)
		if err1 != nil {
			SendResponse(w, true, "100")
			return
		}
		json.NewEncoder(w).Encode(ok)
		return
	} else {
		SendResponse(w, true, "003")
		return
	}
}
