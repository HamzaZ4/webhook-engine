package ledger

import (
	"database/sql"
	"fmt"
)


func WriteLedgerEntries(db * sql.DB, eventId string, amount int64) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction with error: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
    INSERT INTO ledger_entries (account_id, webhook_event_id, amount, entry_type)
    VALUES ($1, $2, $3, 'debit')
	`, "00000000-0000-0000-0000-000000000001", eventId, -amount)
	if err != nil {
		return fmt.Errorf("failed to insert debit: %w", err)
	}

	_, err = tx.Exec(`
    INSERT INTO ledger_entries (account_id, webhook_event_id, amount, entry_type)
    VALUES ($1, $2, $3, 'credit')
	`, "00000000-0000-0000-0000-000000000002", eventId, amount)
	if err != nil {
		return fmt.Errorf("failed to insert debit: %w", err)
	}

	return tx.Commit() 
}