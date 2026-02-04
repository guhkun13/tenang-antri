# Customer Ticket Tracking Implementation

This implementation plan outlines the changes needed to add a customer-facing ticket tracking feature that allows customers to check their queue status in real-time.

## Background

Currently, the AntriQ queue system has interfaces for:
- **Kiosk** - where customers generate tickets
- **Admin** - for queue management
- **Staff** - for calling and serving tickets
- **Display** - for showing currently serving tickets

This feature adds a new **Track** interface where customers can:
1. Enter their ticket number
2. See their queue position
3. See estimated wait time
4. See which counter they'll be served at
5. Get live updates without manual refresh

## Proposed Changes

### Backend Layer

#### [MODIFY] [ticket_repository.go](file:///home/guhkun/work/mine/antriQ/v7/queue-system/internal/repository/ticket_repository.go)

Add new repository method to fetch ticket by ticket number:

```go
// GetByTicketNumber retrieves a ticket by its ticket number string
func (r *TicketRepository) GetByTicketNumber(ctx context.Context, ticketNumber string) (*model.Ticket, error)
```

This will query the database using the `ticket_number` field instead of ID.

---

#### [MODIFY] [ticket_queries.go](file:///home/guhkun/work/mine/antriQ/v7/queue-system/internal/query/ticket_queries.go)

Add SQL query method:

```go
func (q *TicketQueries) GetTicketByNumber(ctx context.Context, ticketNumber string) pgx.Row
```

Executes: `SELECT * FROM tickets WHERE ticket_number = $1`

---

#### [NEW] [tracking_service.go](file:///home/guhkun/work/mine/antriQ/v7/queue-system/internal/service/tracking_service.go)

Create a new service to handle tracking-specific business logic:

```go
type TrackingService struct {
    ticketRepo   *repository.TicketRepository
    categoryRepo *repository.CategoryRepository
    counterRepo  *repository.CounterRepository
}

// GetTicketTrackingInfo retrieves comprehensive tracking information
func (s *TrackingService) GetTicketTrackingInfo(ctx context.Context, ticketNumber string) (*TrackingInfo, error)

// CalculateQueuePosition determines position in queue
func (s *TrackingService) CalculateQueuePosition(ctx context.Context, ticket *model.Ticket) (int, error)

// EstimateWaitTime calculates estimated wait time in minutes
func (s *TrackingService) EstimateWaitTime(ctx context.Context, ticket *model.Ticket, position int) (int, error)
```

**Logic Details:**
- **Queue Position**: Count tickets with same category, status='waiting', and created_at < current ticket's created_at
- **Wait Time Estimation**: 
  - If ticket status is 'serving': 0 minutes (already being served)
  - If status is 'waiting': position × average_service_time_per_ticket (default: 5 minutes)
  - If status is 'completed': show completion message
  - If status is 'cancelled' or 'no_show': show appropriate message

---

#### [NEW] [tracking.go](file:///home/guhkun/work/mine/antriQ/v7/queue-system/internal/dto/tracking.go)

Create DTO for tracking response:

```go
type TrackingInfo struct {
    TicketNumber     string
    CategoryName     string
    CategoryColor    string
    Status           string
    QueuePosition    int
    EstimatedWaitMin int
    CounterNumber    string
    CounterName      string
    CreatedAt        time.Time
}
```

---

#### [NEW] [tracking_handler.go](file:///home/guhkun/work/mine/antriQ/v7/queue-system/internal/handler/tracking_handler.go)

Create handler for tracking routes:

```go
type TrackingHandler struct {
    trackingService *service.TrackingService
}

// ShowTrackingPage renders the tracking page
func (h *TrackingHandler) ShowTrackingPage(c *gin.Context)

// GetTrackingInfo returns tracking information for a ticket (HTMX endpoint)
func (h *TrackingHandler) GetTrackingInfo(c *gin.Context)
```

---

#### [MODIFY] [router.go](file:///home/guhkun/work/mine/antriQ/v7/queue-system/internal/server/router.go)

Add tracking routes in the public routes section (after kiosk routes):

```go
// Tracking routes (public)
track := r.Group("/track")
{
    track.GET("/", trackingHandler.ShowTrackingPage)
    track.GET("/info/:ticket_number", trackingHandler.GetTrackingInfo)
}
```

Initialize the tracking service and handler in the router setup.

---

### Frontend Layer

#### [NEW] [index.html](file:///home/guhkun/work/mine/antriQ/v7/queue-system/web/templates/pages/track/index.html)

Create main tracking page with:
- Input form for ticket number
- Display area for tracking information (initially hidden)
- Live update mechanism using HTMX polling
- Error handling for invalid tickets

**Features:**
- Clean, customer-friendly UI matching existing kiosk design
- Large, readable text for queue position and wait time
- Color-coded status indicators
- Auto-refresh every 10 seconds using HTMX `hx-trigger="every 10s"`
- Responsive design for mobile and desktop

---

#### [NEW] [_tracking_info.html](file:///home/guhkun/work/mine/antriQ/v7/queue-system/web/templates/pages/track/_tracking_info.html)

HTMX partial template showing:
- Ticket number with category color
- Current status badge
- Queue position (if waiting)
- Estimated wait time
- Counter assignment (if assigned/serving)
- Created timestamp

This partial will be swapped in by HTMX on both initial search and periodic updates.

---

#### [MODIFY] [style.css](file:///home/guhkun/work/mine/antriQ/v7/queue-system/web/static/css/style.css) (if exists)

Add CSS styles for:
- `.tracking-container` - main container
- `.tracking-card` - information display card
- `.queue-position` - large position indicator
- `.wait-time-badge` - estimated time display
- `.status-badge` - status indicators with color coding
- `.refresh-indicator` - subtle loading indicator during HTMX refresh

## Verification Plan

### Automated Tests

Currently, the codebase does not have a comprehensive test suite structure visible. New unit tests will be created:

#### Create Backend Tests

**Test file:** `internal/service/tracking_service_test.go`

```bash
go test -v ./internal/service -run TestTrackingService
```

Test cases:
- `TestGetTicketTrackingInfo_ValidTicket` - retrieves tracking info successfully
- `TestGetTicketTrackingInfo_InvalidTicket` - returns error for non-existent ticket
- `TestCalculateQueuePosition_WaitingTicket` - calculates correct position
- `TestEstimateWaitTime_MultiplePositions` - estimates wait time accurately

**Test file:** `internal/repository/ticket_repository_test.go` (extend existing)

```bash
go test -v ./internal/repository -run TestGetByTicketNumber
```

Test cases:
- `TestGetByTicketNumber_Found` - finds ticket by number
- `TestGetByTicketNumber_NotFound` - handles non-existent ticket

### Manual Verification

#### 1. Basic Tracking Flow

**Prerequisites:** Server running with database populated with sample tickets

**Steps:**
1. Start the server: `make run` or `go run cmd/server/main.go`
2. Navigate to `http://localhost:8080/track` in browser
3. Enter a valid ticket number from the kiosk (e.g., "A001")
4. Click "Track Ticket" button
5. **Expected:** Tracking information displays showing:
   - Ticket number with category color
   - Current status (Waiting/Serving/Completed)
   - Queue position if status is waiting
   - Estimated wait time
   - Counter assignment if assigned

#### 2. Live Updates

**Steps:**
1. Open tracking page with a valid waiting ticket
2. Leave page open for 30 seconds
3. In another tab/window, use the staff dashboard to call tickets
4. **Expected:** Tracking page automatically updates queue position and wait time without manual refresh

#### 3. Error Handling

**Steps:**
1. Navigate to `/track`
2. Enter invalid ticket number (e.g., "INVALID123")
3. **Expected:** Error message displays indicating ticket not found
4. Try entering empty string
5. **Expected:** Validation error prompts for ticket number

#### 4. Different Ticket States

Test tracking for tickets in each state:

| State | Expected Display |
|-------|-----------------|
| `waiting` | Shows position in queue + estimated wait time |
| `serving` | Shows "Now Being Served" + counter details |
| `completed` | Shows "Your visit is complete" message |
| `cancelled` | Shows "Ticket has been cancelled" message |
| `no_show` | Shows "Ticket marked as no-show" message |

#### 5. Responsive Design

**Steps:**
1. Open tracking page on desktop browser
2. Resize window to mobile dimensions (375px width)
3. **Expected:** Layout adapts properly, text remains readable
4. Test on actual mobile device if available

## Notes

- Using HTMX for live updates is consistent with the existing kiosk implementation
- 10-second polling interval balances responsiveness with server load
- No authentication required as this is a public-facing feature
- Ticket numbers are already unique in the system, safe to use as lookup keys
