// Package logger dastur uchun log yozish funksionalligini ta'minlaydi
// Bu paket turli darajalardagi xabarlarni tizimli ravishda qayd etishga imkon beradi
package logger

import (
	"log"
	"os"
	"strings"
)

// Logger - maxsus log yozish tuzilmasi
// Bu tuzilma turli darajadagi xabarlarni qayd etish uchun mo'ljallangan
type Logger struct {
	infoLogger  *log.Logger // Ma'lumot xabarlarini qayd etish uchun
	warnLogger  *log.Logger // Ogohlantirish xabarlarini qayd etish uchun
	errorLogger *log.Logger // Xato xabarlarini qayd etish uchun
	debugLogger *log.Logger // Nosozliklarni tuzatish xabarlarini qayd etish uchun
	level       int         // Joriy log darajasi
}

// Log darajalari. Raqamlar katta bo'lgan sari muhimlik darajasi kamayadi
const (
	LevelError = iota // Xato darajasi - faqat xatolarni qayd etish
	LevelWarn         // Ogohlantirish darajasi - xatolar va ogohlantirishlarni qayd etish
	LevelInfo         // Ma'lumot darajasi - xatolar, ogohlantirishlar va ma'lumotlarni qayd etish
	LevelDebug        // Nosozliklarni tuzatish darajasi - barcha xabarlarni qayd etish
)

// New - ko'rsatilgan darajada yangi logger yaratadi
// Bu funksiya dastur boshida bir marta chaqiriladi
func New(level string) *Logger {
	var logLevel int

	// Kiritilgan matn asosida log darajasini aniqlash
	switch strings.ToLower(level) {
	case "debug":
		logLevel = LevelDebug
	case "info":
		logLevel = LevelInfo
	case "warn":
		logLevel = LevelWarn
	case "error":
		logLevel = LevelError
	default:
		// Agar noto'g'ri daraja ko'rsatilgan bo'lsa, standart qiymat sifatida info darajasi tanlanadi
		logLevel = LevelInfo
	}

	// Yangi logger obyekti yaratish va sozlash
	return &Logger{
		// Har bir log turi uchun alohida formatlash
		infoLogger:  log.New(os.Stdout, "MA'LUMOT: ", log.Ldate|log.Ltime),
		warnLogger:  log.New(os.Stdout, "OGOHLANTIRISH: ", log.Ldate|log.Ltime),
		errorLogger: log.New(os.Stderr, "XATO: ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLogger: log.New(os.Stdout, "NOSOZLIK: ", log.Ldate|log.Ltime|log.Lshortfile),
		level:       logLevel,
	}
}

// Debug - nosozliklarni tuzatish xabarlarini qayd etadi
// Bu metod faqat log darajasi Debug bo'lganida ishga tushadi
func (l *Logger) Debug(v ...interface{}) {
	if l.level >= LevelDebug {
		l.debugLogger.Println(v...)
	}
}

// Info - ma'lumot xabarlarini qayd etadi
// Bu metod log darajasi Info yoki undan yuqori bo'lganida ishga tushadi
func (l *Logger) Info(v ...interface{}) {
	if l.level >= LevelInfo {
		l.infoLogger.Println(v...)
	}
}

// Warn - ogohlantirish xabarlarini qayd etadi
// Bu metod log darajasi Warn yoki undan yuqori bo'lganida ishga tushadi
func (l *Logger) Warn(v ...interface{}) {
	if l.level >= LevelWarn {
		l.warnLogger.Println(v...)
	}
}

// Error - xato xabarlarini qayd etadi
// Bu metod log darajasi Error yoki undan yuqori bo'lganida ishga tushadi
func (l *Logger) Error(v ...interface{}) {
	if l.level >= LevelError {
		l.errorLogger.Println(v...)
	}
}

// Debugf - formatli nosozliklarni tuzatish xabarlarini qayd etadi
// Bu metod faqat log darajasi Debug bo'lganida ishga tushadi
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.level >= LevelDebug {
		l.debugLogger.Printf(format, v...)
	}
}

// Infof - formatli ma'lumot xabarlarini qayd etadi
// Bu metod log darajasi Info yoki undan yuqori bo'lganida ishga tushadi
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.level >= LevelInfo {
		l.infoLogger.Printf(format, v...)
	}
}

// Warnf - formatli ogohlantirish xabarlarini qayd etadi
// Bu metod log darajasi Warn yoki undan yuqori bo'lganida ishga tushadi
func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.level >= LevelWarn {
		l.warnLogger.Printf(format, v...)
	}
}

// Errorf - formatli xato xabarlarini qayd etadi
// Bu metod log darajasi Error yoki undan yuqori bo'lganida ishga tushadi
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.level >= LevelError {
		l.errorLogger.Printf(format, v...)
	}
}
