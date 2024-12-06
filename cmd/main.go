package main

import (
	"log"
	"net/http"

	"github.com/senyabanana/library-service/internal/app"
)

func main() {
	application, err := app.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	application.Logger.Info("Server is running on :8080")
	application.Logger.WithError(err).Fatal(http.ListenAndServe(":8080", nil))
}
