package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	. "v1/Helper/Password"
	. "v1/Helper/SendResponse"
	. "v1/Helper/TokenCreation"
	. "v1/Helper/TokenValidation"
	. "v1/Models"
	. "v1/MongoConnection"
	. "v1/Redis"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var authdetails LoginDetail
	err := json.NewDecoder(r.Body).Decode(&authdetails)
	if err != nil {
		SendResponse(w, true, "010")
		return
	}
	var filter interface{}
	var authUser User
	filter = bson.M{"email": authdetails.Email}
	dbErr := Client.Database("formulas").Collection("User").FindOne(Ctx, filter).Decode(&authUser)
	if dbErr != nil {
		if dbErr == mongo.ErrNoDocuments {
			SendResponse(w, true, "030")
			return
		}
	}
	check := CheckPasswordHash(authdetails.Password, authUser.Password)
	if !check {
		SendResponse(w, true, "030")
		return
	}
	ts, err := CreateToken(authUser.Email)
	if err != nil {
		SendResponse(w, true, "006")
		return
	}
	err = CreateAuth(authUser.ID, ts)
	if err != nil {
		SendResponse(w, true, "300")
		return
	}
	tokens := map[string]interface{}{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
		"auth_User":     authUser.Username,
		"Email":         authUser.Email,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}
func Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.Header.Get("Authorization")
	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf") //this should be in an env file
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Signing Method not Matched")
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	if err != nil {
		SendResponse(w, true, "002")
		return
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		SendResponse(w, true, "004")
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			SendResponse(w, true, "004")
			return
		}
		userId := claims["user_id"].(string)
		deleted, delErr := DeleteAuth(refreshUuid)
		if delErr != nil || deleted == 0 {
			SendResponse(w, true, "003")
			return
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := CreateToken(userId)
		if createErr != nil {
			SendResponse(w, true, "006")
			return
		}
		//save the tokens metadata to redis
		var filter interface{}
		var authUser User
		filter = bson.M{"email": userId}
		dbErr := Client.Database("formulas").Collection("User").FindOne(Ctx, filter).Decode(&authUser)
		if dbErr != nil {
			if dbErr == mongo.ErrNoDocuments {
				SendResponse(w, true, "030")
				return
			}
		}
		err := CreateAuth(authUser.ID, ts)
		if err != nil {
			SendResponse(w, true, "300")
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
			"auth_User":     authUser.Username,
			"Email":         authUser.Email,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokens)
	} else {
		SendResponse(w, true, "002")
	}
}
func Logout(w http.ResponseWriter, r *http.Request) {
	metadata, err := ExtractTokenMetadata(r)
	if err != nil {
		SendResponse(w, true, "004")
		return
	}
	delErr := DeleteTokens(metadata)
	if delErr != nil {
		SendResponse(w, true, "400")
		return
	}
	SendResponse(w, false, "1001")
}
