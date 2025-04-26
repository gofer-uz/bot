// Package handlers Telegram bot buyruqlarini qayta ishlash va javoblar generatsiya qilish uchun mo'ljallangan
// Bu paket foydalanuvchi yuborgan barcha buyruqlarni qayta ishlash va tegishli javoblarni tayyorlash logikasini o'z ichiga oladi
package handlers

import (
	"fmt"
	"strings"

	"tg-bot/pkg/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// commandHandler butun dastur davomida ishlatiladigan global obyekt
// Bu barcha buyruqlar va ularning javoblarini saqlash uchun ishlatiladi
var commandHandler *CommandHandler

// CommandFunction muayyan buyruqni bajaradigan funksiya turi
// Har bir buyruq alohida funksiya sifatida implementatsiya qilinadi
type CommandFunction func(bot *tgbotapi.BotAPI, message *tgbotapi.Message, log *logger.Logger)

// commandHandlers barcha buyruq funksiyalarini saqlovchi xarita
// Bu xarita buyruq nomi va uni qayta ishlovchi funksiya o'rtasidagi bog'lanishni ta'minlaydi
var commandHandlers map[string]CommandFunction

// RegisterBotCommands botga barcha mavjud buyruqlarni ro'yxatdan o'tkazadi
// Bu funksiya bot ishga tushganda bir marta chaqiriladi va barcha buyruqlarni sozlaydi
func RegisterBotCommands(bot *tgbotapi.BotAPI, log *logger.Logger) {
	// Agar commandHandler yaratilmagan bo'lsa, yangi instance yaratish
	if commandHandler == nil {
		commandHandler = NewCommandHandler(log)
	}

	// Buyruqlar xaritasini yaratish
	commandHandlers = make(map[string]CommandFunction)

	// Har bir buyruq uchun qayta ishlovchi funksiyani ro'yxatdan o'tkazish
	// START buyrug'i - botni ishga tushirish va salomlashish xabarini yuborish
	commandHandlers["start"] = func(bot *tgbotapi.BotAPI, message *tgbotapi.Message, log *logger.Logger) {
		msg := tgbotapi.NewMessage(message.Chat.ID, commandHandler.GetStartText())
		bot.Send(msg)
	}

	// HELP buyrug'i - mavjud buyruqlar ro'yxati va ularning tavsifi
	commandHandlers["help"] = func(bot *tgbotapi.BotAPI, message *tgbotapi.Message, log *logger.Logger) {
		msg := tgbotapi.NewMessage(message.Chat.ID, commandHandler.GetHelpText())
		bot.Send(msg)
	}

	// RULES buyrug'i - hamjamiyat qoidalari
	commandHandlers["rules"] = func(bot *tgbotapi.BotAPI, message *tgbotapi.Message, log *logger.Logger) {
		msg := tgbotapi.NewMessage(message.Chat.ID, commandHandler.GetRulesText())
		bot.Send(msg)
	}

	// ABOUT buyrug'i - bot va uning maqsadi haqida ma'lumot
	commandHandlers["about"] = func(bot *tgbotapi.BotAPI, message *tgbotapi.Message, log *logger.Logger) {
		msg := tgbotapi.NewMessage(message.Chat.ID, commandHandler.GetAboutText())
		bot.Send(msg)
	}

	// GROUP buyrug'i - Go bo'yicha guruhlar va hamjamiyatlar haqida ma'lumot
	commandHandlers["group"] = func(bot *tgbotapi.BotAPI, message *tgbotapi.Message, log *logger.Logger) {
		msg := tgbotapi.NewMessage(message.Chat.ID, commandHandler.GetGroupText())
		bot.Send(msg)
	}

	// ROADMAP buyrug'i - Go o'rganish yo'l xaritasi
	commandHandlers["roadmap"] = func(bot *tgbotapi.BotAPI, message *tgbotapi.Message, log *logger.Logger) {
		msg := tgbotapi.NewMessage(message.Chat.ID, commandHandler.GetRoadmapText())
		bot.Send(msg)
	}

	// USEFUL buyrug'i - Go bo'yicha foydali resurslar
	commandHandlers["useful"] = func(bot *tgbotapi.BotAPI, message *tgbotapi.Message, log *logger.Logger) {
		msg := tgbotapi.NewMessage(message.Chat.ID, commandHandler.GetUsefulText())
		bot.Send(msg)
	}

	// LATEST buyrug'i - eng so'nggi Go versiyasi haqida ma'lumot
	commandHandlers["latest"] = func(bot *tgbotapi.BotAPI, message *tgbotapi.Message, log *logger.Logger) {
		msg := tgbotapi.NewMessage(message.Chat.ID, commandHandler.GetLatestText())
		bot.Send(msg)
	}

	// VERSION buyrug'i - so'ralgan Go versiyasi haqida batafsil ma'lumot
	commandHandlers["version"] = func(bot *tgbotapi.BotAPI, message *tgbotapi.Message, log *logger.Logger) {
		msg := tgbotapi.NewMessage(message.Chat.ID, commandHandler.GetVersionText(message.CommandArguments()))
		bot.Send(msg)
	}

	// WARN buyrug'i - foydalanuvchiga ogohlantirish xabarini yuborish
	commandHandlers["warn"] = func(bot *tgbotapi.BotAPI, message *tgbotapi.Message, log *logger.Logger) {
		msg := tgbotapi.NewMessage(message.Chat.ID, commandHandler.GetWarnText(message.From.UserName))
		bot.Send(msg)
	}

	log.Info("Bot buyruqlari ro'yxatdan o'tkazildi")
}

// GetCommandHandler ma'lum bir buyruq uchun qayta ishlovchi funksiyani qaytaradi
// Bu funksiya asosiy bot logikasi tomonidan buyruq aniqlanganda chaqiriladi
func GetCommandHandler(command string) CommandFunction {
	handler, exists := commandHandlers[command]
	if !exists {
		return nil
	}
	return handler
}

// HandleCallback inline klaviatura tugmachalaridan kelgan callback so'rovlarini qayta ishlaydi
// Bu funksiya foydalanuvchi inline tugmani bosganda chaqiriladi
func HandleCallback(bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, log *logger.Logger) {
	// Callback so'rovini qabul qilganligimizni Telegram'ga xabar berish
	// Bu foydalanuvchi interfeysi uchun muhim, chunki tugmani bosish animatsiyasini to'xtatadi
	callback_config := tgbotapi.NewCallback(callback.ID, "")
	bot.Send(callback_config)

	// Callback ma'lumotlarini qayta ishlash
	log.Debugf("Callback qabul qilindi: %s", callback.Data)

	switch callback.Data {
	case "about":
		// Bot haqida ma'lumot yuborish
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, commandHandler.GetAboutText())
		bot.Send(msg)
	case "roadmap":
		// Go o'rganish yo'l xaritasini yuborish
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, commandHandler.GetRoadmapText())
		bot.Send(msg)
	default:
		// Noma'lum callback ID kelsa, xatolik haqida ma'lumot berish
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Noma'lum tugma bosildi. Iltimos qaytadan urinib ko'ring.")
		bot.Send(msg)
	}
}

// CommandHandler buyruqlar va ularning mantiqini o'z ichiga oluvchi asosiy tuzilma
// Bu tuzilma barcha bot buyruqlari uchun javoblarni generatsiya qilish funksiyalarini o'z ichiga oladi
type CommandHandler struct {
	logger *logger.Logger // Logger obyekti, xatoliklar va hodisalarni qayd etish uchun
}

// NewCommandHandler yangi CommandHandler obyektini yaratish uchun konstruktor metod
// Bu funksiya har bir buyruq qayta ishlovchisi uchun kerakli narsalarni tayyorlaydi
func NewCommandHandler(logger *logger.Logger) *CommandHandler {
	return &CommandHandler{
		logger: logger,
	}
}

// HandleCommand barcha buyruqlar uchun qayta ishlash funksiyasi (eski usul)
// Bu metod to'g'ridan-to'g'ri Update obektini qabul qiladi va kerakli javoblarni yuboradi
func (h *CommandHandler) HandleCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	h.logger.Debug("Yangi buyruq qabul qilindi: ", update.Message.Command())

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	// Buyruq turini aniqlash va tegishli funksiyani chaqirish
	switch update.Message.Command() {
	case "start":
		msg.Text = h.GetStartText()
	case "help":
		msg.Text = h.GetHelpText()
	case "rules":
		msg.Text = h.GetRulesText()
	case "about":
		msg.Text = h.GetAboutText()
	case "group":
		msg.Text = h.GetGroupText()
	case "roadmap":
		msg.Text = h.GetRoadmapText()
	case "useful":
		msg.Text = h.GetUsefulText()
	case "latest":
		msg.Text = h.GetLatestText()
	case "version":
		msg.Text = h.GetVersionText(update.Message.CommandArguments())
	case "warn":
		msg.Text = h.GetWarnText(update.Message.From.UserName)
	default:
		msg.Text = "Noma'lum buyruq. Yordam olish uchun /help buyrug'ini ishlatib ko'ring."
	}

	// Javob xabarini yuborish
	if _, err := bot.Send(msg); err != nil {
		h.logger.Error("Xabar yuborishda xatolik yuz berdi: ", err)
	}
}

// GetStartText bot bilan ilk muloqotda yuboriladigan salomlashish xabari
// Bu xabar foydalanuvchi botdan qanday foydalanish mumkinligi haqida qisqacha ma'lumot beradi
func (h *CommandHandler) GetStartText() string {
	return `Assalomu alaykum! GoferUz Golang botiga xush kelibsiz 👋

Bu bot Go dasturlash tili bo'yicha ma'lumotlar va resurslarni taqdim etish uchun yaratilgan.

Mavjud buyruqlar ro'yxatini ko'rish uchun /help buyrug'ini yuboring.
O'zbekistondagi Go dasturchilar hamjamiyati haqida ma'lumot olish uchun /group buyrug'ini yuboring.`
}

// GetHelpText mavjud barcha buyruqlar ro'yxati va ularning qisqacha tavsifi
// Bu funksiya foydalanuvchi botdan to'liq foydalanish uchun qo'llanma vazifasini bajaradi
func (h *CommandHandler) GetHelpText() string {
	return `Mavjud komandalar ro'yxati:

/help - ushbu xabarni qayta ko'rsatish
/rules - qoidalarni aks ettirish
/about - ushbu botimizning rivojlantirish qismi
/group - Go ga oid guruh va hamjamiyatlar
/roadmap - boshlang'ich o'rganuvchilar uchun
/useful - Go haqida foydali yoki kerakli ma'lumotlar
/latest - eng oxirgi reliz haqida qisqacha ma'lumot
/version - biron anniq reliz haqida to'liq ma'lumot
/warn - mavzudan chetlashganga ogohlantiruv`
}

// GetRulesText hamjamiyat va guruh uchun qoidalar to'plami
// Bu qoidalar hamjamiyat a'zolari o'rtasida hurmat va professionallikni ta'minlaydi
func (h *CommandHandler) GetRulesText() string {
	return `GoferUz hamjamiyati qoidalari:

1. Hurmat bilan munosabatda bo'ling - boshqa a'zolarga nisbatan doimo hurmat va e'tibor ko'rsating
2. Spam yoki reklama tarqatmang - ruxsatsiz reklama materiallarini jo'natmang
3. Siyosiy va diniy mavzulardan chetlaning - guruh faqat Go dasturlash tili uchun
4. Go dasturlash tili bo'yicha savollarda aniq va foydali bo'ling
5. Maqsadimiz - O'zbekistonda Go dasturlash tilini rivojlantirish va Go jamoasini kengaytirish

Qoidalarga rioya qilmaslik ogohlantirishga, takrorlanishi esa guruhdan chetlashtirishga sabab bo'lishi mumkin.`
}

// GetAboutText bot haqida va uning yaratilishi to'g'risidagi ma'lumotlarni qaytaradi
// Bu funksiya bot arxitekturasi va uni kim tomonidan yaratilgani haqida ma'lumot beradi
func (h *CommandHandler) GetAboutText() string {
	return `Bu bot Go dasturlash tilida yaratilgan va Go dasturlash tiliga bag'ishlangan.
Botning asosiy maqsadi - Go o'rganuvchilar uchun foydali ma'lumotlarni tezkor taqdim etish va 
O'zbekistondagi Go hamjamiyatini qo'llab-quvvatlash.

Bot arxitekturasi:
- cmd/bot - asosiy dastur kirish nuqtasi va bot logikasi
- internal/config - konfiguratsiya sozlamalari va muhit o'zgaruvchilari bilan ishlash
- internal/handlers - barcha buyruqlarni qayta ishlash mantiqiy qismi
- pkg/logger - tizimli log yozish uchun maxsus kutubxona

Botning joriy versiyasi: 1.0.0
Muallif: haywan

Bot Golang tilida yozilgan. Go - bu Google tomonidan yaratilgan, yuqori samaradorlikka ega, 
statik tipli va kompilyatsiya qilinadigan zamonaviy dasturlash tili. Go dasturlari juda tezkor 
ishlaydi, parallel dasturlashni osonlashtiradi va xotiradan samarali foydalanadi.`
}

// GetGroupText Go dasturlash tili bo'yicha guruhlar va hamjamiyatlar haqida ma'lumot
// Bu funksiya Go bo'yicha O'zbekiston va xalqaro darajadagi hamjamiyatlar haqida ma'lumot beradi
func (h *CommandHandler) GetGroupText() string {
	return `Go dasturlash tili bo'yicha guruhlar va hamjamiyatlar:

🌍 O'zbekiston hamjamiyati:
- Telegram: @goferuz - O'zbekistondagi Go dasturchilar guruhi
- Veb-sayt: https://gopher.uz - O'zbekistonlik gopher'lar uchun portal

🌐 Xalqaro hamjamiyatlar:
- GitHub: https://github.com/goferuz - O'zbek Go dasturlari repozitoriyalari
- Forum: https://forum.golangbridge.org/ - Go dasturchilar forumi
- Reddit: https://www.reddit.com/r/golang/ - Go haqidagi Reddit jamiyati
- Slack: https://gophers.slack.com/ - Go dasturchilar uchun Slack kanali
- Discord: https://discord.gg/golang - Go dasturchilarning Discord serveri
- Stack Overflow: https://stackoverflow.com/questions/tagged/go - Go savollari bazasi

Ushbu hamjamiyatlarga qo'shilish orqali siz Go bo'yicha bilimlaringizni oshirish 
va tajribali dasturchilar bilan muloqot qilish imkoniyatiga ega bo'lasiz.`
}

// GetRoadmapText Go dasturlash tilini o'rganish uchun yo'l xaritasini qaytaradi
// Bu funksiya bosqichma-bosqich Go ni o'rganish rejasini taqdim etadi
func (h *CommandHandler) GetRoadmapText() string {
	return `Go dasturlash tilini o'rganish uchun mukammal yo'l xaritasi:

1️⃣ Go asoslari - o'zgaruvchilar, turlari va funksiyalar
   - O'zgaruvchilar va konstantalar deklaratsiyasi (var, const)
   - Asosiy ma'lumot turlari (int, float64, bool, string, rune)
   - Funksiyalar, qaytarish qiymatlari va ko'p qaytarishlar

2️⃣ Ma'lumot tuzilmalari
   - Massivlar (o'zgarmas o'lcham) va slayslar (dinamik o'lcham)
   - Map (xaritalar) - kalit/qiymat juftliklari bilan ishlash
   - Strukturalar (struct) va ularning usullari (methods)

3️⃣ Dastur oqimi boshqaruvi
   - If/else shartli ifodalar
   - For looplarining turli ko'rinishlari
   - Switch va select ifodalar

4️⃣ Paralel dasturlash asoslari
   - Goroutine - Go'ning engil vazn thread'lari
   - Kanallar (channel) orqali ma'lumot almashish
   - Sync paketi va mutex yordamida sinxronizatsiya

5️⃣ Interfeys va xatolar bilan ishlash
   - Interfeys tushunchasi va duck typing
   - Xatolarni qayta ishlash metodologiyasi
   - defer, panic va recover mexanizmlari

6️⃣ Testlash va sifat ta'minoti
   - Go texnologiyasida yozilgan unit testlar
   - Benchmark test'lar orqali samaradorlikni baholash
   - Table-driven test usuli

7️⃣ Paketlar va modullar tizimi
   - Go module tizimi va go.mod fayli
   - Paket strukturasi va importlar
   - Eksport (bosh harf) va shaxsiy (kichik harf) identifikatorlar

8️⃣ Ilg'or mavzular
   - Reflection mexanizmi bilan ishlash
   - CGO - C kodini Go bilan integratsiyalash
   - Context paketi va uni qo'llash usullari

9️⃣ Amaliy loyihalar
   - CLI (buyruq qatori) dasturlari yaratish
   - Web xizmatlar va HTTP server (net/http)
   - Ma'lumotlar bazasi bilan ishlash (SQL va NoSQL)

Boshlash uchun eng yaxshi resurs: https://go.dev/learn/`
}

// GetUsefulText Go dasturlash tili bo'yicha foydali resurslarni qaytaradi
// Bu funksiya Go o'rganish uchun eng yaxshi manbalar to'plamini taqdim etadi
func (h *CommandHandler) GetUsefulText() string {
	return `Go dasturlash tili bo'yicha eng foydali manbalar:

📚 Asosiy manbalar:
- Rasmiy veb-sayt: https://go.dev - barcha rasmiy hujjatlar va yangiliklar
- Tour of Go: https://tour.golang.org/ - interaktiv o'rganish qo'llanmasi
- Go by Example: https://gobyexample.com/ - misollarda Go'ni o'rganish
- Effektiv Go: https://go.dev/doc/effective_go - samarali kod yozish bo'yicha tavsiyalar
- Standard Library: https://pkg.go.dev/std - standart kutubxonalar hujjatlari
- Go Playground: https://play.golang.org/ - brauzerda kod yozish va sinab ko'rish
- Go hamjamiyati blogi: https://go.dev/blog/ - yangiliklar va chuqurlashtirilgan maqolalar

🔍 Qo'shimcha foydali manbalar:
- Awesome Go: https://github.com/avelino/awesome-go - Go kutubxonalari va vositalar to'plami
- Go Design Patterns: https://github.com/tmrts/go-patterns - Go uchun dizayn patternlar
- Go Forums: https://forum.golangbridge.org/ - savol-javoblar va muhokamalar

📖 Tavsiya etiladigan kitoblar:
- "The Go Programming Language" - Alan Donovan va Brian Kernighan
- "Go in Action" - William Kennedy
- "Concurrency in Go" - Katherine Cox-Buday

🎓 Video darslar:
- Golang bo'yicha o'zbek tilidagi darslar: https://youtube.com/playlist?list=PLLIX7niqDict7oqNQesQQT9b3GlqF7JAj`
}

// GetLatestText eng yangi Go versiyasi haqida ma'lumot qaytaradi
// Bu funksiya Go'ning eng so'nggi versiyasidagi asosiy yangiliklar haqida ma'lumot beradi
func (h *CommandHandler) GetLatestText() string {
	// Haqiqiy dasturda bu ma'lumotlar Go rasmiy saytidan yoki API orqali olinishi mumkin
	return `Go 1.22.1 (2024-yil 5-mart) versiyasidagi asosiy yangiliklar:

🔧 Muhim xatolar tuzatildi:
- net/http paketidagi xavfsizlik bilan bog'liq muammolar bartaraf etildi
- crypto paketlaridagi xatolar tuzatildi

🚀 Yaxshilanishlar:
- Kompilyator ishlash tezligi oshirildi
- Runtime samaradorligi yaxshilandi
- Xotiradan foydalanish optimallashtirildi

🔒 Xavfsizlik yangiliklari:
- Standart kutubxonalardagi potensial zaifliklar bartaraf etildi

Batafsil ma'lumot: https://go.dev/doc/devel/release#go1.22.1`
}

// GetVersionText ko'rsatilgan Go versiyasi haqida batafsil ma'lumot qaytaradi
// Bu funksiya foydalanuvchi so'ragan versiya haqida to'liq ma'lumotni beradi
func (h *CommandHandler) GetVersionText(version string) string {
	// Agar versiya ko'rsatilmagan bo'lsa, foydalanuvchiga ko'rsatma berish
	if version == "" {
		return "Iltimos, ma'lumot olmoqchi bo'lgan versiya raqamini kiriting. Masalan: /version 1.22.0"
	}

	// Haqiqiy dasturda bu ma'lumotlar ma'lumotlar bazasi yoki API orqali olinishi kerak
	versions := map[string]string{
		"1.22.1": "Go 1.22.1 (2024-yil 5-mart):\n\n" +
			"🔧 Xatoliklar tuzatishlari:\n" +
			"- net/http paketidagi HTTP sarlavhalarni qayta ishlashdagi xatolar tuzatildi\n" +
			"- crypto/tls paketidagi sertifikat tekshirishda optimizatsiyalar qilindi\n" +
			"- reflect paketidagi xotira sizishlar bartaraf etildi\n\n" +
			"🚀 Yaxshilanishlar:\n" +
			"- Paralel garbage collection algoritmi takomillashtirildi\n" +
			"- GOEXPERIMENT=rangefunc bayroq orqali yangi range funksiyalarini sinash imkoniyati qo'shildi",

		"1.22.0": "Go 1.22.0 (2024-yil 6-fevral):\n\n" +
			"🆕 Yangi imkoniyatlar:\n" +
			"- Butun sonlar ustida iteratsiya qilish uchun yangi range sintaksisi (range 10)\n" +
			"- HTTP router pattern matching qo'llab-quvvatlash bilan yaxshilandi\n" +
			"- For loop'larda xatoliklarni qayta ishlash takomillashtirildi\n\n" +
			"🔄 Muhim o'zgarishlar:\n" +
			"- Orqaga moslik yanada kuchaytirildi\n" +
			"- Xatolik xabarlari tushunarliroq bo'ldi\n" +
			"- Paket importi optimallashtirildi",

		"1.21.0": "Go 1.21.0 (2023-yil 8-avgust):\n\n" +
			"🆕 Yangi imkoniyatlar:\n" +
			"- min() va max() o'rnatilgan funksiyalar qo'shildi\n" +
			"- Loop o'zgaruvchilari semantikasi o'zgartirildi (har bir iteratsiya uchun yangi o'zgaruvchi)\n" +
			"- slog paketi orqali strukturaviy log yozish imkoniyati qo'shildi\n\n" +
			"🔧 Yaxshilanishlar:\n" +
			"- Forward compatible method chaqirishlari\n" +
			"- PGO (Profile-guided optimization) orqali dastur ishlash tezligini oshirish",
	}

	// Versiya topilganda, uni batafsil ma'lumot bilan qaytarish
	if info, ok := versions[version]; ok {
		return info + "\n\nRasmiy hujjatlar va batafsilroq ma'lumot: https://go.dev/doc/devel/release#go" + strings.ReplaceAll(version, ".", "")
	}

	// Agar versiya topilmasa, xato xabarini qaytarish
	return fmt.Sprintf("Kechirasiz, %s versiyasi haqida ma'lumot bazamizda topilmadi. Mavjud versiyalar: 1.21.0, 1.22.0, 1.22.1", version)
}

// GetWarnText ko'rsatilgan foydalanuvchiga ogohlantirish xabari qaytaradi
// Bu funksiya guruh qoidalarini buzgan foydalanuvchilarni ogohlantirish uchun ishlatiladi
func (h *CommandHandler) GetWarnText(username string) string {
	// Agar foydalanuvchi nomi mavjud bo'lmasa, umumiy murojaat ishlatish
	if username == "" {
		username = "Foydalanuvchi"
	}
	// Ogohlantirish xabarini formatlash va qaytarish
	return fmt.Sprintf("⚠️ Diqqat @%s! Iltimos, guruh qoidalariga rioya qiling va mavzudan chetlashmang. Qoidalar bilan tanishish uchun /rules buyrug'ini yuboring.", username)
}
