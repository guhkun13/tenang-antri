package dto

import "time"

// TrackingInfo contains comprehensive tracking information for a ticket
type TrackingInfo struct {
	TicketNumber                string    `json:"ticket_number"`
	CategoryName                string    `json:"category_name"`
	CategoryColor               string    `json:"category_color"`
	Status                      string    `json:"status"`
	QueuePosition               int       `json:"queue_position"`
	EstimatedWaitMin            int       `json:"estimated_wait_min"`
	CounterNumber               string    `json:"counter_number,omitempty"`
	CounterName                 string    `json:"counter_name,omitempty"`
	CounterStatus               string    `json:"counter_status,omitempty"`
	IsCounterServing            bool      `json:"is_counter_serving"`
	CounterCurrentServingTicket string    `json:"counter_current_serving_ticket,omitempty"`
	LastCalledTicketNumber      string    `json:"last_called_ticket_number,omitempty"`
	OperationalHours            string    `json:"operational_hours,omitempty"`
	CreatedAt                   time.Time `json:"created_at"`
}
