package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	tg "github.com/amarnathcjd/gogram/telegram"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/ulugbek0217/sheriff-bot/db/sqlc"
	"github.com/ulugbek0217/sheriff-bot/handlers"
	"github.com/ulugbek0217/sheriff-bot/utils"
)

func main() {
	err := utils.LoadEnv("config/.env")
	if err != nil {
		log.Fatalf("error loading env: %v", err)
	}

	app_id, err := strconv.ParseInt(os.Getenv("APP_ID"), 10, 64)
	if err != nil {
		log.Fatalf("error reading env: %v", err)
	}
	app_hash := os.Getenv("APP_HASH")
	phone := os.Getenv("PHONE")

	client, err := tg.NewClient(tg.ClientConfig{
		AppID:   int32(app_id),
		AppHash: app_hash,
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

	client.Login(phone, &tg.LoginOptions{
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

	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		os.Exit(1)
	}
	defer dbPool.Close()
	// conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	// 	os.Exit(1)
	// }
	// defer conn.Close(context.Background())

	var wg = sync.WaitGroup{}
	app := &handlers.App{
		Store:  db.NewStore(dbPool),
		Pool:   dbPool,
		Admins: utils.GetAdmins(),
		WG:     &wg,
	}

	client.SendMessage(utils.GetAdmins()[0], "Bot has been started")

	client.On("message:/export", app.ExportAccounts, tg.FilterPrivate, utils.FilterAdmins)
	client.On("message:/info", app.GetUserInfo, tg.FilterPrivate, utils.FilterAdmins)
	client.On(tg.OnNewMessage, app.CollectUserInfo, tg.FilterGroup)

	log.Println("Bot is running")
	client.Idle()
}
