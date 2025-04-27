// Package bot Telegram bot ishlashini ta'minlovchi asosiy paket
// Bu paket botning barcha asosiy funksionalligini o'z ichiga oladi
package bot

import (
	"tg-bot/internal/handlers"
	"tg-bot/pkg/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Config interfeysi bot uchun zarur sozlamalarni belgilaydi
// Bu interfeys orqali turli manbalardagi konfiguratsiyalarni ishlatish mumkin
type Config interface {
	// GetTelegramToken Telegram bot tokenini qaytaradi
	GetTelegramToken() string
}

// RunBot Telegram botni ishga tushirish va boshqarish uchun asosiy funksiya
// Bu funksiya botni yaratadi, sozlaydi va yangilanishlarni qabul qilishni boshlaydi
func RunBot(cfg Config, log *logger.Logger) {
	// Yangi bot namunasini yaratish
	bot, err := tgbotapi.NewBotAPI(cfg.GetTelegramToken())
	if err != nil {
		log.Error("Bot yaratishda xatolik yuz berdi:", err)
		return
	}

	// Nosozliklarni tuzatish rejimini yoqish (ixtiyoriy)
	bot.Debug = true
	log.Info("Bot muvaffaqiyatli ishga tushirildi:", bot.Self.UserName)

	// Yangilanishlar konfiguratsiyasini sozlash
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60 // Kutish vaqti (sekundlarda)

	// Yangilanishlar kanalini olish
	updates := bot.GetUpdatesChan(updateConfig)

	// Bot buyruqlarini ro'yxatdan o'tkazish
	handlers.RegisterBotCommands(bot, log)

	// Yangilanishlarni qayta ishlash
	for update := range updates {
		go handleUpdate(bot, update, log) // Har bir yangilanish uchun alohida go-routineda ishlaymiz
	}
}

// handleUpdate har bir kiruvchi yangilanishni qayta ishlaydi
// Bu funksiya xabarlar, buyruqlar va callback so'rovlarni aniqlaydi va ularga javob beradi
func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update, log *logger.Logger) {
	// Yangi a'zo guruhga qo'shilganligini tekshirish
	if update.Message != nil && update.Message.NewChatMembers != nil && len(update.Message.NewChatMembers) > 0 {
		for _, newUser := range update.Message.NewChatMembers {
			// Bot o'zi qo'shilganini e'tiborga olmaslik
			if newUser.ID == bot.Self.ID {
				continue
			}

			// Send a simple mention message
			mentionNewUser(bot, update.Message.Chat.ID, newUser, log)
		}
		return
	}

	// Buyruqlarni qayta ishlash
	if update.Message != nil && update.Message.IsCommand() {
		command := update.Message.Command()
		// Only log commands, not regular messages
		log.Infof("Foydalanuvchi \"%s\" buyrug'ini yubordi", command)

		// Buyruqni tegishli qayta ishlovchiga uzatish
		if handler := handlers.GetCommandHandler(command); handler != nil {
			handler(bot, update.Message, log)
		} else {
			// Agar buyruq ma'lum bo'lmasa, foydalanuvchiga yordam xabarini yuborish
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Noma'lum buyruq. Mavjud buyruqlar ro'yxatini ko'rish uchun /help buyrug'ini ishlatib ko'ring")
			bot.Send(msg)
		}
		return
	}

	// Oddiy xabarlarni qayta ishlash
	if update.Message != nil {
		// Remove debug logging and message echoing for regular messages
		// No need to resend messages that bot receives from groups
		return
	}

	// Callback so'rovlarini qayta ishlash (inline klaviaturalar uchun)
	if update.CallbackQuery != nil {
		log.Debugf("Callback so'rovi qabul qilindi: %s", update.CallbackQuery.Data)
		handlers.HandleCallback(bot, update.CallbackQuery, log)
		return
	}
}

// mentionNewUser mentions a new user with a subtle suggestion to check the bot for community info
// This function is called when a new user joins the group
func mentionNewUser(bot *tgbotapi.BotAPI, chatID int64, user tgbotapi.User, log *logger.Logger) {
	// Get user's mention format
	var userMention string
	if user.UserName != "" {
		userMention = "@" + user.UserName
	} else {
		userMention = user.FirstName
	}

	// Create a subtle mention message
	welcomeMsg := tgbotapi.NewMessage(chatID,
		"Assalomu alaykum "+userMention+"! Bizni hamjamiyat haqida ko'proq bilish uchun botga murojaat qiling.")

	// Create an inline keyboard with a Start button
	startButton := tgbotapi.NewInlineKeyboardButtonURL("Botga tashrif buyirish", "https://t.me/"+bot.Self.UserName+"?start=welcome")
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(startButton),
	)

	// Attach the keyboard to the message
	welcomeMsg.ReplyMarkup = keyboard

	// Send the message
	_, err := bot.Send(welcomeMsg)
	if err != nil {
		log.Error("Error sending mention message to new user:", err)
	}
}
