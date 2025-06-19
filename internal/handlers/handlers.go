package handlers

import (
	tele "gopkg.in/telebot.v3"
)

func HandleStart(c tele.Context) error {
	return c.Send("welcome! use /help to see commands")
}

func HandleHelp(c tele.Context) error {
	return c.Send("available commands:\n/start - start the bot\n/help - list commands")
}
