package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Member struct {
	Id     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"` //Id，mongdb auto-generated
	Token  string             `json:"token,omitempty" bson:"token"`      //Token
	Email  string             `json:"email,omitempty" bson:"email"`      //Email
	Domain string             `json:"domain,omitempty" bson:"domain"`
	Role   string             `json:"role,omitempty" bson:"role"`     //Role: 0-normal,1-validator,3-manager,4-admin
	Status string             `json:"status,omitempty" bson:"status"` //Status: 0-registered, 1-pending, 2-approved
	Name   string             `json:"name,omitempty" bson:"name"`     // Name: 别名
}
