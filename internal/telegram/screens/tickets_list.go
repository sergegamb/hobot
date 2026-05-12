package screens

import (
    "strings"

    "github.com/sergegamb/hobot/internal/managedesk"

    tele "gopkg.in/telebot.v4"
)

type RenderedScreen struct {
    Text string

    Markup *tele.ReplyMarkup
}

type TicketsListScreen struct {
    Client *managedesk.Client
}

func (s *TicketsListScreen) Render() (
    *RenderedScreen,
    error,
) {

    requests, err := s.Client.GetRequests()
    if err != nil {
        return nil, err
    }

    markup := &tele.ReplyMarkup{}

    rows := []tele.Row{}

    var text strings.Builder

    text.WriteString(
        "*All Tickets*",
    )

    for _, req := range requests {
        btn := markup.Data(
            req.Subject,
            "ticket_view",
            req.ID,
        )

        rows = append(
            rows,
            markup.Row(btn),
        )
    }

    markup.Inline(rows...)

    return &RenderedScreen{
        Text: text.String(),
        Markup: markup,
    }, nil
}
