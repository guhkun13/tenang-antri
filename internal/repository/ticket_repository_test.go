package repository

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"

	"tenangantri/internal/model"
	"tenangantri/internal/query"
)

func TestTicketRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	ticketQry := query.NewTicketQueries()
	repo := &ticketRepository{
		pool:      mock,
		ticketQry: ticketQry,
	}

	ticketID := 1
	now := time.Now()
	queueDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	rows := pgxmock.NewRows([]string{"id", "ticket_number", "category_id", "counter_id", "status", "priority", "created_at", "called_at", "completed_at", "wait_time", "service_time", "daily_sequence", "queue_date", "notes"}).
		AddRow(ticketID, "A001", 1, nil, "waiting", 1, now, nil, nil, nil, nil, 1, queueDate, "test notes")

	expectedSQL := `SELECT t.id, t.ticket_number, t.category_id, t.counter_id, t.status, t.priority, t.created_at, t.called_at, t.completed_at, t.wait_time, t.service_time, t.daily_sequence, t.queue_date, t.notes FROM tickets t WHERE t.id = \$1`

	mock.ExpectQuery(expectedSQL).
		WithArgs(ticketID).
		WillReturnRows(rows)

	ctx := context.Background()
	ticket, err := repo.GetByID(ctx, ticketID)

	assert.NoError(t, err)
	assert.NotNil(t, ticket)
	assert.Equal(t, ticketID, ticket.ID)
	assert.Equal(t, "A001", ticket.TicketNumber)
	assert.Equal(t, 1, ticket.DailySequence)
	assert.True(t, queueDate.Equal(ticket.QueueDate))

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTicketRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	ticketQry := query.NewTicketQueries()
	repo := &ticketRepository{
		pool:      mock,
		ticketQry: ticketQry,
	}

	now := time.Now()
	queueDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	ticket := &model.Ticket{
		TicketNumber:  "A001",
		Category:      &model.Category{ID: 1},
		Status:        "waiting",
		Priority:      1,
		Notes:         "test notes",
		DailySequence: 1,
		QueueDate:     queueDate,
	}

	mock.ExpectQuery("INSERT INTO tickets").
		WithArgs(ticket.TicketNumber, ticket.Category.ID, ticket.Status, ticket.Priority, ticket.Notes, ticket.DailySequence, ticket.QueueDate).
		WillReturnRows(pgxmock.NewRows([]string{"id", "created_at"}).AddRow(1, now))

	ctx := context.Background()
	createdTicket, err := repo.Create(ctx, ticket)

	assert.NoError(t, err)
	assert.NotNil(t, createdTicket)
	assert.Equal(t, 1, createdTicket.ID)
	assert.NotEmpty(t, createdTicket.CreatedAt)

	assert.NoError(t, mock.ExpectationsWereMet())
}
