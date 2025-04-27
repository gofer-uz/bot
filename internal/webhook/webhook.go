// Package webhook - Telegram webhook funksionalligini ta'minlovchi paket
// Bu paket webhook serverini yaratish va boshqarish uchun ishlatiladi
package webhook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"tg-bot/pkg/logger"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Config webhook sozlamalarini o'z ichiga olgan interfeys
type Config interface {
	// WebhookURL manzilini qaytaradi
	WebhookURL() string
	// WebhookPort portini qaytaradi
	WebhookPort() string
	// GetTelegramToken Telegram bot tokenini qaytaradi
	GetTelegramToken() string
}

// Server webhook serverini yaratish va boshqarish uchun tuzilma
type Server struct {
	bot        *tgbotapi.BotAPI
	config     Config
	logger     *logger.Logger
	httpServer *http.Server
	updateChan chan tgbotapi.Update
}

// NewServer yangi webhook server yaratadi
func NewServer(bot *tgbotapi.BotAPI, config Config, log *logger.Logger) *Server {
	return &Server{
		bot:        bot,
		config:     config,
		logger:     log,
		updateChan: make(chan tgbotapi.Update, 100), // Update kanalini bufer bilan yaratamiz
	}
}

// Setup webhook serverini sozlaydi va ishga tushiradi
func (s *Server) Setup() error {
	// First, remove any existing webhook
	removeConfig := tgbotapi.DeleteWebhookConfig{
		DropPendingUpdates: false, // Keep existing updates
	}

	_, err := s.bot.Request(removeConfig)
	if err != nil {
		s.logger.Warnf("Failed to remove existing webhook: %v", err)
		// Continue anyway
	}

	// Webhook manzilini olish
	webhookURLStr := s.config.WebhookURL()

	// Parse webhook URL string into a URL object
	webhookURL, err := url.Parse(webhookURLStr)
	if err != nil {
		s.logger.Errorf("Invalid webhook URL: %v", err)
		return fmt.Errorf("invalid webhook URL: %w", err)
	}

	// Append bot token as path if not already present
	if len(webhookURL.Path) <= 1 {
		// Make sure the path includes the token
		webhookURL.Path = "/" + s.bot.Token
	}

	// Log the complete webhook URL
	s.logger.Infof("Setting webhook URL to: %s", webhookURL.String())

	// Create webhook config without certificate for domains with valid SSL
	webhookConfig := tgbotapi.WebhookConfig{
		URL:            webhookURL,
		MaxConnections: 40,
	}

	// Webhook ni o'rnatish
	s.logger.Info("Registering webhook with Telegram...")
	resp, err := s.bot.Request(webhookConfig)
	if err != nil {
		s.logger.Errorf("Webhook registration error: %v", err)
		return fmt.Errorf("webhook o'rnatishda xatolik: %w", err)
	}

	// Log the full response
	s.logger.Infof("Webhook registration response: %v", resp)

	// Webhook info ni tekshirish
	info, err := s.bot.GetWebhookInfo()
	if err != nil {
		return fmt.Errorf("webhook ma'lumotlarini olishda xatolik: %w", err)
	}

	// Log webhook info for debugging
	s.logger.Infof("Webhook info: %+v", info)

	// Webhook holatini tekshirish
	if info.LastErrorDate != 0 {
		s.logger.Warnf("Webhook xatoligi: %s", info.LastErrorMessage)
	} else {
		s.logger.Info("Webhook successfully registered with no errors!")
	}

	// Webhook manzili va endpointini yaratish
	webhookEndpoint := "/" + s.bot.Token

	// HTTP handler sozlash
	http.HandleFunc(webhookEndpoint, func(w http.ResponseWriter, r *http.Request) {
		// Log all incoming requests
		s.logger.Infof("Received webhook request from: %s %s", r.RemoteAddr, r.URL.Path)

		if r.Method != http.MethodPost {
			s.logger.Warnf("Rejected non-POST request: %s", r.Method)
			http.Error(w, "Faqat POST so'rovlari qabul qilinadi", http.StatusMethodNotAllowed)
			return
		}

		// Telegram update ni qabul qilish
		update, err := s.bot.HandleUpdate(r)
		if err != nil {
			s.logger.Errorf("Update ni qayta ishlashda xatolik: %v", err)
			http.Error(w, "Update ni qayta ishlashda xatolik", http.StatusBadRequest)
			return
		}

		// Log successful update
		s.logger.Infof("Received valid update ID: %d", update.UpdateID)

		// Update ni kanalga yuborish
		s.updateChan <- *update
		w.WriteHeader(http.StatusOK)
	})

	// Add a health check endpoint to verify server is working
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Log the health check request
		s.logger.Infof("Health check request received from: %s", r.RemoteAddr)

		// Return bot information and webhook status
		webhookInfo, _ := s.bot.GetWebhookInfo()
		responseData := map[string]interface{}{
			"status": "ok",
			"bot": map[string]interface{}{
				"username": s.bot.Self.UserName,
				"id":       s.bot.Self.ID,
			},
			"webhook": map[string]interface{}{
				"url":         webhookInfo.URL,
				"is_set":      webhookInfo.URL != "",
				"last_error":  webhookInfo.LastErrorMessage,
				"error_date":  webhookInfo.LastErrorDate,
				"pending":     webhookInfo.PendingUpdateCount,
				"ip_address":  webhookInfo.IPAddress,
				"server_time": time.Now().Format(time.RFC3339),
			},
		}

		// Set JSON content type
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Write JSON response
		jsonBytes, _ := json.Marshal(responseData)
		w.Write(jsonBytes)
	})

	// Add a simple test endpoint
	http.HandleFunc("/webhook-test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Webhook server is running!"))
	})

	// Info chiqarish
	s.logger.Infof("Webhook registered with URL: %s", webhookURL.String())
	s.logger.Infof("Webhook endpoint listening on: %s", webhookEndpoint)
	s.logger.Info("Health check endpoint available at: /health")

	return nil
}

// Start webhook serverni ishga tushiradi
func (s *Server) Start() error {
	// Webhook portini olish
	port := s.config.WebhookPort()

	// HTTP serverni sozlash
	s.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Server ishga tushirish
	s.logger.Infof("Webhook server %s portida ishlamoqda", port)

	// Go-routine da serverni ishga tushirish
	go func() {
		// HTTP server
		s.logger.Info("HTTP server ishga tushirilmoqda")
		err := s.httpServer.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			s.logger.Errorf("HTTP server xatoligi: %v", err)
		}
	}()

	return nil
}

// Updates webhook orqali kelgan yangilanishlar kanalini qaytaradi
func (s *Server) Updates() <-chan tgbotapi.Update {
	return s.updateChan
}

// Stop webhook serverni to'xtatadi
func (s *Server) Stop() error {
	if s.httpServer != nil {
		s.logger.Info("Webhook server to'xtatilmoqda...")
		return s.httpServer.Close()
	}
	return nil
}
