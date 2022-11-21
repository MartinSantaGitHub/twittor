package request

import (
	jwt "github.com/dgrijalva/jwt-go"
)

/* Claim is the model to process the JWT */
type Claim struct {
	Email string `json:"email"`
	Id    string `json:"_id,omitempty"`
	jwt.StandardClaims
}
