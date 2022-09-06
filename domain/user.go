package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserDAO struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Email          string             `json:"email"`
	Password       string             `json:"password"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
	ActivationCode string             `json:"activation_code" bson:"activation_code"`
	Activated      bool               `json:"activated"`
	RefreshToken   string             `json:"refresh_token" bson:"refresh_token"`
}
