package dto

import "tenangantri/internal/model"

// StaffDashboardResponse represents the staff dashboard data
type StaffDashboardResponse struct {
	User             *model.User          `json:"user"`
	Counter          *model.Counter       `json:"counter"`
	CurrentTicket    *model.Ticket        `json:"current_ticket"`
	WaitingTickets   []model.Ticket       `json:"waiting_tickets"`
	QueueStats       []CategoryQueueStats `json:"queue_stats"`
	CompletedTickets []model.Ticket       `json:"completed_tickets"`
	CategoryIDs      []int                `json:"category_ids"`
}

// StaffQueueStatusResponse represents the queue status for staff
type StaffQueueStatusResponse struct {
	Counter        *model.Counter       `json:"counter"`
	CurrentTicket  *model.Ticket        `json:"current_ticket"`
	WaitingTickets []model.Ticket       `json:"waiting_tickets"`
	QueueStats     []CategoryQueueStats `json:"queue_stats"`
}
