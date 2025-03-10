package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
	. "v1/Models"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var client *redis.Client

func init() {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	client = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
}
func CreateVerificationCode(v VerificationData) error {
	at := time.Unix(v.ExpiresAt.Unix(), 0) //converting Unix to UTC(to Time object)
	now := time.Now()
	errAccess := client.Set(context.Background(), v.Email, v.Code, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	return nil
}
func FetchVerification(Email string) (string, error) {
	userid, err := client.Get(context.Background(), Email).Result()
	if err != nil {
		return "anon", err
	}
	return userid, nil
}
func CreateAuth(userid primitive.ObjectID, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()
	useridstring := userid.String()

	errAccess := client.Set(context.Background(), td.AccessUuid, useridstring, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := client.Set(context.Background(), td.RefreshUuid, useridstring, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}
func FetchAuth(authD *RedisSession, data string) (string, error) {
	userid, err := client.Get(context.Background(), authD.AccessUuid).Result()
	if err != nil {
		return "anon", err
	}
	if data != userid {
		return "anon", errors.New("unauthorized")
	}
	return userid, nil
}
func FetchFormula(url string) (string, error) {
	fetchedformula, err := client.Get(context.Background(), url).Result()
	if err != nil {
		return " ", err
	}
	return fetchedformula, nil
}
func SetFormula(url string, data *Formula, new bool, user bool) error {
	json, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	var calculateTime int64
	if data.Used > 20 && data.Used < 100 {
		calculateTime = time.Now().Add(time.Hour * 24).Unix()
	} else if data.Used > 99 && data.Used < 250 {
		calculateTime = time.Now().Add(time.Hour * 24 * 7).Unix()
	} else if data.Used > 250 {
		calculateTime = time.Now().Add(time.Hour * 24 * 14).Unix()
	} else {
		calculateTime = time.Now().Add(time.Hour * 1).Unix()
	}
	if new == true {
		calculateTime = time.Now().Add(time.Hour * 24).Unix()
	}
	if user == true {
		calculateTime = time.Now().Add(time.Hour * 24 * 7).Unix()
	}
	at := time.Unix(calculateTime, 0) //converting Unix to UTC(to Time object)
	now := time.Now()
	errAccess := client.Set(context.Background(), url, json, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	return nil
}
func DeleteDataRedis(url string) (int64, error) {
	deleted, err := client.Del(context.Background(), url).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
func DeleteAuth(givenUuid string) (int64, error) {
	deleted, err := client.Del(context.Background(), givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
func DeleteTokens(authD *RedisSession) error {
	//get the refresh uuid
	refreshUuid := fmt.Sprintf("%s++%s", authD.AccessUuid, authD.UserId)
	//delete access token
	deletedAt, err := client.Del(context.Background(), authD.AccessUuid).Result()
	if err != nil {
		return err
	}
	//delete refresh token
	deletedRt, err := client.Del(context.Background(), refreshUuid).Result()
	if err != nil {
		return err
	}
	//When the record is deleted, the return value is 1
	if deletedAt != 1 || deletedRt != 1 {
		return errors.New("something went wrong")
	}
	return nil
}
