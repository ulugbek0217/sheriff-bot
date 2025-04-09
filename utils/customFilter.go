package utils

import (
	tg "github.com/amarnathcjd/gogram/telegram"
)

var FilterAdmins = tg.FilterFunc(func(m *tg.NewMessage) bool {
	admins := GetAdmins()
	senderID := m.Sender.ID

	for _, id := range admins {
		if id == senderID {
			return true
		}
	}
	return false
})
