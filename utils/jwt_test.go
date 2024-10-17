package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWT(t *testing.T) {
	const testOpenID = "0F25t7h0T3"
	const testsecret = "tM6uI6gF"

	claims := jwt.MapClaims{}
	claims["open_id"] = testOpenID

	token, err := GenerateToken(claims, testsecret, time.Minute)
	if err != nil {
		t.Error(err)
		return
	}

	parsedClaims, err := ParseToken(token, testsecret)
	if err != nil {
		t.Error(err)
		return
	}

	parseOpenID, ok := parsedClaims["open_id"].(string)
	if !ok {
		t.Error("interface type assertions fail")
		return
	}

	if parsedClaims["open_id"] != parseOpenID {
		t.Error("parse JWT fail")
	}
}
