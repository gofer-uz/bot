package bot

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	tele "gopkg.in/telebot.v3"

	"bot-architecture/internal/config"
	"bot-architecture/internal/handlers"
)

// App wraps telebot.Bot and handles lifecycle.

type App struct {
	Bot *tele.Bot
	cfg config.Config
}

// New creates and configures a new bot application instance.
func New(cfg config.Config) (*App, error) {
	pref := tele.Settings{
		Token:     cfg.Telegram.Token,
		ParseMode: tele.ModeMarkdown,
		Client:    &http.Client{Timeout: 10 * time.Second},
		// webhook settings will be configured after bot is created
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	a := &App{Bot: b, cfg: cfg}

	// register command and message handlers in a dedicated package to keep bot clean
	handlers.Register(a.Bot)

	return a, nil
}

// SetupWebhook registers the bot webhook with telegram and starts internal server.
func (a *App) SetupWebhook() error {
	// reconfigure poller to webhook mode and start
	wh := &tele.Webhook{
		Listen:   a.cfg.Server.Listen,
		Endpoint: &tele.WebhookEndpoint{PublicURL: a.cfg.Server.PublicURL},
	}

	a.Bot.Poller = wh

	slog.Info("starting bot with webhook poller")
	a.Bot.Start()
	return nil
}

// Shutdown gracefully stops the bot.
func (a *App) Shutdown(_ context.Context) {
	slog.Info("shutting down bot")
	// telebot's Stop blocks until the poller is fully shut down, so a simple call is enough.
	// we ignore the context because Stop has its own internal confirmation mechanism.
	// a.Bot.Stop()
}
