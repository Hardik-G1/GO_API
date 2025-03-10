package helper

import (
	"encoding/json"
	. "v1/Models"
	. "v1/MongoConnection"
	. "v1/Redis"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetPrivate(url string) (bool, primitive.ObjectID, *Formula) {
	redisData, _ := FetchFormula(url)
	if redisData == " " {
		var result *Formula
		var filter interface{}
		filter = bson.M{"url": url}
		data, err := Client.Database("formulas").Collection("data").Find(Ctx, filter)
		if err != nil {
			return false, primitive.NilObjectID, &Formula{}
		}
		for data.Next(Ctx) {
			err := data.Decode(&result)
			if err != nil {
				return false, primitive.NilObjectID, &Formula{}
			}
		}
		_, err = json.Marshal(result)
		if err != nil {
			return false, primitive.NilObjectID, &Formula{}
		}
		if result.Private {
			return true, result.ID, result
		} else {
			return false, primitive.NilObjectID, &Formula{}
		}
	} else {
		var result *Formula
		_ = json.Unmarshal([]byte(redisData), &result)
		if result.Private {
			return true, result.ID, result
		} else {
			return false, primitive.NilObjectID, &Formula{}
		}
	}
}
func IdinArray(a interface{}, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
