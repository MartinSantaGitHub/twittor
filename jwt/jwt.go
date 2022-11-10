package jwt

import (
	"helpers"
	"models"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

/* GenerateJWT generates the encryption with JWT */
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
