// Package config dastur konfiguratsiya sozlamalarini boshqarish uchun mo'ljallangan
// Bu paket .env faylidan va tizim muhit o'zgaruvchilaridan sozlamalarni yuklaydi
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config tuzilmasi dastur uchun barcha kerakli sozlamalarni saqlaydi
// Bu tuzilma bot ishga tushganda bir marta yuklanadi va butun dastur davomida ishlatiladi
type Config struct {
	TelegramToken string // Telegram bot tokeni - Botfather tomonidan berilgan maxsus identifikator
	LogLevel      string // Log darajasi - qancha batafsil ma'lumot saqlanishini belgilaydi (debug, info, warn, error)
}

// GetTelegramToken Telegram bot tokenini qaytaruvchi metod
// Bu metod Config interfeysi talablarini qondirish uchun ishlatiladi
func (c *Config) GetTelegramToken() string {
	return c.TelegramToken
}

// LoadConfig konfiguratsiya sozlamalarini .env faylidan va tizim muhit o'zgaruvchilaridan yuklaydi
// Bu funksiya dastur ishga tushganda eng birinchi chaqirilishi kerak
func LoadConfig() *Config {
	// .env faylini yuklash (agar mavjud bo'lsa)
	// Bu asosan mahalliy ishlab chiqish muhitida foydali hisoblanadi
	err := godotenv.Load()
	if err != nil {
		// Xatolik yuzaga kelsa, bu faqat ogohlantirish sifatida ko'rib chiqiladi,
		// chunki .env fayli bo'lmasligi mumkin (masalan, ishlab chiqarish muhitida)
		log.Println("Ogohlantirish: .env fayli topilmadi, tizim muhit o'zgaruvchilari ishlatiladi")
	}

	// Yangi konfiguratsiya obyektini yaratish
	cfg := &Config{
		// Muhit o'zgaruvchilaridan qiymatlarni olish, agar mavjud bo'lmasa standart qiymatlarni ishlatish
		TelegramToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
	}

	// Telegram tokeni mavjudligini tekshirish, chunki u bot ishlashi uchun muhim
	if cfg.TelegramToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN o'zgaruvchisi topilmadi. Iltimos, .env faylida yoki tizim muhit o'zgaruvchilarida sozlang.")
	}

	return cfg
}

// getEnv muhit o'zgaruvchisini qiymatini oladi, agar topilmasa belgilangan standart qiymatni qaytaradi
// Bu yordamchi funksiya konfiguratsiya yuklash jarayonida muhit o'zgaruvchilarini xavfsiz olish uchun ishlatiladi
func getEnv(key, defaultValue string) string {
	// Avval muhit o'zgaruvchisi qiymatini olishga harakat qilish
	value := os.Getenv(key)
	// Agar qiymat bo'sh bo'lsa, standart qiymatni qaytarish
	if value == "" {
		return defaultValue
	}
	return value
}
