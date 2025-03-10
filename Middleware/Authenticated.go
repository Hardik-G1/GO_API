package middleware

import (
	"net/http"
	. "v1/Helper/SendResponse"
	. "v1/Helper/TokenValidation"
	. "v1/MongoConnection"
	. "v1/Redis"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AuthenticateUser(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metadata, err := ExtractTokenMetadata(r)
		if err != nil {
			SendResponse(w, true, "004")
		} else {
			filter := bson.M{"email": metadata.UserId}
			opts := options.FindOne().SetProjection(bson.D{{Key: "id", Value: 1}})
			var userID interface{}
			err = Client.Database("formulas").Collection("User").FindOne(Ctx, filter, opts).Decode(&userID)
			if err != nil {
				SendResponse(w, true, "030")
				handler.ServeHTTP(w, r)
			}
			userId, err := FetchAuth(metadata, userID.(string))
			if err != nil {
				SendResponse(w, true, "002")
				handler.ServeHTTP(w, r)
			}
			r.Header.Set("user", userId)
		}
		handler.ServeHTTP(w, r)
	}
}
