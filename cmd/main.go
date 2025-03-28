package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/nktauserum/aisearch-telegram/internal/app"
)

func main() {
	godotenv.Load(".env")
	app := app.NewApplication(os.Getenv("TG_TOKEN"))

	err := app.Run()
	if err != nil {
		log.Fatalf("error while running the application: %v", err)
	}
}
