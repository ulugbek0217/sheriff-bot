package handlers

import (
	"context"
	"fmt"
	"log"
	"math"

	tg "github.com/amarnathcjd/gogram/telegram"
	db "github.com/ulugbek0217/sheriff-bot/db/sqlc"
	"github.com/ulugbek0217/sheriff-bot/utils"
)

func (app *App) ExportAccounts(message *tg.NewMessage) error {
	accountsQuantity, err := app.Store.AccountsQuantity(context.Background())
	if err != nil {
		log.Fatalf("error getting accounts quantity: %v", err)
	}
	for i := 1; i <= int(math.Ceil(float64(accountsQuantity)/50.0)); i++ {
		arg := db.ListAccountsParams{
			Limit:  int32((50)),
			Offset: int32(i - 1),
		}

		accounts, err := app.Store.ListAccounts(context.Background(), arg)
		if err != nil {
			message.Reply(fmt.Errorf("error listing accounts: %v", err))
			return err
		}

		err = utils.ExportToExcel(accounts)
		if err != nil {
			message.Reply(fmt.Errorf("error exporting to excel: %v", err))
			return err
		}
	}

	_, err = message.RespondMedia("files/accounts.xlsx")
	if err != nil {
		message.Reply(fmt.Errorf("error sending exported accounts: %v", err))
		return err
	}

	return nil
}
