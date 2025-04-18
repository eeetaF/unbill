package initialization

import (
	"log"
	"time"

	data "backend/config"
	"backend/vacuum_cleaner"
)

func InitializeApp() {
	if err := data.LoadConfig("app.ini"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := data.GetConfig().InitKey(); err != nil {
		log.Fatalf("Cannot initialize program: InitKey: %v", err)
	}

	go vacuum_cleaner.AutoClearOldFiles(data.GetConfig().MaxAge, max(data.GetConfig().MaxAge/24, time.Second))
}
