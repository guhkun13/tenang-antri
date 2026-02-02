package model

import (
	"time"
)

// DailyStats represents daily statistics
type DailyStats struct {
	Date             time.Time `json:"date" db:"date"`
	TotalTickets     int       `json:"total_tickets" db:"total_tickets"`
	CompletedTickets int       `json:"completed_tickets" db:"completed_tickets"`
	NoShowTickets    int       `json:"no_show_tickets" db:"no_show_tickets"`
	CancelledTickets int       `json:"cancelled_tickets" db:"cancelled_tickets"`
	AvgWaitTime      *int      `json:"avg_wait_time,omitempty" db:"avg_wait_time"`
	AvgServiceTime   *int      `json:"avg_service_time,omitempty" db:"avg_service_time"`
	PeakHour         *int      `json:"peak_hour,omitempty" db:"peak_hour"`
}

// DashboardStats represents dashboard statistics
type DashboardStats struct {
	TotalTicketsToday     int                  `json:"total_tickets_today"`
	CurrentlyServing      int                  `json:"currently_serving"`
	WaitingTickets        int                  `json:"waiting_tickets"`
	ActiveCounters        int                  `json:"active_counters"`
	PausedCounters        int                  `json:"paused_counters"`
	AvgWaitTime           int                  `json:"avg_wait_time"`
	AvgServiceTime        int                  `json:"avg_service_time"`
	TicketsByStatus       map[string]int       `json:"tickets_by_status"`
	QueueLengthByCategory []CategoryQueueStats `json:"queue_length_by_category"`
	HourlyDistribution    []HourlyStats        `json:"hourly_distribution"`
}

// CategoryQueueStats represents queue stats for a category
type CategoryQueueStats struct {
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"`
	Prefix       string `json:"prefix"`
	WaitingCount int    `json:"waiting_count"`
	ColorCode    string `json:"color_code"`
}

// HourlyStats represents hourly ticket statistics
type HourlyStats struct {
	Hour  int `json:"hour"`
	Count int `json:"count"`
}

// DisplayTicket represents a ticket for display board
type DisplayTicket struct {
	TicketNumber   string `json:"ticket_number"`
	CounterNumber  string `json:"counter_number"`
	CategoryPrefix string `json:"category_prefix"`
	ColorCode      string `json:"color_code"`
	Status         string `json:"status"`
}

// WebSocketMessage represents a WebSocket message
type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
