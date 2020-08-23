package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Employee struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name,omitempty"`
	Department string             `json:"department" bson:"department,omitempty"`
	Address    AddressInfo        `json:"address" bson:"address,omitempty"`
	Skills     []string           `json:"skills" bson:"skills,omitempty"`
	IsActive   bool               `json:"isActive" bson:"isActive,omitempty"`
}

type AddressInfo struct {
	HouseNumber int    `json:"houseNumber" bson:"houseNumber,omitempty"`
	Street      string `json:"street" bson:"street,omitempty"`
	City        string `json:"city" bson:"city,omitempty"`
	State       string `json:"state" bson:"state,omitempty"`
	Pincode     string `json:"pincode" bson:"pincode,omitempty"`
}
