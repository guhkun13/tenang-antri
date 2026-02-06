package query

import (
	"context"
	"strings"
	"testing"
)

func TestTicketQueries_GenerateTicketNumber(t *testing.T) {
	q := NewTicketQueries()
	ctx := context.Background()

	sql := q.GenerateTicketNumber(ctx)

	if !strings.Contains(sql, "daily_sequence") {
		t.Errorf("Expected SQL to contain 'daily_sequence', got: %s", sql)
	}
	if !strings.Contains(sql, "queue_date = CURRENT_DATE") {
		t.Errorf("Expected SQL to contain 'queue_date = CURRENT_DATE', got: %s", sql)
	}
}

func TestTicketQueries_CreateTicket(t *testing.T) {
	q := NewTicketQueries()
	ctx := context.Background()

	sql := q.CreateTicket(ctx)

	expectedColumns := []string{"daily_sequence", "queue_date"}
	for _, col := range expectedColumns {
		if !strings.Contains(sql, col) {
			t.Errorf("Expected SQL to contain column '%s', got: %s", col, sql)
		}
	}
}

func TestTicketQueries_GetTicketWithDetails(t *testing.T) {
	q := NewTicketQueries()
	ctx := context.Background()

	sql := q.GetTicketWithDetails(ctx)

	if !strings.Contains(sql, "t.daily_sequence") {
		t.Errorf("Expected SQL to contain 't.daily_sequence', got: %s", sql)
	}
	if !strings.Contains(sql, "t.queue_date") {
		t.Errorf("Expected SQL to contain 't.queue_date', got: %s", sql)
	}
}

func TestStatsQueries_GetCurrentlyServingTickets(t *testing.T) {
	q := NewStatsQueries()
	ctx := context.Background()

	sql := q.GetCurrentlyServingTickets(ctx)

	if !strings.Contains(sql, "t.daily_sequence") {
		t.Errorf("Expected SQL to contain 't.daily_sequence', got: %s", sql)
	}
	if !strings.Contains(sql, "t.queue_date") {
		t.Errorf("Expected SQL to contain 't.queue_date', got: %s", sql)
	}
}

func TestTicketQueries_ListTickets(t *testing.T) {
	q := NewTicketQueries()
	ctx := context.Background()

	// Test default query
	res := q.ListTickets(ctx, nil)
	if !strings.Contains(res.Query, "t.daily_sequence") || !strings.Contains(res.Query, "t.queue_date") {
		t.Errorf("Default ListTickets query missing daily_sequence/queue_date: %s", res.Query)
	}

	// Test filters - though ListTickets doesn't use queue_date for filtering yet based on my previous edits
	// (it used created_at for date_from/date_to). Let's check if I should update that too.
}
