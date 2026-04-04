// Package godesktop provides functionality to parse and interact with Linux desktop entries (.desktop files).
// It allows discovering, managing, and launching applications from standard and custom directories.
package godesktop

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/lotos-linux/desktop/utils"
)

// Manager manages desktop application discovery and retrieval.
// It handles application directories, locale preferences, and caching of application lookups.
type Manager struct {
	// locale stores the language/locale code for localized application names and comments
	locale string

	// defAppDirs contains the default system application directories
	defAppDirs []string
	// customAppDirs contains user-added custom application directories
	customAppDirs []string

	// classNameHash caches the mapping between application class names and their .desktop file paths
	classNameHash map[string]string
}

// New creates a new Manager instance with optional locale specification.
//
// Parameters:
//   - locale: Optional locale code (e.g., "en_US", "ru_RU"). If provided, localized application
//     names and comments will be returned in this language. If omitted, system default is used.
//
// Returns:
//   - *Manager: A configured Manager instance ready for use.
func New(locale ...string) *Manager {
	var lang string
	if len(locale) == 1 {
		lang = locale[0]
	}

	return &Manager{
		locale: lang,

		defAppDirs:    utils.DefaultAppDirs(),
		customAppDirs: make([]string, 0),

		classNameHash: nil,
	}
}

// AddCustomAppDirs adds custom directories to search for .desktop files.
// These directories will be searched before the default system directories.
//
// Parameters:
//   - dir: One or more directory paths to add to the search path.
func (m *Manager) AddCustomAppDirs(dir ...string) {
	m.customAppDirs = append(m.customAppDirs, dir...)
}

// GetAllApps returns all discovered applications from both default and custom directories.
// Duplicate applications (same .desktop file appearing in multiple directories) are filtered out.
//
// Returns:
//   - []*App: A slice of App instances representing all discovered applications.
func (m *Manager) GetAllApps() []*App {
	res := []*App{}

	dirs := utils.ProcessDirectories(append(m.customAppDirs, m.defAppDirs...))
	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, file := range files {
			path := filepath.Join(dir, file.Name())
			app, err := m.NewAppInFile(path)
			if err != nil {
				continue
			}
			res = append(res, app)
		}
	}

	return res
}

// FindByClass finds an application by its class name (the base name of the .desktop file).
// It first checks the cached hash map, then searches through directories if not found.
//
// Parameters:
//   - className: The class name of the application (e.g., "firefox", "org.gnome.Nautilus").
//
// Returns:
//   - *App: The found application instance.
//   - error: Error if the application is not found in any directory.
func (m *Manager) FindByClass(className string) (*App, error) {
	path, exist := m.getClassNameHash()[className]
	if exist {
		return m.NewAppInFile(path)
	}

	dirs := utils.ProcessDirectories(append(m.customAppDirs, m.defAppDirs...))
	for _, dir := range dirs {
		path := utils.SearchDesktopFile(className, dir)
		if path != "" {
			return m.NewAppInFile(path)
		}
	}

	return nil, errors.New("desktop file not found")
}

// getClassNameHash returns a cached hash map mapping class names to their .desktop file paths.
// The cache is created lazily on first call and reused for subsequent lookups.
func (m *Manager) getClassNameHash() map[string]string {
	if m.classNameHash == nil {
		dirs := utils.ProcessDirectories(append(m.customAppDirs, m.defAppDirs...))
		m.classNameHash = utils.CreateClassNameHash(dirs)
	}
	return m.classNameHash
}
