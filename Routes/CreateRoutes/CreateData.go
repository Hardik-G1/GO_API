package createroutes

import (
	"encoding/json"
	"net/http"
	"time"
	. "v1/Helper/SendResponse"
	. "v1/Models"
	. "v1/MongoConnection"
	. "v1/Redis"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateData(w http.ResponseWriter, r *http.Request) {
	var form *Formula
	json.NewDecoder(r.Body).Decode(&form)
	form.ID = primitive.NewObjectID()
	form.CreatedAt = time.Now()
	form.LatestUse = time.Now()
	form.Stars = 0
	form.Used = 0
	form.Forked = false
	form.ForkFrom = "owner"
	form.Report = []string{}
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
	}
	LoggedInUser := r.Header.Get("user")
	var authUser User
	filter = bson.M{"username": form.User}
	dbErr := Client.Database("formulas").Collection("User").FindOne(Ctx, filter).Decode(&authUser)
	if dbErr != nil {
		if dbErr == mongo.ErrNoDocuments {
			SendResponse(w, true, "010")
			return
		}
	}
	if form.User == "" && LoggedInUser == "" {
		form.Private = false
		form.ExpiresAt = time.Now().Add(time.Hour * 24 * 14)
		_, err := Client.Database("formulas").Collection("data").InsertOne(Ctx, form)
		if err != nil {
			SendResponse(w, true, "050")
			return
		}
		err = SetFormula(form.Url, form, true, false)
		if err != nil {
			SendResponse(w, true, "500")
		}
	} else if authUser.ID.String() == LoggedInUser {
		form.ExpiresAt = time.Now().Add(time.Hour * 24 * 365 * 5)
		newForm := map[string]interface{}{
			"id":         form.ID,
			"name":       form.Name,
			"Structure":  form.Structure,
			"user":       form.User,
			"private":    form.Private,
			"url":        form.Url,
			"used":       form.Used,
			"created_at": form.CreatedAt,
			"latestUse":  form.LatestUse,
			"stars":      form.Stars,
			"forked":     form.Forked,
			"forkFrom":   form.ForkFrom,
			"report":     form.Report,
			"expiresAt":  form.ExpiresAt,
		}

		_, err := Client.Database("formulas").Collection("data").InsertOne(Ctx, newForm)
		if err != nil {
			SendResponse(w, true, "050")
			return
		}
		var Userfilter interface{}
		Userfilter = bson.M{"username": form.User}
		updateData := bson.D{
			{Key: "$push",
				Value: bson.D{
					{Key: "createdData", Value: form.ID},
				},
			},
		}
		_, errUser := Client.Database("formulas").Collection("User").UpdateOne(Ctx, Userfilter, updateData)
		if errUser != nil {
			var filter interface{}
			filter = bson.M{"id": form.ID}
			_ = Client.Database("formulas").Collection("data").FindOneAndDelete(Ctx, filter)
			SendResponse(w, true, "050")
			return
		}
		err = SetFormula(form.Url, form, true, true)
		if err != nil {
			SendResponse(w, true, "500")
		}
	} else if authUser.ID.String() != LoggedInUser {
		SendResponse(w, true, "003")
		return
	}
	SendResponse(w, false, "1100")
}
