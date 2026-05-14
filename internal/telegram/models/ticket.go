package models

import "fmt"

// TicketStatus represents the display status of a ticket
type TicketStatus string

const (
	StatusOpen    TicketStatus = "open"
	StatusClosed  TicketStatus = "closed"
	StatusPending TicketStatus = "pending"
)

// TicketDisplayFields holds all fields needed for ticket display in Telegram
type TicketDisplayFields struct {
	ID          string
	Subject     string
	Status      TicketStatus
	Priority    string
	Category    string
	Created     string
	Updated     string
	Requester   string
	Description string
}

// TicketListItem represents a ticket in list view
type TicketListItem struct {
	ID      string
	Subject string
	Status  TicketStatus
}

// TicketDetail represents a ticket in detail view
type TicketDetail struct {
	Fields TicketDisplayFields
}

// StatusEmoji returns an emoji representation of the ticket status
func StatusEmoji(status TicketStatus) string {
	switch status {
	case StatusOpen:
		return "🟢"
	case StatusClosed:
		return "✅"
	case StatusPending:
		return "⏳"
	default:
		return "❓"
	}
}

// FormatListButtonText formats a ticket for display in a list button
// Format: "[ID] Subject - Status"
func FormatListButtonText(item *TicketListItem) string {
	return fmt.Sprintf("[%s] %s %s", item.ID, item.Subject, StatusEmoji(item.Status))
}

// FormatDetailMessage formats a ticket for display in detail view
func FormatDetailMessage(detail *TicketDetail) string {
	f := detail.Fields
	message := fmt.Sprintf("*Ticket #%s*\n\n", f.ID)
	message += fmt.Sprintf("*Subject:* %s\n", f.Subject)
	message += fmt.Sprintf("*Status:* %s %s\n", StatusEmoji(TicketStatus(f.Status)), f.Status)

	if f.Priority != "" {
		message += fmt.Sprintf("*Priority:* %s\n", f.Priority)
	}

	if f.Category != "" {
		message += fmt.Sprintf("*Category:* %s\n", f.Category)
	}

	if f.Requester != "" {
		message += fmt.Sprintf("*Requester:* %s\n", f.Requester)
	}

	if f.Created != "" {
		message += fmt.Sprintf("*Created:* %s\n", f.Created)
	}

	if f.Updated != "" {
		message += fmt.Sprintf("*Updated:* %s\n", f.Updated)
	}

	if f.Description != "" {
		message += fmt.Sprintf("\n*Description:*\n%s\n", f.Description)
	}

	return message
}

// FormatFilterSelectionMessage formats the filter selection message
func FormatFilterSelectionMessage() string {
	return "*Available Filters/Displays*\n\nSelect a filter to display:"
}

// TruncateSubject truncates a subject to a maximum length with ellipsis
func TruncateSubject(subject string, maxLength int) string {
	if len(subject) <= maxLength {
		return subject
	}
	return subject[:maxLength-3] + "..."
}

func AppendedSubject(subject string) string {
	return subject + "                            "
}
