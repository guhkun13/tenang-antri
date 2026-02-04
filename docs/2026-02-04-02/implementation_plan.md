# Implementation Plan - Enhance Tracking Ticket Feature

The goal is to provide users with more detailed information when tracking their tickets, including counter status and the last ticket called in their category.

## Proposed Changes

### [Backend] Tracking DTO
#### [MODIFY] [tracking.go](file:///home/guhkun/work/mine/antih/v7/source/internal/dto/tracking.go)
- Add fields to [TrackingInfo](file:///home/guhkun/work/mine/antih/v7/source/internal/dto/tracking.go#6-17):
    - `CounterStatus` (string)
    - `IsCounterServing` (bool)
    - `CounterCurrentServingTicket` (string)
    - `LastCalledTicketNumber` (string)

### [Backend] Tracking Service
#### [MODIFY] [tracking_service.go](file:///home/guhkun/work/mine/antih/v7/source/internal/service/tracking_service.go)
- Enhance [GetTicketTrackingInfo](file:///home/guhkun/work/mine/antih/v7/source/internal/service/tracking_service.go#34-86):
    - If a counter is assigned to the ticket, fetch its status.
    - Check if the counter is currently serving another ticket.
    - Fetch the last called ticket number for the ticket's category.

### [Backend] Ticket Repository
#### [MODIFY] [ticket_repository.go](file:///home/guhkun/work/mine/antih/v7/source/internal/repository/ticket_repository.go)
- Add `GetLastCalledByCategoryID(ctx, categoryID)` method to retrieve the most recently called ticket number.

### [Backend] Ticket Queries
#### [MODIFY] [ticket_queries.go](file:///home/guhkun/work/mine/antih/v7/source/internal/query/ticket_queries.go)
- Add `GetLastCalledTicketByCategory(ctx, categoryID)` SQL query.

### [Frontend] Tracking Templates
#### [MODIFY] [_tracking_info.html](file:///home/guhkun/work/mine/antih/v7/source/web/templates/pages/track/_tracking_info.html)
- Update UI to show:
    - Counter active/inactive status.
    - If serving another ticket, show "Serving [Ticket Number]".
    - Show "Last Called: [Ticket Number]".
    - Rename "Your Position" to "Remaining Queue".
    - Remove "Estimated Wait Time".

## Verification Plan

### Manual Verification
1.  Open the tracking page.
2.  Create a ticket.
3.  Track the ticket:
    - Verify "Remaining Queue" is displayed.
    - Verify "Estimated Wait Time" is gone.
    - Verify "Last Called" shows correctly.
4.  As staff, call the ticket:
    - Verify counter status shows "Active" or similar on the tracking page.
    - Verify current serving ticket number is shown if applicable.
5.  Pause/Close counter:
    - Verify counter status updates on the tracking page.
