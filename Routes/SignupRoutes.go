package routes

import (
	"encoding/json"
	"net/http"
	"time"
	. "v1/Helper/GenerateCode"
	. "v1/Helper/Mailer"
	. "v1/Helper/Password"
	. "v1/Helper/SendResponse"
	. "v1/Models"
	. "v1/MongoConnection"
	. "v1/Redis"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	LoggedInUser := r.Header.Get("user")
	if LoggedInUser == "" {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			SendResponse(w, true, "010")
			return
		}
		user.Password, _ = GeneratehashPassword(user.Password)
		user.CreatedAt = time.Now()
		newForm := map[string]interface{}{
			"id":          primitive.NewObjectID(),
			"email":       user.Email,
			"password":    user.Password,
			"username":    user.Username,
			"isVerified":  user.IsVerified,
			"createdAt":   user.CreatedAt,
			"starred":     user.Starred,
			"createdData": user.CreatedData,
			"forked":      user.Forked,
			"follows":     user.Follows,
			"followedBy":  user.FollowedBy,
		}
		var filter interface{}
		filter = bson.M{"username": user.Username}
		data, dberr := Client.Database("formulas").Collection("User").CountDocuments(Ctx, filter)
		if dberr != nil {
			SendResponse(w, true, "100")
			return
		}
		if data > 0 {
			SendResponse(w, true, "020")
			return
		}
		filter = bson.M{"email": mux.Vars(r)["mail"]}
		data, dberr = Client.Database("formulas").Collection("User").CountDocuments(Ctx, filter)
		if dberr != nil {
			SendResponse(w, true, "100")
			return
		}
		if data > 0 {
			SendResponse(w, true, "020")
			return
		}
		var verificationData VerificationData
		verificationData.Email = user.Email
		verificationData.Code = GenerateRandomString(6)
		verificationData.ExpiresAt = time.Now().Add(time.Minute * 10)
		redisErr := CreateVerificationCode(verificationData)
		if redisErr != nil {
			SendResponse(w, true, "200")
			return
		}
		_, dberr = Client.Database("formulas").Collection("User").InsertOne(Ctx, newForm)
		if dberr != nil {
			SendResponse(w, true, "100")
		}
		SendAccountConfirmation(verificationData.Code, verificationData.Email)
		SendResponse(w, false, "007")
	} else {
		SendResponse(w, false, "001")
	}
}
func VerifyMail(w http.ResponseWriter, r *http.Request) {
	var data VerificationData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		SendResponse(w, true, "010")
		return
	}
	code, redisErr := FetchVerification(data.Email)
	if redisErr != nil {
		SendResponse(w, true, "1200")
		return
	}
	if code == data.Code {
		var filter interface{}
		filter = bson.M{"email": data.Email}
		value := true
		update := bson.D{
			{Key: "$set",
				Value: bson.D{
					{Key: "isVerified", Value: value},
				},
			},
		}
		ok, dbErr := Client.Database("formulas").Collection("User").UpdateOne(Ctx, filter, update)
		if dbErr != nil {
			SendResponse(w, true, "100")
			return
		}
		if ok.ModifiedCount == 1 {
			SendResponse(w, false, "1005")
			return
		}
		SendResponse(w, true, "005")
	} else {
		SendResponse(w, true, "005")
	}
}
func ResendVerificationMail(w http.ResponseWriter, r *http.Request) {
	LoggedInUser := r.Header.Get("user")
	if LoggedInUser != "" {
		var verificationData VerificationData
		err := json.NewDecoder(r.Body).Decode(&verificationData)
		if err != nil {
			SendResponse(w, true, "010")
		}
		var filter interface{}
		var authUser User
		filter = bson.M{"email": verificationData.Email}

		dbErr := Client.Database("formulas").Collection("User").FindOne(Ctx, filter).Decode(&authUser)
		if dbErr != nil {
			if dbErr == mongo.ErrNoDocuments {
				SendResponse(w, true, "030")
			}
		}
		if authUser.ID.String() == LoggedInUser {
			verificationData.Code = GenerateRandomString(6)
			verificationData.ExpiresAt = time.Now().Add(time.Minute * 10)
			redisErr := CreateVerificationCode(verificationData)
			if redisErr != nil {
				SendResponse(w, true, "500")
			}
			SendResponse(w, false, "007")
			SendAccountConfirmation(verificationData.Code, verificationData.Email)
		} else {
			SendResponse(w, true, "003")
		}
	} else {
		SendResponse(w, true, "003")
	}
}
