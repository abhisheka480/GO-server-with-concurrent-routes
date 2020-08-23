package models

import "github.com/dgrijalva/jwt-go"

//login credentials struct
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//Payload struct
type Claims struct {
	Username           string `json:"username"`
	jwt.StandardClaims        //expiration time of token
}
