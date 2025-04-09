package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tg "github.com/amarnathcjd/gogram/telegram"
	"github.com/jackc/pgx/v5"
	db "github.com/ulugbek0217/sheriff-bot/db/sqlc"
	"github.com/ulugbek0217/sheriff-bot/handlers"
	"github.com/ulugbek0217/sheriff-bot/utils"
)

func main() {
	client, err := tg.NewClient(tg.ClientConfig{
		AppID:   19973721,
		AppHash: "a0a5674cceaa4283aa00ee243f6089b8",
		DeviceConfig: tg.DeviceConfig{
			DeviceModel:   "Sheriff",
			SystemVersion: "Linux",
			AppVersion:    "1.0",
		},
		ParseMode:  "html",
		DataCenter: 2,
		LogLevel:   tg.LogInfo,
		Session:    "session.dat",
	})

	if err != nil {
		log.Fatalf("error creating client: %v", err)
	}

	// Connect to Telegram
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	client.Login("+998338943615", &tg.LoginOptions{
		// Handle OTP code callback
		CodeCallback: func() (string, error) {
			fmt.Print("Enter Telegram code: ")
			var code string
			fmt.Scanln(&code)
			return code, nil
		},
		// Handle 2FA password callback (if enabled)
		PasswordCallback: func() (string, error) {
			fmt.Print("Enter 2FA password: ")
			var password string
			fmt.Scanln(&password)
			return password, nil
		},
	})

	err = client.Start()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Disconnect()

	err = utils.LoadEnv("config/.env")
	if err != nil {
		log.Fatalf("error loading env: %v", err)
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	app := &handlers.App{
		Store:  db.NewStore(conn),
		Admins: utils.GetAdmins(),
	}

	client.SendMessage(utils.GetAdmins()[0], "Bot has been started")

	client.On("message:/export", app.ExportAccounts, tg.FilterPrivate, utils.FilterAdmins)
	client.On("message:/info", app.GetUserInfo, tg.FilterPrivate, utils.FilterAdmins)
	client.On(tg.OnNewMessage, app.CollectUserInfo, tg.FilterGroup)

	log.Println("Bot is running")
	client.Idle()
}
