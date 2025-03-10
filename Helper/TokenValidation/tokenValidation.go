package tokenvalidation

import (
	"fmt"
	"net/http"
	"os"
	. "v1/Models"

	"github.com/golang-jwt/jwt"
)

func ExtractToken(r *http.Request) string {
	access := r.Header.Get("Authorization")
	if access == "" {
		return ""
	}
	return access
}
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	accessToken := ExtractToken(r)
	if accessToken == "" {
		return nil, fmt.Errorf("No Token passed")
	}
	tokenSecret, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Signing Method not Matched")
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return tokenSecret, nil
}
func ExtractTokenMetadata(r *http.Request) (*RedisSession, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		Uuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId := claims["user_id"].(string)

		return &RedisSession{
			AccessUuid: Uuid,
			UserId:     userId,
		}, nil
	}
	return nil, err
}
