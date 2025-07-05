package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var userState = make(map[int64]string)
var userData = make(map[int64]map[string]string)

var srv *sheets.Service

// Deklarasikan variabel global untuk diisi dari environment
var BotToken string
var spreadsheetPLN string
var spreadsheetPLTS string
var spreadsheetGenset125 string
var spreadsheetGenset400 string
var googleServiceAccountKeyPath string

// --- NEW: Tambahkan variabel untuk link masing-masing spreadsheet ---
var linkSpreadsheetPLN string
var linkSpreadsheetPLTS string
var linkSpreadsheetGenset125 string
var linkSpreadsheetGenset400 string

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	BotToken = os.Getenv("BOT_TOKEN")
	if BotToken == "" {
		log.Fatal("BOT_TOKEN environment variable not set")
	}

	spreadsheetPLN = os.Getenv("SPREADSHEET_PLN")
	if spreadsheetPLN == "" {
		log.Fatal("SPREADSHEET_PLN environment variable not set")
	}

	spreadsheetPLTS = os.Getenv("SPREADSHEET_PLTS")
	if spreadsheetPLTS == "" {
		log.Fatal("SPREADSHEET_PLTS environment variable not set")
	}

	spreadsheetGenset125 = os.Getenv("SPREADSHEET_GENSET_125")
	if spreadsheetGenset125 == "" {
		log.Fatal("SPREADSHEET_GENSET_125 environment variable not set")
	}

	spreadsheetGenset400 = os.Getenv("SPREADSHEET_GENSET_400")
	if spreadsheetGenset400 == "" {
		log.Fatal("SPREADSHEET_GENSET_400 environment variable not set")
	}

	googleServiceAccountKeyPath = os.Getenv("GOOGLE_SERVICE_ACCOUNT_KEY_PATH")
	if googleServiceAccountKeyPath == "" {
		log.Fatal("GOOGLE_SERVICE_ACCOUNT_KEY_PATH environment variable not set")
	}

	// --- NEW: Inisialisasi link spreadsheet dari environment ---
	linkSpreadsheetPLN = os.Getenv("LINK_SPREADSHEET_PLN")
	if linkSpreadsheetPLN == "" {
		log.Fatal("LINK_SPREADSHEET_PLN environment variable not set")
	}
	linkSpreadsheetPLTS = os.Getenv("LINK_SPREADSHEET_PLTS")
	if linkSpreadsheetPLTS == "" {
		log.Fatal("LINK_SPREADSHEET_PLTS environment variable not set")
	}
	linkSpreadsheetGenset125 = os.Getenv("LINK_SPREADSHEET_GENSET_125")
	if linkSpreadsheetGenset125 == "" {
		log.Fatal("LINK_SPREADSHEET_GENSET_125 environment variable not set")
	}
	linkSpreadsheetGenset400 = os.Getenv("LINK_SPREADSHEET_GENSET_400")
	if linkSpreadsheetGenset400 == "" {
		log.Fatal("LINK_SPREADSHEET_GENSET_400 environment variable not set")
	}
}

func initGoogleSheet() {
	b, err := os.ReadFile(googleServiceAccountKeyPath)
	if err != nil {
		log.Fatalf("Gagal membaca kredensial dari %s: %v", googleServiceAccountKeyPath, err)
	}

	config, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Gagal parsing JSON: %v", err)
	}

	client := config.Client(context.Background())
	srv, err = sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Tidak bisa membuat service Sheets: %v", err)
	}
}

func appendToSheet(spreadsheetID string, data []interface{}) error {
	vr := &sheets.ValueRange{
		Values: [][]interface{}{data},
	}
	_, err := srv.Spreadsheets.Values.Append(spreadsheetID, "Sheet1!A1", vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Printf("Gagal menulis ke Google Sheet: %v", err)
		return err
	}
	return nil
}

func main() {
	initGoogleSheet()

	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false
	log.Printf("Bot berhasil dijalankan")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		userID := update.Message.From.ID
		text := update.Message.Text

		if _, ok := userData[userID]; !ok {
			userData[userID] = make(map[string]string)
		}

		switch userState[userID] {
		case "":
			if strings.HasPrefix(text, "/start") {
				msgText := "⚡️ *Mulai Pengecekan Catu Daya* ⚡️\n\n" +
					"Silakan pilih sumber catu daya yang akan diperiksa:\n\n" +
					"1. 🏢  PLN\n" +
					"2. ☀️  PLTS\n" +
					"3. ⛽️  Genset 125 kVA\n" +
					"4. ⛽️  Genset 400 kVA\n" +
					"5. 📖  Lihat Rekap Report\n\n" +
					"_Ketik angka pilihan Anda (contoh: 2)_"

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				msg.ParseMode = "Markdown"
				bot.Send(msg)

				userState[userID] = "JENIS"
			}

		case "JENIS":
			switch text {
			case "1":
				userData[userID]["alat"] = "PLN"
				userState[userID] = "VOLT"
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "💡 Masukkan data voltase (Volt):"))
			case "2":
				userData[userID]["alat"] = "PLTS"
				userState[userID] = "VOLT"
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "💡 Masukkan data voltase (Volt):"))
			case "3":
				userData[userID]["alat"] = "Genset 125 kVA"
				userState[userID] = "VOLT"
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "💡 Masukkan data voltase (Volt):"))
			case "4":
				userData[userID]["alat"] = "Genset 400 kVA"
				userState[userID] = "VOLT"
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "💡 Masukkan data voltase (Volt):"))
			case "5":
				// --- NEW: Pilihan untuk melihat rekap report ---
				msgText := "Pilih rekap report yang ingin Anda lihat:\n\n" +
					"1. 🏢  Rekap PLN\n" +
					"2. ☀️  Rekap PLTS\n" +
					"3. ⛽️  Rekap Genset 125 kVA\n" +
					"4. ⛽️  Rekap Genset 400 kVA\n\n" +
					"_Ketik angka pilihan Anda (contoh: 1)_"

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
				bot.Send(msg)
				userState[userID] = "VIEW_REPORT_CHOICE" // Set state baru
			default:
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Pilihan tidak valid. Ketik 1–5."))
				continue
			}

		// --- NEW: State untuk menangani pilihan rekap report ---
		case "VIEW_REPORT_CHOICE":
			var reportLink string
			var reportName string
			switch text {
			case "1":
				reportLink = linkSpreadsheetPLN
				reportName = "Rekap PLN"
			case "2":
				reportLink = linkSpreadsheetPLTS
				reportName = "Rekap PLTS"
			case "3":
				reportLink = linkSpreadsheetGenset125
				reportName = "Rekap Genset 125 kVA"
			case "4":
				reportLink = linkSpreadsheetGenset400
				reportName = "Rekap Genset 400 kVA"
			default:
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Pilihan tidak valid. Ketik 1–4 untuk melihat rekap report."))
				continue
			}

			userState[userID] = "" // Reset state setelah memberikan link

			linkMessage := fmt.Sprintf("Anda bisa melihat %s di sini:", reportName)

			// Buat inline keyboard dengan tombol URL
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonURL("Buka Spreadsheet", reportLink),
				),
			)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, linkMessage)
			msg.ReplyMarkup = keyboard // Lampirkan keyboard
			bot.Send(msg)

			// Opsional: Kirim ulang menu utama setelah memberikan link
			msgText := "⚡️ *Mulai Pengecekan Catu Daya* ⚡️\n\n" +
				"Silakan pilih sumber catu daya yang akan diperiksa:\n\n" +
				"1. 🏢  PLN\n" +
				"2. ☀️  PLTS\n" +
				"3. ⛽️  Genset 125 kVA\n" +
				"4. ⛽️  Genset 400 kVA\n" +
				"5. 📖  Lihat Rekap Report\n\n" +
				"_Ketik angka pilihan Anda (contoh: 2)_"

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
			msg.ParseMode = "Markdown"
			bot.Send(msg)

		case "VOLT":
			userData[userID]["volt"] = text
			userState[userID] = "ARUS"
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "💡 Masukkan data arus (Ampere):"))

		case "ARUS":
			userData[userID]["arus"] = text
			userState[userID] = "JAM"
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "💡 Masukkan jam pengecekan (format hh:mm):"))

		case "JAM":
			userData[userID]["jam"] = text
			userState[userID] = "KONDISI"
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "💡 Bagaimana kondisi peralatan?\n1. Normal\n2. Tidak Normal / Rusak"))

		case "KONDISI":
			if text == "1" {
				userData[userID]["kondisi"] = "Normal"
			} else {
				userData[userID]["kondisi"] = "Tidak Normal / Rusak"
			}
			userState[userID] = "" // Reset state pengguna

			data := userData[userID]
			hasil := fmt.Sprintf(
				"🧾 *Rekap Data Pengecekan*\n\n"+
					"*Catu Daya:* %s\n"+
					"*Voltase:* %s V\n"+
					"*Arus:* %s A\n"+
					"*Jam:* %s\n"+
					"*Kondisi:* %s\n\n"+
					"Mohon tunggu, sedang menyimpan data...",
				data["alat"], data["volt"], data["arus"], data["jam"], data["kondisi"],
			)
			recapMsg := tgbotapi.NewMessage(update.Message.Chat.ID, hasil)
			recapMsg.ParseMode = "Markdown"
			bot.Send(recapMsg)

			var targetSpreadsheet string
			switch data["alat"] {
			case "PLN":
				targetSpreadsheet = spreadsheetPLN
			case "PLTS":
				targetSpreadsheet = spreadsheetPLTS
			case "Genset 125 kVA":
				targetSpreadsheet = spreadsheetGenset125
			case "Genset 400 kVA":
				targetSpreadsheet = spreadsheetGenset400
			}

			timestamp := time.Now().Format("2006-01-02 15:04:05")
			nama := update.Message.From.FirstName
			err := appendToSheet(targetSpreadsheet, []interface{}{timestamp, nama, data["alat"], data["volt"], data["arus"], data["jam"], data["kondisi"]})

			var alertMsg tgbotapi.MessageConfig
			if err != nil {
				alertMsg = tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Gagal menyimpan data ke spreadsheet. Silakan hubungi admin.")
			} else {
				alertMsg = tgbotapi.NewMessage(update.Message.Chat.ID, "✅ Data berhasil disimpan ke Database!")
			}
			bot.Send(alertMsg)
		}
	}
}
