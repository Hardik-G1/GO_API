package readroutes

import (
	"encoding/json"
	"net/http"
	"time"
	. "v1/Helper/SendResponse"
	. "v1/Models"
	. "v1/MongoConnection"
	. "v1/Redis"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateUsed(url string) {
	var filter interface{}
	filter = bson.M{"url": url}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "latestUse", Value: time.Now()},
			},
		},
		{Key: "$inc",
			Value: bson.D{
				{Key: "used", Value: 1},
			},
		},
	}
	_, _ = Client.Database("formulas").Collection("data").UpdateOne(Ctx, filter, update)
}
func GetData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	url := mux.Vars(r)["url"]
	redisData, _ := FetchFormula(url)
	if redisData == " " {
		var result *Formula
		var filter interface{}
		filter = bson.M{"url": url}
		data := Client.Database("formulas").Collection("data").FindOne(Ctx, filter)
		err := data.Decode(&result)
		if err != nil {
			SendResponse(w, true, "040")
			return
		}
		jsonResponse, err := json.Marshal(result)
		if err != nil {
			SendResponse(w, true, "040")
			return
		}
		if result.Private {
			LoggedInUser := r.Header.Get("user")
			if LoggedInUser == "" {
				SendResponse(w, true, "003")
				return
			} else if result.User == LoggedInUser {
				latest := result.LatestUse
				now := time.Now()
				diff := now.Sub(latest).Minutes()
				if diff >= 15 {
					UpdateUsed(result.Url)
				}
				SetFormula(result.Url, result, false, false)
				w.Write(jsonResponse)
			} else if result.User != LoggedInUser {
				SendResponse(w, true, "003")
				return
			}
		} else {
			latest := result.LatestUse
			now := time.Now()
			diff := now.Sub(latest).Minutes()
			if diff >= 15 {
				UpdateUsed(result.Url)
			}
			SetFormula(result.Url, result, false, false)
			w.Write(jsonResponse)
		}
	} else {
		var result *Formula
		_ = json.Unmarshal([]byte(redisData), &result)
		if result.Private {
			LoggedInUser := r.Header.Get("user")
			if LoggedInUser == "" {
				SendResponse(w, true, "003")
				return
			} else if result.User == LoggedInUser {
				latest := result.LatestUse
				now := time.Now()
				diff := now.Sub(latest).Minutes()
				if diff >= 15 {
					UpdateUsed(result.Url)
				}
				jsonResponse, err := json.Marshal(result)
				if err != nil {
					SendResponse(w, true, "040")
					return
				}
				w.Write(jsonResponse)
			} else if result.User != LoggedInUser {
				SendResponse(w, true, "003")
				return
			}
		} else {
			latest := result.LatestUse
			now := time.Now()
			diff := now.Sub(latest).Minutes()
			if diff >= 15 {
				UpdateUsed(result.Url)
			}
			jsonResponse, err := json.Marshal(result)
			if err != nil {
				SendResponse(w, true, "040")
				return
			}
			w.Write(jsonResponse)
		}
	}
}
