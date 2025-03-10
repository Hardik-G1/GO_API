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

func StarData(w http.ResponseWriter, r *http.Request) {
	LoggedInUser := r.Header.Get("user")
	//get the parameter
	url := mux.Vars(r)["url"]
	set := mux.Vars(r)["set"]
	if LoggedInUser == "" {
		SendResponse(w, true, "008")
		return
	}
	a, b, _ := GetPrivate(url)
	if a == true {
		var filter interface{}
		filter = bson.M{"user": LoggedInUser}
		update1 := bson.D{
			{Key: "$addToSet",
				Value: bson.D{
					{Key: "starred", Value: b},
				},
			},
		}
		update2 := bson.D{
			{Key: "$pull",
				Value: bson.D{
					{Key: "starred", Value: b},
				},
			},
		}
		update := update1
		if set == "0" {
			update = update2
		}
		ok, err1 := Client.Database("formulas").Collection("User").UpdateOne(Ctx, filter, update)
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
