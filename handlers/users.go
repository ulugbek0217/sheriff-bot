package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	tg "github.com/amarnathcjd/gogram/telegram"
	db "github.com/ulugbek0217/sheriff-bot/db/sqlc"
)

type App struct {
	Store  db.Store
	Admins []int64
}

type UserInfo struct {
	UserID            int64  `json:"user_id"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Username          string `json:"username"`
	Phone             string `json:"phone"`
	About             string `json:"about"`
	Birthday          string `json:"birthday"`
	PersonalChannelID string `json:"personal_channel_id"`
}

func (app *App) CollectUserInfo(message *tg.NewMessage) error {
	// Check if the sender is not a bot and not myself
	me, _ := message.Client.GetMe()
	if message.Sender.Bot || message.Sender.ID == me.ID {
		return nil
	}

	user, err := message.Client.UsersGetFullUser(&tg.InputUserObj{UserID: message.SenderID(), AccessHash: message.Sender.AccessHash})
	if err != nil {
		log.Printf("error UsersGetFullUser: %v", err)
	}
	fu := user.FullUser
	if fu == nil {
		fmt.Println("cannot get user")
		return err
	}

	var birthday string
	if fu.Birthday != nil {
		birthday = fmt.Sprintf("%02d-%02d-%d", fu.Birthday.Day, fu.Birthday.Month, fu.Birthday.Year)
	} else {
		birthday = "null"
	}

	var personalChannelID string = fmt.Sprintf("%d", fu.PersonalChannelID)

	var userInfo UserInfo = UserInfo{
		UserID:            message.Sender.ID,
		FirstName:         message.Sender.FirstName,
		LastName:          message.Sender.LastName,
		Username:          message.Sender.Username,
		Phone:             message.Sender.Phone,
		About:             fu.About,
		Birthday:          birthday,
		PersonalChannelID: personalChannelID,
	}

	// b, err := json.MarshalIndent(userInfo, "", "\t")
	// if err != nil {
	// 	log.Fatalf("cannot marshal to json: %v", err)
	// }

	// fmt.Println(string(b))

	err = app.Store.CreateAccount(context.Background(), db.CreateAccountParams{
		UserID:            userInfo.UserID,
		FirstName:         userInfo.FirstName,
		LastName:          userInfo.LastName,
		Username:          userInfo.Username,
		Phone:             userInfo.Phone,
		About:             userInfo.About,
		Birthday:          userInfo.Birthday,
		PersonalChannelID: userInfo.PersonalChannelID,
	})

	if err != nil {
		log.Printf("cannot create account: %v", err)
		return err
	}

	return nil
}

func (app *App) GetUserInfo(message *tg.NewMessage) error {
	// fmt.Println("get user info handler")
	args := message.Args()
	// fmt.Println(args)
	userID, err := strconv.ParseInt(args, 10, 64)
	if err != nil {
		log.Printf("err parsing userID from args: %v", err)
		return err
	}
	user, err := app.Store.GetAccount(context.Background(), userID)
	if err != nil {
		fmt.Printf("err getting user from db: %v", err)
		return err
	}

	b, err := json.MarshalIndent(user, "", "\t")
	if err != nil {
		log.Fatalf("cannot marshal to json: %v", err)
		return err
	}

	_, err = message.Reply(`<pre lang="json">`+string(b)+`</pre>`, tg.SendOptions{ParseMode: "html"})
	if err != nil {
		fmt.Printf("err sending info: %v", err)
	}

	return nil
}
