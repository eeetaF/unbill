package vacuum_cleaner

import (
	"log"
	"os"
	"path/filepath"
	"time"

	runtime_data "backend/runtime-data"
)

func AutoClearOldFiles(maxAge time.Duration, interval time.Duration) {
	dir := runtime_data.GetConfig().DataDir
	for {
		files, err := os.ReadDir(dir)
		if err != nil {
			log.Printf("Auto-clear: failed to read dir %s: %v", dir, err)
			time.Sleep(interval)
			continue
		}

		now := time.Now()
		for _, f := range files {
			if f.IsDir() {
				continue
			}

			path := filepath.Join(dir, f.Name())
			info, err := f.Info()
			if err != nil {
				log.Printf("Auto-clear: failed to get info for %s: %v", path, err)
				continue
			}

			if now.Sub(info.ModTime()) > maxAge {
				err = os.Remove(path)
				if err != nil {
					log.Printf("Auto-clear: failed to remove %s: %v", path, err)
				} else {
					log.Printf("Auto-clear: removed %s", path)
				}
			}
		}

		time.Sleep(interval)
	}
}
