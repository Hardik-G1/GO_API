package updateroutes

import (
	"encoding/json"
	"net/http"
	"time"
	. "v1/Helper"
	. "v1/Helper/SendResponse"
	. "v1/Models"
	. "v1/MongoConnection"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func EditData(w http.ResponseWriter, r *http.Request) {
	LoggedInUser := r.Header.Get("user")
	//get the parameter
	url := mux.Vars(r)["url"]
	if LoggedInUser == "" {
		SendResponse(w, true, "008")
		return
	}
	_, id, formula := GetPrivate(url)
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
	if formula.User == authUser.Username {
		var form *Formula
		json.NewDecoder(r.Body).Decode(&form)
		form.ID = formula.ID
		form.CreatedAt = formula.CreatedAt
		form.LatestUse = formula.LatestUse
		form.Stars = formula.Stars
		form.Used = formula.Used
		form.Forked = formula.Forked
		form.ForkFrom = formula.ForkFrom
		form.Report = formula.Report
		form.User = LoggedInUser
		//url check
		newUrl := ""
		if form.Url != formula.Url && form.Url != "" {
			var filter interface{}
			filter = bson.M{"url": form.Url}
			data, err := Client.Database("formulas").Collection("data").CountDocuments(Ctx, filter)
			if err != nil {
				SendResponse(w, true, "100")
				return
			}
			if data > 0 {
				SendResponse(w, true, "020")
				return
			} else {
				newUrl = form.Url
			}
		} else {
			newUrl = formula.Url
		}
		newPrivate := formula.Private
		if form.Private != formula.Private {
			newPrivate = form.Private
		}
		newStructure := formula.Structure
		if len(form.Structure) != len(formula.Structure) {
			newStructure = form.Structure
		}
		form.ExpiresAt = time.Now().Add(time.Hour * 24 * 365 * 5)
		newForm := map[string]interface{}{
			"id":         form.ID,
			"name":       form.Name,
			"Structure":  newStructure,
			"user":       form.User,
			"private":    newPrivate,
			"url":        newUrl,
			"used":       form.Used,
			"created_at": form.CreatedAt,
			"latestUse":  form.LatestUse,
			"stars":      form.Stars,
			"forked":     form.Forked,
			"forkFrom":   form.ForkFrom,
			"report":     form.Report,
			"expiresAt":  form.ExpiresAt,
		}
		var filterUpdate interface{}
		filterUpdate = bson.M{"id": id}
		_ = Client.Database("formulas").Collection("data").FindOneAndReplace(Ctx, filterUpdate, newForm)
		SendResponse(w, false, "1100")
		return
	}
	SendResponse(w, true, "003")
	return
}
