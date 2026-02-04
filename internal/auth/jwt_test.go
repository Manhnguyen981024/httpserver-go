package auth

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJwt(t *testing.T) {
	secret := "abc123"
	userId := uuid.New()
	tokenString, err := MakeJWT(userId, secret, time.Minute*5)
	if err != nil {
		t.Fatalf("MakeJWT() error: %v", err)
	}

	log.Printf("Token is: %v", tokenString)

	returnedUserId, err := ValidateJWT(tokenString, secret)
	if err != nil {
		t.Fatalf("ValidateJWT() error: %v", err)
	}

	if returnedUserId != userId {
		t.Fatalf("expected userId %v, got %v", userId, returnedUserId)
	}

	t.Log("Test JWT")
}

func TestGetBearerHeader(t *testing.T) {
	header := &http.Header{}
	header.Add("Authorization", "Bearer abc")
	token, err := GetBearerToken(*header)
	if err != nil {
		t.Fatalf("GetBearerToken() error: %v", err)
	}
	if token != "abc" {
		t.Fatalf("expected token abc, got %v", token)
	}
	t.Log("Test Get Bearer Token OK")
}
