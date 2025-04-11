package handlers

import (
	"context"
	"fmt"
	"math"
	"sync"

	tg "github.com/amarnathcjd/gogram/telegram"
	db "github.com/ulugbek0217/sheriff-bot/db/sqlc"
	"github.com/ulugbek0217/sheriff-bot/utils"
)

func (app *App) ExportAccounts(message *tg.NewMessage) error {
	app.WG.Add(1)
	defer app.WG.Done()
	accountsQuantity, err := app.Store.AccountsQuantity(context.Background(), app.Pool)
	if err != nil {
		return fmt.Errorf("error getting accounts quantity: %v", err)
	}

	exportWG := sync.WaitGroup{}

	pages := int(math.Ceil(float64(accountsQuantity) / 50.0))
	accounts := make(chan []db.Account, pages)
	errors := make(chan error, pages)
	for i := 1; i <= pages; i++ {
		exportWG.Add(1)
		go func(pageID int) {
			defer exportWG.Done()

			arg := db.ListAccountsParams{
				Limit:  int32((50)),
				Offset: int32(pageID-1) * 50,
			}

			accountsPerPage, err := app.Store.ListAccounts(context.Background(), app.Pool, arg)
			if err != nil {
				errors <- fmt.Errorf("error listing accounts on page %d: %v", pageID, err)
				return
			}
			accounts <- accountsPerPage

		}(i)
	}

	go func() {
		exportWG.Wait()
		close(accounts)
		close(errors)
	}()

	var accountsList = []db.Account{}

	for accList := range accounts {
		accountsList = append(accountsList, accList...)
	}
	// action, _ := message.SendAction("sending file")
	// defer action.Cancel()

	err = utils.ExportToExcel(accountsList)
	if err != nil {
		message.Reply("error exporting to excel: %v")
		return fmt.Errorf("error exporting to excel: %v", err)
	}

	_, err = message.RespondMedia("files/accounts.xlsx")
	if err != nil {
		message.Reply(fmt.Errorf("error sending exported accounts: %v", err))
		return err
	}

	return nil
}
