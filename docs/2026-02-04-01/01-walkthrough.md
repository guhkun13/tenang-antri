# Customer Ticket Tracking - Implementation Walkthrough

## Overview

Successfully implemented a customer-facing ticket tracking feature for the AntriQ queue management system. Customers can now track their queue status in real-time using their ticket number at the `/track` URL.

## Changes Made

### Backend Implementation

#### 1. Data Transfer Object

**New file:** [tracking.go](file:///home/guhkun/work/mine/antriQ/v7/queue-system/internal/dto/tracking.go)

Created `TrackingInfo` DTO to structure tracking response data with fields for ticket number, category, status, queue position, estimated wait time, and counter information.

---

#### 2. Database Query Layer

**Modified:** [ticket_queries.go](file:///home/guhkun/work/mine/antriQ/v7/queue-system/internal/query/ticket_queries.go#L59-L71)

Added `GetTicketByNumber()` query method to retrieve tickets by their ticket number string (e.g., "A001", "B005") with full category and counter details via SQL joins.

---

#### 3. Repository Layer

**Modified:** [ticket_repository.go](file:///home/guhkun/work/mine/antriQ/v7/queue-system/internal/repository/ticket_repository.go#L79-L117)

Added `GetByTicketNumber()` repository method that calls the query layer and maps results to the `Ticket` model with complete category and counter information.

---

#### 4. Service Layer

**New file:** [tracking_service.go](file:///home/guhkun/work/mine/antriQ/v7/queue-system/internal/service/tracking_service.go)

Created `TrackingService` with three core methods:

- **`GetTicketTrackingInfo()`** - Retrieves complete tracking information for a ticket
- **`CalculateQueuePosition()`** - Calculates position in queue by counting waiting tickets with same category created before the current ticket
- **`EstimateWaitTime()`** - Estimates wait time using formula: `position × 5 minutes` (configurable average service time)

---

#### 5. Handler Layer

**New file:** [tracking_handler.go](file:///home/guhkun/work/mine/antriQ/v7/queue-system/internal/handler/tracking_handler.go)

Created `TrackingHandler` with two endpoints:

- **`ShowTrackingPage()`** - Renders main tracking page at `/track`
- **`GetTrackingInfo()`** - Returns HTMX partial with tracking data at `/track/info/:ticket_number`

---

#### 6. Router Configuration

**Modified:** [router.go](file:///home/guhkun/work/mine/antriQ/v7/queue-system/internal/server/router.go#L37-L53)

- Initialized `TrackingService` and `TrackingHandler`
- Added public tracking routes group:
  - `GET /track/` → Main tracking page
  - `GET /track/info/:ticket_number` → HTMX endpoint for tracking info

---

### Frontend Implementation

#### 1. Main Tracking Page

**New file:** [index.html](file:///home/guhkun/work/mine/antriQ/v7/queue-system/web/templates/pages/track/index.html)

Features:
- Clean purple gradient design matching existing kiosk aesthetics
- Live clock display
- Informational card explaining how to use the feature
- Prominent ticket number input with uppercase formatting
- HTMX integration for dynamic updates
- JavaScript form handler that sets up auto-refresh on submit

---

#### 2. Tracking Info Partial

**New file:** [_tracking_info.html](file:///home/guhkun/work/mine/antriQ/v7/queue-system/web/templates/pages/track/_tracking_info.html)

Dynamic display based on ticket status:

- **Waiting**: Shows queue position and estimated wait time in large, easy-to-read cards
- **Serving**: Displays animated bell icon with counter number and name
- **Completed**: Shows completion message
- **Cancelled/No Show**: Shows appropriate status message
- **Error**: Displays user-friendly error message for invalid tickets

Includes auto-refresh indicator showing "Auto-updating every 10 seconds"

---

## Features Implemented

### ✅ Ticket Lookup

Customers can search for their ticket by entering the ticket number (e.g., A001, B005).

### ✅ Queue Position Display

Shows the customer's exact position in the queue, calculated based on tickets ahead of them in the same category.

### ✅ Wait Time Estimation

Provides estimated wait time in minutes using the formula: `position × 5 minutes per ticket`.

### ✅ Counter Assignment

When a ticket is being served, displays the counter number and name where the customer should go.

### ✅ Live Updates

HTMX polling automatically refreshes tracking information every 10 seconds without requiring manual page refresh.

### ✅ Status-Based Display

Different visual presentations for each ticket status (waiting, serving, completed, cancelled, no_show).

### ✅ Error Handling

Graceful error handling for:
- Invalid ticket numbers
- Non-existent tickets
- Empty input

### ✅ Responsive Design

Mobile-friendly layout that adapts to different screen sizes.

---

## Verification & Testing

### Test 1: Valid Ticket Tracking

**Steps:**
1. Generated ticket B005 from kiosk for Bakso service
2. Navigated to `/track` and entered ticket number
3. Clicked "Track" button

**Results:**
- ✅ Ticket information displayed correctly
- ✅ Status: "Waiting in Queue"
- ✅ Queue position: 3
- ✅ Estimated wait time: 15 minutes
- ✅ Service category shown with correct color (Bakso - orange/red gradient)

![Ticket Tracking](file:///home/guhkun/.gemini/antigravity/brain/f97a4537-07ab-45b8-9e01-8359b777fb62/tracking_info_initial_1770175373745.png)

---

### Test 2: Live Auto-Refresh

**Steps:**
1. Kept tracking page open for ticket B005
2. Waited for 12+ seconds without interaction
3. Observed page updates

**Results:**
- ✅ Page automatically refreshed
- ✅ Clock time updated from 11:22 AM to 11:23 AM
- ✅ Tracking information remained current
- ✅ No manual refresh required

**Before (11:22 AM):**

![Before Refresh](file:///home/guhkun/.gemini/antigravity/brain/f97a4537-07ab-45b8-9e01-8359b777fb62/tracking_info_initial_1770175373745.png)

**After (11:23 AM):**

![After Refresh](file:///home/guhkun/.gemini/antigravity/brain/f97a4537-07ab-45b8-9e01-8359b777fb62/tracking_info_after_wait_1770175400996.png)

---

### Test 3: Error Handling

**Steps:**
1. Entered invalid ticket number "INVALID999"
2. Clicked "Track" button

**Results:**
- ✅ Error message displayed: "Ticket Not Found"
- ✅ Helpful description: "Ticket not found. Please check your ticket number and try again."
- ✅ Clean error UI with warning icon
- ✅ No application crash or exceptions

![Error Handling](file:///home/guhkun/.gemini/antigravity/brain/f97a4537-07ab-45b8-9e01-8359b777fb62/invalid_ticket_error_1770175451607.png)

---

### Test 4: Build Verification

**Command:** `go build -o queue-server cmd/server/main.go`

**Result:** ✅ Build successful with no compilation errors

**Server logs:**
```
[GIN-debug] GET /track/ --> trackingHandler.ShowTrackingPage-fm
[GIN-debug] GET /track/info/:ticket_number --> trackingHandler.GetTrackingInfo-fm
```

---

## Browser Recording

Interactive demonstration of the tracking feature:

![Ticket Tracking Demo](file:///home/guhkun/.gemini/antigravity/brain/f97a4537-07ab-45b8-9e01-8359b777fb62/track_ticket_test_1770174188153.webp)

---

## Technical Highlights

### Queue Position Calculation

The queue position is calculated by counting tickets with:
- Same category ID
- Status = "waiting"
- Created earlier than the current ticket

This ensures accurate positioning regardless of ticket number sequence.

### Wait Time Algorithm

Estimated wait time = `queue_position × 5 minutes`

The 5-minute average service time is configurable and can be enhanced in the future to:
- Calculate actual average from historical data
- Vary by category
- Adjust based on time of day

### HTMX Integration

The implementation leverages HTMX for seamless live updates:
- Initial search triggered by form submission
- Auto-refresh setup via JavaScript
- HTMX polls `/track/info/:ticket_number` every 10 seconds
- Partial template swapping for efficient updates

---

## API Endpoints

### `GET /track`
**Description:** Main tracking page  
**Authentication:** None (public)  
**Response:** HTML page with ticket input form

### `GET /track/info/:ticket_number`
**Description:** Get tracking information for a ticket  
**Authentication:** None (public)  
**Parameters:** 
- `ticket_number` (URL param) - Ticket number to track (e.g., "A001")

**Response:** HTMX partial HTML with tracking info or error message

**Example requests logged:**
```
2026-02-04T11:22:46 GET /track/info/B005 status=200 latency=3.602726ms
2026-02-04T11:22:56 GET /track/info/B005 status=200 latency=2.527789ms
2026-02-04T11:24:04 GET /track/info/INVALID999 status=200 latency=0.721188ms (ticket not found)
```

---

## Files Created

- `internal/dto/tracking.go` - Tracking DTO
- `internal/service/tracking_service.go` - Business logic
- `internal/handler/tracking_handler.go` - HTTP handlers
- `web/templates/pages/track/index.html` - Main page
- `web/templates/pages/track/_tracking_info.html` - HTMX partial

## Files Modified

- `internal/query/ticket_queries.go` - Added GetTicketByNumber query
- `internal/repository/ticket_repository.go` - Added GetByTicketNumber method
- `internal/server/router.go` - Added tracking routes and service initialization

---

## Summary

The customer ticket tracking feature is **fully functional and tested**. Customers can now:
- Access the tracking page at `/track`
- Enter their ticket number
- View their real-time queue position and estimated wait time
- See which counter they'll be served at when called
- Receive automatic updates every 10 seconds without manual refresh

The implementation follows the existing codebase patterns, uses HTMX for live updates consistent with other pages, and provides a polished, user-friendly experience.
