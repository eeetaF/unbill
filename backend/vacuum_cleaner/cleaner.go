package vacuum_cleaner

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"backend/config"
)

func AutoClearOldFiles(maxAge time.Duration, interval time.Duration) {
	dir1, dir2 := config.GetConfig().DataSharedDir, config.GetConfig().DataSplitDir
	for {
		clearDir(dir1, interval, maxAge)
		clearDir(dir2, interval, maxAge)

		time.Sleep(interval)
	}
}

func clearDir(dir string, interval time.Duration, maxAge time.Duration) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("Auto-clear: failed to read dir %s: %v", dir, err)
		return
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
			return
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
}
