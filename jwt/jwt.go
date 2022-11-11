package jwt

import (
	"errors"
	"helpers"
	"models"
	"strings"
	"time"

	db "db/users"

	jwt "github.com/dgrijalva/jwt-go"
)

/* Email Email used in all the endpoints */
var Email string

/* UserId It's the User Id that is going to be used in all the endpoints */
var UserId string

/* GenerateJWT Generates the encryption with JWT */
func GenerateJWT(user models.User) (string, error) {
	myKey := []byte(helpers.GetEnvVariable("JWT_SIGNING_KEY"))
	payload := jwt.MapClaims{
		"email":      user.Email,
		"name":       user.Name,
		"last_name":  user.LastName,
		"birth_date": user.BirthDate,
		"biography":  user.Biography,
		"location":   user.Location,
		"web_site":   user.WebSite,
		"_id":        user.Id.Hex(),
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokenStr, err := token.SignedString(myKey)

	if err != nil {
		return tokenStr, err
	}

	return tokenStr, nil
}

/* ProcessJWT process the JWT received in the request */
func ProcessJWT(token string) (*models.Claim, bool, string, error) {
	myKey := []byte(helpers.GetEnvVariable("JWT_SIGNING_KEY"))
	claims := &models.Claim{}
	splitToken := strings.Split(token, "Bearer")

	if len(splitToken) != 2 {
		return claims, false, "", errors.New("token format invalid")
	}

	token = strings.TrimSpace(splitToken[1])

	tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return myKey, nil
	})

	if !tkn.Valid {
		return claims, false, "", errors.New("invalid token")
	}

	if err != nil {
		return claims, false, "", err
	}

	_, isFound, id := db.IsUser(claims.Email)

	if !isFound {
		return claims, false, "", errors.New("user not found")
	}

	Email = claims.Email
	UserId = claims.Id.Hex()

	return claims, isFound, id, nil
}
