package utils

import (
	"os"
	"path/filepath"
)

func CreateClassNameHash(dirs []string) map[string]string {
	res := make(map[string]string)

	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, file := range files {
			path := filepath.Join(dir, file.Name())

			data, err := GetMap(path, "Desktop Entry")
			if err != nil {
				continue
			}

			general, exist := data["Desktop Entry"]
			if !exist {
				continue
			}

			className, exist := general["StartupWMClass"]
			if !exist {
				continue
			}

			if className != "" {
				res[className] = path
			}
		}
	}

	return res
}
