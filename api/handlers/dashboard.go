package handlers

import (
	"fmt"

	"github.com/masudur-rahman/expense-tracker-bot/configs"
	"github.com/masudur-rahman/expense-tracker-bot/services/all"

	"gopkg.in/telebot.v3"
)

// Dashboard sends an inline button linking to the web dashboard.
func Dashboard(ctx telebot.Context) error {
	svc := all.GetServices()
	user, err := svc.User.GetUserByID(ctx.Sender().ID)
	if err != nil {
		return ctx.Send("Please /start the bot first to create your profile.")
	}

	if user.Username == "" && user.MobileNumber == "" {
		return ctx.Send("Set a Telegram username or share your phone number first so the dashboard can identify you.")
	}

	cfg := configs.TrackerConfig.WebDashboard
	if !cfg.Enabled {
		return ctx.Send("Web dashboard is not enabled.")
	}

	dashboardURL := buildDashboardURL(cfg)

	btn := &telebot.ReplyMarkup{}
	webBtn := btn.URL("Open Dashboard", dashboardURL)
	btn.Inline(btn.Row(webBtn))

	return ctx.Send("Open your expense dashboard:", btn)
}

// buildDashboardURL constructs the frontend URL from config.
func buildDashboardURL(cfg configs.WebDashboardConfig) string {
	origin := cfg.CORSOrigin
	if origin == "" || origin == "*" {
		port := cfg.Port
		if port == "" {
			port = "3000"
		}
		return fmt.Sprintf("http://localhost:%s", port)
	}
	return origin
}
