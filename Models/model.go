package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Formula struct {
	ID        primitive.ObjectID `bson:"id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Structure [][]string         `bson:"Structure" json:"Structure"`
	User      string             `bson:"user" json:"user"`
	Private   bool               `bson:"private" json:"private"`
	Url       string             `bson:"url" json:"url"`
	Used      int                `bson:"used" json:"used"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	LatestUse time.Time          `bson:"latestUse" json:"latestUse"`
	Stars     int                `bson:"stars" json:"stars"`
	Forked    bool               `bson:"forked" json:"forked"`
	ForkFrom  string             `bson:"forkFrom" json:"forkFrom"`
	Report    []string           `bson:"report" json:"report"`
	ExpiresAt time.Time          `bson:"expiresAt" json:"expiresAt"`
}
type User struct {
	ID          primitive.ObjectID `json:"id" bson:"id"`
	Email       string             `json:"email" validate:"required" bson:"email"`
	Password    string             `json:"password" validate:"required" bson:"password"`
	Username    string             `json:"username" bson:"username"`
	IsVerified  bool               `json:"isVerified" bson:"isVerified"`
	CreatedAt   time.Time          `json:"createdAt,omitempty" bson:"createdAt"`
	Starred     []string           `json:"starred" bson:"starred"`
	CreatedData []string           `json:"createdData" bson:"createdData"`
	Forked      []string           `json:"forked" bson:"forked"`
	Follows     []string           `json:"follows" bson:"follows"`
	FollowedBy  []string           `json:"followedBy" bson:"followedBy"`
}
type VerificationData struct {
	Email     string    `json:"email" validate:"required" bson:"email"`
	Code      string    `json:"code" validate:"required" bson:"code"`
	ExpiresAt time.Time `json:"expiresat" bson:"expiresat"`
}
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}
type LoginDetail struct {
	Email    string `json:"email" validate:"required" bson:"email"`
	Password string `json:"password" validate:"required" bson:"password"`
}
type RedisSession struct {
	AccessUuid string
	UserId     string
}
