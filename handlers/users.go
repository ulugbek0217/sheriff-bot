package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	tg "github.com/amarnathcjd/gogram/telegram"
	"github.com/jackc/pgx/v5"
	db "github.com/ulugbek0217/sheriff-bot/db/sqlc"
)

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
	app.WG.Add(1)
	defer app.WG.Done()

	// Check if the sender is not a bot and not myself
	me, _ := message.Client.GetMe()
	if message.Sender.Bot || message.Sender.ID == me.ID || message.SenderID() == message.ChannelID() {
		return nil
	}

	user, err := message.Client.UsersGetFullUser(&tg.InputUserObj{UserID: message.SenderID(), AccessHash: message.Sender.AccessHash})
	if err != nil {
		return fmt.Errorf("error getting UsersGetFullUser: %v", err)
	}

	var fu *tg.UserFull = user.FullUser
	var un = user.Users[0].(*tg.UserObj)

	var userInfo db.CreateAccountParams = db.CreateAccountParams{}

	userInfo.UserID = un.ID
	userInfo.FirstName = un.FirstName
	userInfo.LastName = un.LastName
	userInfo.Username = un.Username
	userInfo.Phone = un.Phone

	if fu.Birthday != nil { // Safe check
		userInfo.Birthday = fmt.Sprintf("%02d-%02d-%d", fu.Birthday.Day, fu.Birthday.Month, fu.Birthday.Year)
	} else {
		userInfo.Birthday = "null"
	}

	userInfo.PersonalChannelID = fmt.Sprintf("%d", fu.PersonalChannelID)

	err = app.Store.CreateAccount(context.Background(), app.Pool, userInfo)

	if err != nil {
		log.Printf("cannot create account: %v", err)
		return err
	}

	return nil
}

func (app *App) GetUserInfo(message *tg.NewMessage) error {
	var userID int64 = 0
	var userHash int64 = 0

	if len(message.Args()) > 0 {
		i, ok := strconv.Atoi(message.Args())
		if ok != nil {
			user, err := message.Client.ResolveUsername(message.Args())
			if err != nil {
				message.Reply("Error: " + err.Error())
				return nil
			}
			ux, ok := user.(*tg.UserObj)
			if !ok {
				message.Reply("Error: User not found")
				return nil
			}
			userID = ux.ID
			userHash = ux.AccessHash
		} else {
			userID = int64(i)
			user, err := message.Client.GetUser(int64(i))
			if err != nil {
				message.Reply("Error: " + err.Error())
				return nil
			}
			userHash = user.AccessHash
		}
	}
	user, err := message.Client.UsersGetFullUser(&tg.InputUserObj{
		UserID:     userID,
		AccessHash: userHash,
	})
	if err != nil {
		message.Reply("Error: " + err.Error())
		return nil
	}

	var un = user.Users[0].(*tg.UserObj)

	app.WG.Add(1)
	defer app.WG.Done()

	// args := message.Args()

	// userID, err := strconv.ParseInt(args, 10, 64)
	// if err != nil {
	// 	message.Reply("Invalid ID")
	// 	return fmt.Errorf("err parsing userID from args: %v", err)
	// }
	result, err := app.Store.GetAccount(context.Background(), app.Pool, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			message.Reply("No user with such ID")
			return nil
		}
		return fmt.Errorf("err getting user from db: %v", err)
	}

	b, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		return fmt.Errorf("cannot marshal to json: %v", err)
	}

	var keyb = tg.NewKeyboard()
	sendableUser, err := message.Client.GetSendableUser(un)
	if err == nil {
		keyb.AddRow(
			tg.Button.Mention("Go >> User Profile", sendableUser),
		)
	} else {
		keyb.AddRow(
			tg.Button.URL("Go >> User Profile", "tg://user?id="+strconv.FormatInt(un.ID, 10)),
		)
	}

	_, err = message.Reply(`<pre lang="json">`+string(b)+`</pre>`, tg.SendOptions{
		ParseMode:   "html",
		ReplyMarkup: keyb.Build(),
	})
	if err != nil {
		return fmt.Errorf("err sending info: %v", err)
	}

	return nil
}
