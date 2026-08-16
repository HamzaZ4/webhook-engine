package handler

import (
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/HamzaZ4/webhook-engine/internal/db"
)

func TestIdempotency(t *testing.T){
	database := db.Connect()
	defer database.Close()

	_, err := database.Exec("TRUNCATE TABLE webhook_events CASCADE")
	if err != nil {
		t.Fatalf("Failed to clean the db before the test : %v", err)
	}
	var wg sync.WaitGroup
	h := &WebhookHandler{DB: database}
	for range 500 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			body := strings.NewReader(`{"id":"pmntt_test_001","type":"payment.succeeded","data":{}}`)
    		req := httptest.NewRequest("POST", "/webhooks", body)
    		req.Header.Set("Content-Type", "application/json")
    		rr := httptest.NewRecorder()
    		h.Handle(rr, req)
		}()
	}
	wg.Wait()

	var count int
	database.QueryRow("SELECT COUNT(*) FROM webhook_events WHERE payment_event_id = $1", "pmntt_test_001").Scan(&count)
	if count != 1 {
    	t.Errorf("Expected 1 row, got %d", count)
	}
}