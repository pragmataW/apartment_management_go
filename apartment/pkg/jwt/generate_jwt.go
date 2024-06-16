package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

func (j jwtGenerator) GenerateJWT() (string, error) {
	claims := jwt.MapClaims{
		"flatNo": j.Claim.FlatNo,
		"role":   j.Claim.Role,
		"exp":    j.Claim.Exp,
		"email":  j.Claim.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenizedStr, err := token.SignedString([]byte(j.JwtKey))
	if err != nil {
		return "", err
	}

	return tokenizedStr, nil
}