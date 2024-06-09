package utils

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App
var AuthClient *auth.Client


func InitFirebase(){
	opt := option.WithCredentialsFile("config/firebase-service-account.json")
	var err error
	FirebaseApp, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}
	AuthClient, err = FirebaseApp.Auth(context.Background())
	if err != nil {
		log.Fatalf("error initalizing firebase auth client: %v\n", err)
	}
}

func ValidateToken(idToken string) (*auth.Token, error) {
	token, err := AuthClient.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}