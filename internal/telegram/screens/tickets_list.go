package screens

import (
	"fmt"
	"log"
	"strings"

	"github.com/sergegamb/hobot/internal/managedesk"
	"github.com/sergegamb/hobot/internal/telegram/context"
	"github.com/sergegamb/hobot/internal/telegram/models"

	tele "gopkg.in/telebot.v4"
)

type RenderedScreen struct {
	Text   string
	Markup *tele.ReplyMarkup
}

// PaginationState manages pagination state
type PaginationState struct {
	CurrentPage int
	PageSize    int
	TotalItems  int
}

// TicketFilterState manages filter state
type TicketFilterState struct {
	SelectedFilter   string
	AvailableFilters []string
}

// TicketListState manages the overall state of the tickets list
type TicketListState struct {
	Pagination PaginationState
	Filter     TicketFilterState
}

type TicketsListScreen struct {
	Client       *managedesk.Client
	State        *TicketListState
	SessionStore *context.SessionStore
	UserID       int64
}

// NewTicketsListScreenWithSession creates a new tickets list screen with session state
// Loads saved user state (page and filter) from session store if available
func NewTicketsListScreenWithSession(
	client *managedesk.Client,
	sessionStore *context.SessionStore,
	userID int64,
) *TicketsListScreen {
	userState := sessionStore.GetUserState(userID)

	return &TicketsListScreen{
		Client: client,
		State: &TicketListState{
			Pagination: PaginationState{
				CurrentPage: userState.CurrentPage,
				PageSize:    10,
				TotalItems:  0,
			},
			Filter: TicketFilterState{
				SelectedFilter:   userState.SelectedFilter,
				AvailableFilters: []string{"all", "open", "closed", "pending"},
			},
		},
		SessionStore: sessionStore,
		UserID:       userID,
	}
}

// SetPage updates the current page in pagination state and persists to session
func (s *TicketsListScreen) SetPage(page int) {
	if page < 1 {
		page = 1
	}
	s.State.Pagination.CurrentPage = page

	// Persist to session if available
	if s.SessionStore != nil && s.UserID > 0 {
		s.SessionStore.UpdateUserPage(s.UserID, page)
	}
}

// SetFilter updates the selected filter and persists to session
func (s *TicketsListScreen) SetFilter(filter string) {
	s.State.Filter.SelectedFilter = filter

	// Persist to session if available (also resets page to 1)
	if s.SessionStore != nil && s.UserID > 0 {
		s.SessionStore.UpdateUserFilter(s.UserID, filter)
		s.State.Pagination.CurrentPage = 1 // Reset page when filter changes
	}
}

// GetPaginatedRange returns the start and end indices for the current page
func (s *TicketsListScreen) GetPaginatedRange() (start, end int) {
	start = (s.State.Pagination.CurrentPage - 1) * s.State.Pagination.PageSize
	end = start + s.State.Pagination.PageSize
	if end > s.State.Pagination.TotalItems {
		end = s.State.Pagination.TotalItems
	}
	return start, end
}

// Render builds and returns the tickets list screen
func (s *TicketsListScreen) Render() (
	*RenderedScreen,
	error,
) {
	// Create list info with pagination parameters
	listInfo := &managedesk.ListInfo{
		RowCount:      s.State.Pagination.PageSize,
		StartIndex:    (s.State.Pagination.CurrentPage - 1) * s.State.Pagination.PageSize,
		SortField:     "id",
		SortOrder:     "desc",
		GetTotalCount: true,
	}

	response, err := s.Client.GetRequestsWithListInfo(listInfo)
	if err != nil {
		return nil, err
	}

	requests := response.Requests

	// Update total items count from API response
	s.State.Pagination.TotalItems = response.ListInfo.TotalCount

	// Build message and keyboard
	text := s.buildListHeaderMessage(requests)
	markup := s.buildListKeyboard(requests)

	return &RenderedScreen{
		Text:   text,
		Markup: markup,
	}, nil
}

// buildListHeaderMessage builds the header message for tickets list with filter and pagination info
func (s *TicketsListScreen) buildListHeaderMessage(requests []managedesk.Request) string {
	start, end := s.GetPaginatedRange()

	var message strings.Builder

	message.WriteString("*Tickets*\n")
	message.WriteString(fmt.Sprintf("*Filter:* `%s`\n", s.State.Filter.SelectedFilter))
	message.WriteString(fmt.Sprintf("*Page:* %d | *Showing:* %d-%d of %d\n\n",
		s.State.Pagination.CurrentPage,
		start+1,
		end,
		s.State.Pagination.TotalItems,
	))

	return message.String()
}

// buildListKeyboard builds the inline keyboard for tickets list with filter, tickets, and pagination buttons
func (s *TicketsListScreen) buildListKeyboard(requests []managedesk.Request) *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}
	rows := []tele.Row{}

	// Add filter selection button at top
	filterBtn := markup.Data(
		fmt.Sprintf("🔍 Filter: %s", s.State.Filter.SelectedFilter),
		"filter_select",
	)
	rows = append(rows, markup.Row(filterBtn))

	// Add paginated tickets
	// The API already returns only the paginated items for the current page,
	// so iterate through all returned items (indices 0 to len(requests))
	for i := 0; i < len(requests); i++ {
		req := requests[i]
		buttonText := models.FormatListButtonText(&models.TicketListItem{
			ID: req.ID,
			// Subject: models.TruncateSubject(req.Subject, 40),
			Subject: models.AppendedSubject(req.Subject),
			Status:  models.TicketStatus(req.Status.Name),
		})

		btn := markup.Data(buttonText, "ticket_view", req.ID)
		rows = append(rows, markup.Row(btn))
	}

	// Add pagination buttons
	s.appendPaginationButtons(markup, &rows)

	markup.Inline(rows...)

	return markup
}

// appendPaginationButtons adds previous/next pagination buttons to the keyboard
func (s *TicketsListScreen) appendPaginationButtons(markup *tele.ReplyMarkup, rows *[]tele.Row) {
	log.Printf("[appendPaginationButtons] TotalItems: %d, PageSize: %d, CurrentPage: %d",
		s.State.Pagination.TotalItems, s.State.Pagination.PageSize, s.State.Pagination.CurrentPage)

	// Only show pagination if there are multiple pages or not on first page
	if s.State.Pagination.TotalItems <= s.State.Pagination.PageSize && s.State.Pagination.CurrentPage == 1 {
		log.Printf("[appendPaginationButtons] Only 1 page or less, skipping pagination buttons")
		return
	}

	var paginationRow tele.Row

	if s.State.Pagination.CurrentPage > 1 {
		log.Printf("[appendPaginationButtons] Adding PREVIOUS button")
		prevBtn := markup.Data(
			"⬅️ Previous",
			"page_prev",
		)
		paginationRow = append(paginationRow, prevBtn)
	}

	totalPages := (s.State.Pagination.TotalItems + s.State.Pagination.PageSize - 1) / s.State.Pagination.PageSize
	log.Printf("[appendPaginationButtons] Total pages: %d", totalPages)

	if s.State.Pagination.CurrentPage < totalPages {
		log.Printf("[appendPaginationButtons] Adding NEXT button (page_next)")
		nextBtn := markup.Data(
			"Next ➡️",
			"page_next",
		)
		paginationRow = append(paginationRow, nextBtn)
	}

	if len(paginationRow) > 0 {
		*rows = append(*rows, paginationRow)
		log.Printf("[appendPaginationButtons] Pagination buttons added successfully")
	}
}

// RenderDetailScreen builds the detail view for a single ticket
func (s *TicketsListScreen) RenderDetailScreen(request *managedesk.Request) (
	*RenderedScreen,
	error,
) {
	detail := &models.TicketDetail{
		Fields: models.TicketDisplayFields{
			ID:      request.ID,
			Subject: request.Subject,
			Status:  models.TicketStatus(request.Status.Name),
			// Priority:    request.Priority,
			Category: request.Category,
			// Created:     request.CreatedTime.Format("2006-01-02 15:04"),
			// Updated:     request.UpdatedTime.Format("2006-01-02 15:04"),
			Requester:   request.Requester.DisplayName,
			Description: request.Description,
		},
	}

	message := models.FormatDetailMessage(detail)
	markup := s.buildDetailKeyboard()

	return &RenderedScreen{
		Text:   message,
		Markup: markup,
	}, nil
}

// buildDetailKeyboard builds the inline keyboard for ticket detail view with back button
func (s *TicketsListScreen) buildDetailKeyboard() *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}

	backBtn := markup.Data("⬅️ Back to List", "tickets_list_back")
	markup.Inline(markup.Row(backBtn))

	return markup
}

// RenderFilterSelectionScreen builds the filter selection view
func (s *TicketsListScreen) RenderFilterSelectionScreen() *RenderedScreen {
	message := models.FormatFilterSelectionMessage()
	markup := s.buildFilterKeyboard()

	return &RenderedScreen{
		Text:   message,
		Markup: markup,
	}
}

// buildFilterKeyboard builds the inline keyboard for filter selection with all available filters
func (s *TicketsListScreen) buildFilterKeyboard() *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}

	filters := []struct {
		name  string
		value string
	}{
		{"🟢 Open", managedesk.FilterOpen},
		{"✅ Closed", managedesk.FilterClosed},
		{"⏳ Pending", managedesk.FilterPending},
		{"📋 All", managedesk.FilterAll},
	}

	var rows []tele.Row

	for _, f := range filters {
		btn := markup.Data(f.name, "filter_apply", f.value)
		rows = append(rows, markup.Row(btn))
	}

	backBtn := markup.Data("⬅️ Back", "tickets_list_back")
	rows = append(rows, markup.Row(backBtn))

	markup.Inline(rows...)

	return markup
}
