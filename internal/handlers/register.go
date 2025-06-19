package handlers

import tele "gopkg.in/telebot.v3"

// Register attaches all command and message handlers to the given bot instance.
// keep this file minimal; add new handlers in separate files within this package.
func Register(b *tele.Bot) {
	b.Handle("/start", HandleStart)

	b.Handle("/help", HandleHelp)

	registerUseful(b)
}
