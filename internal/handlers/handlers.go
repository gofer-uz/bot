package handlers

import (
	tele "gopkg.in/telebot.v3"
)

func HandleStart(c tele.Context) error {
	return c.Send("salom")
}

func HandleHelp(c tele.Context) error {
	msg := "nimalar qilsa bo'ladi:\n" +
		"• /useful – turkumlar bo'yicha foydali havolalar\n" +
		"\n"

	msg += "buyruqlar ro'yxati:\n" +
		"/start – botni boshlash yoki qayta ishga tushirish\n" +
		"/help – yordam va buyruqlar\n" +
		"/useful – foydali havolalar menyusi"

	return c.Send(msg)
}
