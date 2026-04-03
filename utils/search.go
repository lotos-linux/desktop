package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func SearchDesktopFile(className string, dir string) string {
	_, err := os.Stat(filepath.Join(dir, className+".desktop"))
	if err == nil {
		return filepath.Join(dir, className+".desktop")
	}

	// If file non found
	files, _ := os.ReadDir(dir)

	// "krita" > "org.kde.krita.desktop" / "lutris" > "net.lutris.Lutris.desktop"
	for _, file := range files {
		if strings.Count(file.Name(), ".") > 1 && strings.Contains(strings.ToLower(file.Name()), className) {
			return filepath.Join(dir, file.Name())
		}

	}

	// "VirtualBox Manager" > "virtualbox.desktop"
	for _, file := range files {
		if file.Name() == strings.Split(strings.ToLower(className), " ")[0]+".desktop" {
			return filepath.Join(dir, file.Name())
		}
	}

	// "GitHub Desktop" > "github-desktop.desktop"
	for _, file := range files {
		fileName := file.Name()
		fileName = strings.ToLower(fileName)
		classNameLower := strings.ToLower(className)
		classNameLower = strings.ReplaceAll(classNameLower, " ", "-")

		if fileName == classNameLower+".desktop" {
			return filepath.Join(dir, file.Name())
		}
	}

	// Chrome/Chromium webapp: "chrome-messenger.com__-Default" > "Messenger.desktop" (by martonbtoth)
	if strings.HasPrefix(className, "chrome-") || strings.HasPrefix(className, "chromium-") {
		// Extract domain from class name (e.g., "chrome-messenger.com__-Default" -> "messenger.com")
		parts := strings.SplitN(className, "-", 2)
		if len(parts) == 2 {
			domain := strings.Split(parts[1], "__")[0] // Remove "__-Default" suffix
			domain = strings.TrimSuffix(domain, "-")
			domainParts := strings.Split(domain, ".")
			if len(domainParts) > 0 {
				// Try matching by domain name (e.g., "messenger" from "messenger.com")
				baseName := domainParts[0]
				for _, file := range files {
					fileName := file.Name()
					fileNameLower := strings.ToLower(fileName)
					if strings.Contains(fileNameLower, strings.ToLower(baseName)) && strings.HasSuffix(fileNameLower, ".desktop") {
						return filepath.Join(dir, fileName)
					}
				}
			}
		}
	}

	return ""
}
