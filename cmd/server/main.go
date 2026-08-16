package main

import (
	"log"
	"net/http"

	"github.com/HamzaZ4/webhook-engine/internal/db"
	"github.com/HamzaZ4/webhook-engine/internal/handler"
)


func main(){
	database := db.Connect()
	webhookHandler := &handler.WebhookHandler{DB: database}

	mux := http.NewServeMux()

	mux.HandleFunc("/webhooks", webhookHandler.Handle)

	log.Fatal(http.ListenAndServe(":8081", mux))

}