// Package config dastur konfiguratsiya sozlamalarini boshqarish uchun mo'ljallangan
// Bu paket config.yaml faylidan sozlamalarni yuklaydi
package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config tuzilmasi dastur uchun barcha kerakli sozlamalarni saqlaydi
// Bu tuzilma bot ishga tushganda bir marta yuklanadi va butun dastur davomida ishlatiladi
type Config struct {
	TelegramToken string `yaml:"telegram_token"` // Telegram bot tokeni - Botfather tomonidan berilgan maxsus identifikator
	LogLevel      string `yaml:"log_level"`      // Log darajasi - qancha batafsil ma'lumot saqlanishini belgilaydi (debug, info, warn, error)
	Mode          string `yaml:"mode"`           // Bot ishlash rejimi - webhook yoki polling
	Webhook       struct {
		URL  string `yaml:"url"`  // Webhook URL manzili - faqat webhook rejimida ishlatiladi
		Port string `yaml:"port"` // Webhook porti - faqat webhook rejimida ishlatiladi
	} `yaml:"webhook"`
}

// GetTelegramToken Telegram bot tokenini qaytaruvchi metod
// Bu metod Config interfeysi talablarini qondirish uchun ishlatiladi
func (c *Config) GetTelegramToken() string {
	return c.TelegramToken
}

// IsWebhookMode botning webhook rejimida ishlashini tekshiradi
func (c *Config) IsWebhookMode() bool {
	return strings.ToLower(c.Mode) == "webhook"
}

// WebhookURL webhook URL manzilini qaytaradi
func (c *Config) WebhookURL() string {
	return c.Webhook.URL
}

// WebhookPort webhook portini qaytaradi
func (c *Config) WebhookPort() string {
	return c.Webhook.Port
}

// LoadConfig konfiguratsiya sozlamalarini config.yaml faylidan yuklaydi
// Bu funksiya dastur ishga tushganda eng birinchi chaqirilishi kerak
func LoadConfig() *Config {
	// Standart qiymatlar bilan konfiguratsiya obyektini yaratish
	cfg := &Config{
		LogLevel: "info",
		Mode:     "polling",
	}
	cfg.Webhook.Port = "8443" // Webhook uchun standart port

	// Birinchi navbatda "config.yaml" ni tekshiramiz
	configPaths := []string{
		"config.yaml",                                      // Asosiy direktoriyada
		"configs/config.yaml",                              // configs papkasida
		filepath.Join("configs", "config.yaml"),            // Absolute path with configs folder
		filepath.Join("internal", "config", "config.yaml"), // internal/config papkasida
	}

	var configFile string
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			configFile = path
			break
		}
	}

	// Agar config.yaml topilmasa
	if configFile == "" {
		log.Println("Ogohlantirish: config.yaml fayli topilmadi, standart qiymatlar ishlatiladi")
		return createDefaultConfigIfNeeded(cfg)
	}

	// YAML faylini o'qish
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		log.Printf("Config faylini o'qishda xatolik: %v", err)
		return createDefaultConfigIfNeeded(cfg)
	}

	// YAML faylini Config strukturasiga o'girish
	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		log.Printf("YAML formatini qayta ishlashda xatolik: %v", err)
		return createDefaultConfigIfNeeded(cfg)
	}

	// Telegram tokeni mavjudligini tekshirish, chunki u bot ishlashi uchun muhim
	if cfg.TelegramToken == "" {
		log.Fatal("Telegram token topilmadi. Iltimos, config.yaml faylida 'telegram_token' parametrini sozlang.")
	}

	// Webhook rejimida webhook URL manzili tekshiriladi
	if cfg.IsWebhookMode() && cfg.WebhookURL() == "" {
		log.Fatal("Webhook rejimida ishlash uchun webhook URL manzili kerak. Iltimos, config.yaml faylida 'webhook.url' parametrini sozlang.")
	}

	return cfg
}

// createDefaultConfigIfNeeded agar config fayli topilmasa yangi standart config yaratadi
func createDefaultConfigIfNeeded(cfg *Config) *Config {
	// Telegram tokenini muhit o'zgaruvchilaridan olishga harakat qilish
	if token := os.Getenv("TELEGRAM_BOT_TOKEN"); token != "" {
		cfg.TelegramToken = token
	}

	// Taklif qilinadigan standart konfiguratsiya yaratish
	defaultConfig := `# Bot konfiguratsiyasi
telegram_token: "" # Botfather tomonidan berilgan token
log_level: "info"  # debug, info, warn, error
mode: "polling"    # webhook yoki polling

# Webhook sozlamalari (faqat webhook rejimida ishlatiladi)
webhook:
  url: ""          # https://example.com/your_token
  port: "8443"     # 8443, 443, 80, 88 yoki 8080
`

	// Standart config faylini yaratish (configs papkasida)
	configDir := "configs"
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			log.Printf("Configs papkasini yaratishda xatolik: %v", err)
		}
	}

	configPath := filepath.Join(configDir, "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("Standart config.yaml fayli yaratilmoqda:", configPath)
		if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
			log.Printf("Standart config faylini yaratishda xatolik: %v", err)
		}
	}

	// Telegram tokeni mavjudligini tekshirish
	if cfg.TelegramToken == "" {
		log.Fatal("Telegram token topilmadi. Iltimos, configs/config.yaml faylida 'telegram_token' parametrini sozlang.")
	}

	return cfg
}
