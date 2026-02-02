package helper

import (
	"database/sql"

	"github.com/jackc/pgx/v5"

	"queue-system/internal/model"
)

// ScanUser scans a row into a User struct
func ScanUser(row pgx.Row) (*model.User, error) {
	user := &model.User{}
	err := row.Scan(
		&user.ID, &user.Username, &user.Password, &user.FullName,
		&user.Email, &user.Phone, &user.Role, &user.IsActive,
		&user.CounterID, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)
	return user, err
}

// ScanCategory scans a row into a Category struct
func ScanCategory(row pgx.Row) (*model.Category, error) {
	category := &model.Category{}
	err := row.Scan(
		&category.ID, &category.Name, &category.Prefix, &category.Priority,
		&category.ColorCode, &category.Description, &category.Icon,
		&category.IsActive, &category.CreatedAt, &category.UpdatedAt,
	)
	return category, err
}

// ScanCounter scans a row into a Counter struct
func ScanCounter(row pgx.Row) (*model.Counter, error) {
	counter := &model.Counter{}
	err := row.Scan(
		&counter.ID, &counter.Number, &counter.Name, &counter.Location,
		&counter.Status, &counter.IsActive, &counter.CurrentStaffID,
		&counter.CreatedAt, &counter.UpdatedAt,
	)
	return counter, err
}

// ScanTicket scans a row into a Ticket struct
func ScanTicket(row pgx.Row) (*model.Ticket, error) {
	ticket := &model.Ticket{}
	err := row.Scan(
		&ticket.ID, &ticket.TicketNumber, &ticket.CategoryID, &ticket.CounterID,
		&ticket.Status, &ticket.Priority, &ticket.CreatedAt, &ticket.CalledAt,
		&ticket.CompletedAt, &ticket.WaitTime, &ticket.ServiceTime, &ticket.Notes,
	)
	return ticket, err
}

// ScanTicketWithDetails scans a row into a Ticket struct with related data
func ScanTicketWithDetails(row pgx.Row) (*model.Ticket, error) {
	ticket := &model.Ticket{Category: &model.Category{}, Counter: &model.Counter{}}

	var catID *int
	var catName, catPrefix, catColor *string
	var coID *int
	var coNumber, coName *string

	err := row.Scan(
		&ticket.ID, &ticket.TicketNumber, &ticket.CategoryID, &ticket.CounterID,
		&ticket.Status, &ticket.Priority, &ticket.CreatedAt, &ticket.CalledAt,
		&ticket.CompletedAt, &ticket.WaitTime, &ticket.ServiceTime, &ticket.Notes,
		&catID, &catName, &catPrefix, &catColor,
		&coID, &coNumber, &coName,
	)
	if err != nil {
		return nil, err
	}

	if catID != nil {
		ticket.Category.ID = *catID
		ticket.Category.Name = *catName
		ticket.Category.Prefix = *catPrefix
		ticket.Category.ColorCode = *catColor
	}
	if coID != nil {
		ticket.Counter.ID = *coID
		ticket.Counter.Number = *coNumber
		ticket.Counter.Name = *coName
	}

	return ticket, nil
}

// ScanDisplayTicket scans a row into a DisplayTicket struct
func ScanDisplayTicket(row pgx.Row) (*model.DisplayTicket, error) {
	ticket := &model.DisplayTicket{}
	err := row.Scan(&ticket.TicketNumber, &ticket.CounterNumber, &ticket.CategoryPrefix, &ticket.ColorCode, &ticket.Status)
	return ticket, err
}

// ScanCategoryQueueStats scans a row into a CategoryQueueStats struct
func ScanCategoryQueueStats(row pgx.Row) (*model.CategoryQueueStats, error) {
	stats := &model.CategoryQueueStats{}
	err := row.Scan(&stats.CategoryID, &stats.CategoryName, &stats.Prefix, &stats.ColorCode, &stats.WaitingCount)
	return stats, err
}

// ScanHourlyStats scans a row into a HourlyStats struct
func ScanHourlyStats(row pgx.Row) (*model.HourlyStats, error) {
	stats := &model.HourlyStats{}
	err := row.Scan(&stats.Hour, &stats.Count)
	return stats, err
}

// PtrToInt converts sql.NullInt64 to *int
func PtrToInt(nullInt sql.NullInt64) *int {
	if nullInt.Valid {
		val := int(nullInt.Int64)
		return &val
	}
	return nil
}

// PtrToString converts sql.NullString to *string
func PtrToString(nullStr sql.NullString) *string {
	if nullStr.Valid {
		return &nullStr.String
	}
	return nil
}

// PtrToBool converts sql.NullBool to *bool
func PtrToBool(nullBool sql.NullBool) *bool {
	if nullBool.Valid {
		return &nullBool.Bool
	}
	return nil
}

// NullString creates a sql.NullString from string
func NullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

// NullInt64 creates a sql.NullInt64 from *int
func NullInt64(i *int) sql.NullInt64 {
	if i != nil {
		return sql.NullInt64{Int64: int64(*i), Valid: true}
	}
	return sql.NullInt64{Valid: false}
}
