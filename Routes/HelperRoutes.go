package routes

import (
	"net/http"
	. "v1/Helper/SendResponse"
	. "v1/MongoConnection"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func CheckUsernamePresent(w http.ResponseWriter, r *http.Request) {
	var filter interface{}
	filter = bson.M{"username": mux.Vars(r)["name"]}
	data, dberr := Client.Database("formulas").Collection("User").CountDocuments(Ctx, filter)
	if dberr != nil {
		SendResponse(w, true, "100")
		return
	}
	if data > 0 {
		SendResponse(w, true, "020")
		return
	}
	SendResponse(w, false, "1020")
}
func CheckEmailPresent(w http.ResponseWriter, r *http.Request) {
	var filter interface{}
	filter = bson.M{"email": mux.Vars(r)["mail"]}
	data, dberr := Client.Database("formulas").Collection("User").CountDocuments(Ctx, filter)
	if dberr != nil {
		SendResponse(w, true, "100")
		return
	}
	if data > 0 {
		SendResponse(w, true, "020")
		return
	}
	SendResponse(w, false, "1020")
}
func CheckUrlisAvailable(w http.ResponseWriter, r *http.Request) {
	url := mux.Vars(r)["url"]
	var filter interface{}
	filter = bson.M{"url": url}
	data, err := Client.Database("formulas").Collection("data").CountDocuments(Ctx, filter)
	if err != nil {
		SendResponse(w, true, "100")
		return
	}
	if data > 0 {
		SendResponse(w, true, "020")
		return
	}
	SendResponse(w, false, "1020")
}
