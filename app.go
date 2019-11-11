package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

func main() {
	http.HandleFunc("/", top)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func top(w http.ResponseWriter, r *http.Request) {
	// Firebase Auth クライアントの作成
	ctx := r.Context()
	client, err := newClient(ctx)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	log.Printf("HEADER %+v", r.Header)

	// ヘッダから ID Token の取得
	idToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	// ID Token の検証
	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Fatalf("error verifying ID token: %v\n", err)
	}

	w.Header().Set("Access-Control-Allow-Headers", "Authorization")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	fmt.Fprintf(w, `{"result":"verified id token. %v"}`, token)
}

func newClient(ctx context.Context) (*auth.Client, error) {
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, err
	}
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}
	return client, err
}
