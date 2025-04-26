// Package main - dasturning asosiy kirish nuqtasi
/*
Bu fayl dasturning asosiy kirish nuqtasi hisoblanadi.
Bu yerda loyiha ishga tushiriladi.
*/
package main

import (
	"tg-bot/cmd/bot"
	"tg-bot/internal/config"
	"tg-bot/pkg/logger"
)

func main() {
	// Konfiguratsiyani yuklash
	cfg := config.LoadConfig()

	// Logger yaratish
	log := logger.New(cfg.LogLevel)
	log.Info("Bot ishga tushmoqda...")

	// Bot funksiyasini chaqirish
	bot.RunBot(cfg, log)
}
