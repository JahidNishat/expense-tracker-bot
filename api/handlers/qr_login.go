package handlers

import (
	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	"gopkg.in/telebot.v3"
)

// HandleQRLogin confirms a QR login session from a Telegram deep link.
func HandleQRLogin(sessionID string, ctx telebot.Context) error {
	authSvc := all.GetServices().Auth
	if authSvc == nil {
		return ctx.Send("Web dashboard is not enabled.")
	}

	if err := authSvc.ConfirmQRSession(sessionID, ctx.Sender().ID); err != nil {
		return ctx.Send("Login failed: " + err.Error())
	}

	return ctx.Send("Dashboard login confirmed. You can return to your browser.")
}
