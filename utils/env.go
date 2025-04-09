package utils

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func LoadEnv(path string) error {
	err := godotenv.Overload(path)
	return err
}

func GetAdmins() []int64 {
	admins_env := os.Getenv("ADMINS")
	admins_str := strings.Split(admins_env, ",")
	var admins []int64
	for _, i := range admins_str {
		admin_id, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			log.Fatalf("cannot parse admin ids: %v", err)
		}
		admins = append(admins, admin_id)
	}
	return admins
}
