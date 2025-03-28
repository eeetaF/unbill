package runtime_data

import (
	"log"
	"time"

	"backend/vacuum_cleaner"
)

func InitializeApp() {
	if err := LoadConfig("app.ini"); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := GetConfig().InitKey(); err != nil {
		log.Fatalf("Cannot initialize program: InitKey: %v", err)
	}

	go vacuum_cleaner.AutoClearOldFiles(GetConfig().MaxAge, max(GetConfig().MaxAge/24, time.Second))
}
