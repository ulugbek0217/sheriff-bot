package utils

import (
	"fmt"
	"os"
	"strconv"

	db "github.com/ulugbek0217/sheriff-bot/db/sqlc"
	"github.com/xuri/excelize/v2"
)

func ExportToExcel(accounts []db.Account) error {
	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", "user_id")
	f.SetCellValue("Sheet1", "B1", "first_name")
	f.SetCellValue("Sheet1", "C1", "last_name")
	f.SetCellValue("Sheet1", "D1", "username")
	f.SetCellValue("Sheet1", "E1", "phone")
	f.SetCellValue("Sheet1", "F1", "about")
	f.SetCellValue("Sheet1", "G1", "birthday")
	f.SetCellValue("Sheet1", "H1", "peronal_channel_id")
	f.SetCellValue("Sheet1", "I1", "created_at")
	for i := 0; i < len(accounts); i++ {
		f.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), accounts[i].UserID)
		f.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), accounts[i].FirstName)
		f.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), accounts[i].LastName)
		f.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), accounts[i].Username)
		f.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), accounts[i].Phone)
		f.SetCellValue("Sheet1", "F"+strconv.Itoa(i+2), accounts[i].About)
		f.SetCellValue("Sheet1", "G"+strconv.Itoa(i+2), accounts[i].Birthday)
		f.SetCellValue("Sheet1", "H"+strconv.Itoa(i+2), accounts[i].PersonalChannelID)
		f.SetCellValue("Sheet1", "I"+strconv.Itoa(i+2), accounts[i].CreatedAt)
	}
	var dir string = "files"
	if !FolderExists(dir) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return fmt.Errorf("error creating files folder: %v", err)
		}
	}

	if err := f.SaveAs(dir + "/accounts.xlsx"); err != nil {
		return fmt.Errorf("eror saving to .xlsx: %v", err)
	}
	return nil
}
