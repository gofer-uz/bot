// Package bot Telegram bot ishlashini ta'minlovchi asosiy paket
// Bu paket botning barcha asosiy funksionalligini o'z ichiga oladi
package bot

import (
	"tg-bot/internal/handlers"
	"tg-bot/internal/webhook"
	"tg-bot/pkg/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Config interfeysi bot uchun zarur sozlamalarni belgilaydi
// Bu interfeys orqali turli manbalardagi konfiguratsiyalarni ishlatish mumkin
type Config interface {
	// GetTelegramToken Telegram bot tokenini qaytaradi
	GetTelegramToken() string
	// IsWebhookMode botning webhook rejimida ishlashini tekshiradi
	IsWebhookMode() bool
}

// WebhookConfig webhook rejimini konfiguratsiya qilish uchun interfeys
type WebhookConfig interface {
	Config
	// WebhookURL manzilini qaytaradi
	WebhookURL() string
	// WebhookPort portini qaytaradi
	WebhookPort() string
}

// deleteWebhook mavjud webhook konfiguratsiyasini Telegram serveridan o'chiradi
// Bu funksiya webhook va polling rejimlari orasida toza o'tishni ta'minlash uchun muhim
func deleteWebhook(bot *tgbotapi.BotAPI, log *logger.Logger) {
	log.Info("Mavjud webhook konfiguratsiyasi o'chirilmoqda...")

	// Webhook o'chirish konfiguratsiyasini yaratish
	removeConfig := tgbotapi.DeleteWebhookConfig{
		DropPendingUpdates: false, // Mavjud yangilanishlarni saqlab qolish
	}

	// Webhook o'chirilishini so'rash
	_, err := bot.Request(removeConfig)
	if err != nil {
		log.Warnf("Webhook o'chirishda xatolik yuz berdi: %v", err)
		return
	}

	log.Info("Webhook muvaffaqiyatli o'chirildi")
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

	// Bot buyruqlarini ro'yxatdan o'tkazish
	handlers.RegisterBotCommands(bot, log)

	// Bot rejimiga qarab ishlash
	if cfg.IsWebhookMode() {
		log.Info("Bot webhook rejimida ishlamoqda")
		// Use the webhook package implementation
		webhookConfig, ok := cfg.(WebhookConfig)
		if !ok {
			log.Error("Webhook konfiguratsiyasi noto'g'ri")
			return
		}
		runWebhookMode(bot, webhookConfig, log)
	} else {
		log.Info("Bot polling rejimida ishlamoqda")
		// Always delete any existing webhook before starting polling mode
		deleteWebhook(bot, log)
		runPollingMode(bot, log)
	}
}

// runWebhookMode botni webhook rejimida ishga tushiradi
// Bu rejim ishlab chiqarish muhiti uchun tavsiya etiladi
func runWebhookMode(bot *tgbotapi.BotAPI, cfg WebhookConfig, log *logger.Logger) {
	// Import webhook package and use the implemented Server
	webhookServer := webhook.NewServer(bot, cfg, log)

	// Setup the webhook server
	if err := webhookServer.Setup(); err != nil {
		log.Errorf("Webhook serverini sozlashda xatolik: %v", err)
		return
	}

	// Start the webhook server
	if err := webhookServer.Start(); err != nil {
		log.Errorf("Webhook serverini ishga tushirishda xatolik: %v", err)
		return
	}

	// Get update channel from webhook
	updates := webhookServer.Updates()

	// Log webhook info
	info, err := bot.GetWebhookInfo()
	if err == nil {
		if info.LastErrorDate != 0 {
			log.Warnf("Webhook xatoligi: %s", info.LastErrorMessage)
		} else {
			log.Info("Webhook muvaffaqiyatli o'rnatildi")
		}
	}

	// Yangilanishlarni qayta ishlash
	for update := range updates {
		go handleUpdate(bot, update, log)
	}
}

// runPollingMode botni polling rejimida ishga tushiradi
// Bu rejim rivojlantirish muhiti uchun tavsiya etiladi
func runPollingMode(bot *tgbotapi.BotAPI, log *logger.Logger) {
	// Yangilanishlar konfiguratsiyasini sozlash
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60 // Kutish vaqti (sekundlarda)

	// Yangilanishlar kanalini olish
	updates := bot.GetUpdatesChan(updateConfig)

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
