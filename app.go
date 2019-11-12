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
	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", top)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	addr := fmt.Sprintf(":%s", port)

	if os.Getenv("ENABLE_CORS") != "" {
		c := cors.New(cors.Options{
			AllowedOrigins: []string{"http://localhost:4200"},
			AllowedHeaders: []string{"Authorization"},
			Debug:          true,
		})
		log.Fatal(http.ListenAndServe(addr, c.Handler(mux)))
	} else {
		log.Fatal(http.ListenAndServe(addr, mux))
	}
}

func top(w http.ResponseWriter, r *http.Request) {
	// Firebase Auth クライアントの作成
	ctx := r.Context()
	client, err := newClient(ctx)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// ヘッダから ID Token の取得
	idToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

	// ID Token の検証
	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Fatalf("error verifying ID token: %v\n", err)
	}

	fmt.Fprintf(w, `{"result":"verified id token. %+v"}`, token)
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
