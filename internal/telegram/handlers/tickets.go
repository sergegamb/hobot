package handlers

import (
	"log"

	"github.com/sergegamb/hobot/internal/managedesk"
	"github.com/sergegamb/hobot/internal/telegram/context"
	"github.com/sergegamb/hobot/internal/telegram/screens"

	tele "gopkg.in/telebot.v4"
)

func TicketsMenuHandler(
	client *managedesk.Client,
	sessionStore *context.SessionStore,
) tele.HandlerFunc {

	return func(c tele.Context) error {
		userID := int64(c.Sender().ID)
		log.Printf("[TicketsMenuHandler] Started - UserID: %d", userID)

		// Create screen with persisted user state
		screen := screens.NewTicketsListScreenWithSession(
			client,
			sessionStore,
			userID,
		)

		rendered, err := screen.Render()
		if err != nil {
			log.Printf("[TicketsMenuHandler] ERROR rendering tickets menu: %v", err)
			return err
		}

		log.Printf("[TicketsMenuHandler] Successfully rendered tickets list")
		return c.Send(
			rendered.Text,
			rendered.Markup,
			tele.ModeMarkdown,
		)
	}
}

// TicketViewHandler handles ticket selection from the list
func TicketViewHandler(
	client *managedesk.Client,
	sessionStore *context.SessionStore,
) tele.HandlerFunc {

	return func(c tele.Context) error {
		userID := int64(c.Sender().ID)

		// Extract ticket ID from callback data
		ticketID := c.Data()
		log.Printf("[TicketViewHandler] UserID: %d, TicketID: %s", userID, ticketID)

		// Fetch ticket details
		request, err := client.GetRequestByID(ticketID)
		if err != nil {
			log.Printf("[TicketViewHandler] ERROR fetching ticket %s: %v", ticketID, err)
			return c.Respond(&tele.CallbackResponse{
				Text: "Error loading ticket details",
				// Alert: true,
			})
		}

		// Use screen's detail renderer with session
		screen := screens.NewTicketsListScreenWithSession(
			client,
			sessionStore,
			userID,
		)
		rendered, err := screen.RenderDetailScreen(request)
		if err != nil {
			log.Printf("[TicketViewHandler] ERROR rendering detail screen for ticket %s: %v", ticketID, err)
			return c.Respond(&tele.CallbackResponse{
				Text: "Error rendering ticket details",
				// Alert: true,
			})
		}

		log.Printf("[TicketViewHandler] Successfully rendered detail screen for ticket %s", ticketID)
		// Edit the message with ticket details
		return c.Edit(
			rendered.Text,
			rendered.Markup,
			tele.ModeMarkdown,
		)
	}
}

// FilterSelectHandler shows available filters
func FilterSelectHandler(
	client *managedesk.Client,
	sessionStore *context.SessionStore,
) tele.HandlerFunc {

	return func(c tele.Context) error {
		userID := int64(c.Sender().ID)
		log.Printf("[FilterSelectHandler] UserID: %d", userID)

		// Use screen's filter selection renderer with session
		screen := screens.NewTicketsListScreenWithSession(
			client,
			sessionStore,
			userID,
		)
		rendered := screen.RenderFilterSelectionScreen()
		log.Printf("[FilterSelectHandler] Rendering filter selection screen")

		editErr := c.Edit(
			rendered.Text,
			rendered.Markup,
			tele.ModeMarkdown,
		)
		if editErr != nil {
			log.Printf("[FilterSelectHandler] ERROR editing message: %v", editErr)
			return editErr
		}

		log.Printf("[FilterSelectHandler] Successfully edited message")
		return nil
	}
}

// FilterApplyHandler applies selected filter and returns to list
func FilterApplyHandler(
	client *managedesk.Client,
	sessionStore *context.SessionStore,
) tele.HandlerFunc {

	return func(c tele.Context) error {
		userID := int64(c.Sender().ID)

		// Extract filter value from callback data
		filter := c.Data()
		log.Printf("[FilterApplyHandler] UserID: %d, Filter: %s", userID, filter)

		// Create screen with new filter and session
		screen := screens.NewTicketsListScreenWithSession(
			client,
			sessionStore,
			userID,
		)
		screen.SetFilter(filter)
		// SetFilter already resets page to 1 via session store

		rendered, err := screen.Render()
		if err != nil {
			log.Printf("[FilterApplyHandler] ERROR rendering with filter %s: %v", filter, err)
			return c.Respond(&tele.CallbackResponse{
				Text: "Error applying filter",
				// Alert: true,
			})
		}

		log.Printf("[FilterApplyHandler] Successfully rendered filter %s", filter)
		editErr := c.Edit(
			rendered.Text,
			rendered.Markup,
			tele.ModeMarkdown,
		)
		if editErr != nil {
			log.Printf("[FilterApplyHandler] ERROR editing message: %v", editErr)
			return editErr
		}

		log.Printf("[FilterApplyHandler] Successfully edited message with filter %s", filter)
		return nil
	}
}

// PageNextHandler handles next page button
func PageNextHandler(
	client *managedesk.Client,
	sessionStore *context.SessionStore,
) tele.HandlerFunc {

	return func(c tele.Context) error {
		userID := int64(c.Sender().ID)
		log.Printf("[PageNextHandler] Starting - UserID: %d", userID)

		// Get current state from session
		userState := sessionStore.GetUserState(userID)
		currentPage := userState.CurrentPage
		log.Printf("[PageNextHandler] Current state - Page: %d, Filter: %s", currentPage, userState.SelectedFilter)

		// Create screen with session and move to next page
		screen := screens.NewTicketsListScreenWithSession(
			client,
			sessionStore,
			userID,
		)
		nextPage := currentPage + 1
		log.Printf("[PageNextHandler] Moving to page: %d", nextPage)
		screen.SetPage(nextPage)

		rendered, err := screen.Render()
		if err != nil {
			log.Printf("[PageNextHandler] ERROR rendering page %d: %v", nextPage, err)
			return c.Respond(&tele.CallbackResponse{
				Text: "Error loading next page",
				// Alert: true,
			})
		}

		log.Printf("[PageNextHandler] Successfully rendered page %d", nextPage)
		editErr := c.Edit(
			rendered.Text,
			rendered.Markup,
			tele.ModeMarkdown,
		)
		if editErr != nil {
			log.Printf("[PageNextHandler] ERROR editing message: %v", editErr)
			return editErr
		}

		log.Printf("[PageNextHandler] Successfully edited message for page %d", nextPage)
		return nil
	}
}

// PagePrevHandler handles previous page button
func PagePrevHandler(
	client *managedesk.Client,
	sessionStore *context.SessionStore,
) tele.HandlerFunc {

	return func(c tele.Context) error {
		userID := int64(c.Sender().ID)
		log.Printf("[PagePrevHandler] Starting - UserID: %d", userID)

		// Get current state from session
		userState := sessionStore.GetUserState(userID)
		currentPage := userState.CurrentPage
		log.Printf("[PagePrevHandler] Current state - Page: %d, Filter: %s", currentPage, userState.SelectedFilter)

		// Only go back if current page > 1
		if currentPage <= 1 {
			log.Printf("[PagePrevHandler] Already on first page, no action taken")
			return c.Respond(&tele.CallbackResponse{
				Text: "Already on first page",
			})
		}

		// Create screen with session and move to previous page
		screen := screens.NewTicketsListScreenWithSession(
			client,
			sessionStore,
			userID,
		)
		prevPage := currentPage - 1
		log.Printf("[PagePrevHandler] Moving to page: %d", prevPage)
		screen.SetPage(prevPage)

		rendered, err := screen.Render()
		if err != nil {
			log.Printf("[PagePrevHandler] ERROR rendering page %d: %v", prevPage, err)
			return c.Respond(&tele.CallbackResponse{
				Text: "Error loading previous page",
				// Alert: true,
			})
		}

		log.Printf("[PagePrevHandler] Successfully rendered page %d", prevPage)
		editErr := c.Edit(
			rendered.Text,
			rendered.Markup,
			tele.ModeMarkdown,
		)
		if editErr != nil {
			log.Printf("[PagePrevHandler] ERROR editing message: %v", editErr)
			return editErr
		}

		log.Printf("[PagePrevHandler] Successfully edited message for page %d", prevPage)
		return nil
	}
}

// TicketsListBackHandler returns to tickets list from detail or filter view
// Preserves user's filter and page state
func TicketsListBackHandler(
	client *managedesk.Client,
	sessionStore *context.SessionStore,
) tele.HandlerFunc {

	return func(c tele.Context) error {
		userID := int64(c.Sender().ID)

		// Return to tickets list with preserved state
		screen := screens.NewTicketsListScreenWithSession(
			client,
			sessionStore,
			userID,
		)

		rendered, err := screen.Render()
		if err != nil {
			return c.Respond(&tele.CallbackResponse{
				Text: "Error loading tickets list",
				// Alert: true,
			})
		}

		return c.Edit(
			rendered.Text,
			rendered.Markup,
			tele.ModeMarkdown,
		)
	}
}
