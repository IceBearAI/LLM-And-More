package jwt

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = os.Getenv("JWT_KEY")

func init() {
	if jwtKey == "" {
		jwtKey = "9y7dfAlfJ%hJISD*"
	}
}

func JwtKeyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	} else {
		return []byte(GetJwtKey()), nil
	}
}

func GetJwtKey() string {
	return jwtKey
}
