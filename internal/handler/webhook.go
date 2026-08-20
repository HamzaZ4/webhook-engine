package handler

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/HamzaZ4/webhook-engine/internal/ledger"
)

type WebhookHandler struct { 
	DB *sql.DB
}

type PaymentEvent struct{
	ID string `json:"id"`
	Type string `json:"type"`
	Data json.RawMessage `json:"data"`
}
func hashBody(body []byte) string {
	h := sha256.Sum256(body)
	return fmt.Sprintf("%x", h)
}

func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request){
	var event PaymentEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	payload, err := json.Marshal(event)
	if err != nil {
		http.Error(w, "failed to marshal payload", http.StatusInternalServerError)
		return
	}

	hashed_payload := hashBody(payload)
	var stored_hash string
	
	err = h.DB.QueryRow(`
		SELECT request_hash FROM webhook_events
		WHERE payment_event_id = $1	
	`, event.ID).Scan(&stored_hash)
	if err == nil {
		if stored_hash != hashed_payload {
			http.Error(w, "idempotency key reused with different body", http.StatusUnprocessableEntity)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}else if err != sql.ErrNoRows {
		log.Println("Failed to query event", err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	

	var webhookEventUUID string
	err = h.DB.QueryRow(`
		INSERT INTO webhook_events (payment_event_id, payload, status, request_hash)
		Values ($1, $2, 'received', $3)
		ON CONFLICT (payment_event_id) DO NOTHING
		RETURNING id
	`, event.ID, payload, hashed_payload ).Scan(&webhookEventUUID)

	if err != nil {
		log.Println("Failed to insert event:", err)
		http.Error(w, "failed to store event", http.StatusInternalServerError)
		return
	}

	if err := ledger.WriteLedgerEntries(h.DB, webhookEventUUID, 1000); err != nil {
		log.Println("Failed to write transaction to the ledger, errors : %w", err )
		http.Error(w, "Failed to write ledger", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}