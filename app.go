package godesktop

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lotos-linux/desktop/utils"
)

// App represents a desktop application with its metadata and capabilities.
// It contains information parsed from a .desktop file including name, icon,
// execution command, and available actions.
type App struct {
	// Name is the localized application name
	Name string
	// Comment is a short description of the application (localized)
	Comment string
	// Icon is the icon name or path for the application
	Icon string
	// Exec is the command line to launch the application
	Exec string
	// SingleWindow indicates if the application should run as a single window instance
	SingleWindow bool

	// Categories specifies the application categories (e.g., "System;Utility;")
	// as defined in the Desktop Entry Specification
	Categories []string
	// MimeTypes lists the MIME types supported by the application
	MimeTypes []string
	// NoDisplay indicates if the application should be hidden from menus
	NoDisplay bool
	// Terminal specifies if the application needs to run in a terminal window
	Terminal bool
	// OnlyShowIn lists desktop environments where this application should be shown
	OnlyShowIn []string
	// NotShowIn lists desktop environments where this application should NOT be shown
	NotShowIn []string
	// Keywords provides search terms for finding the application (localized)
	Keywords []string
	// StartupWMClass is the window manager class used for startup notification
	StartupWMClass string
	// StartupNotify indicates if the desktop should show startup notification
	StartupNotify bool
	// Path specifies the working directory for the application
	Path string
	// DBusActivatable indicates if the application can be activated via D-Bus
	DBusActivatable bool

	// actions contains the application's desktop actions (e.g., "New Window", "Compose")
	actions []Action
	// raw stores the complete parsed .desktop file data
	raw map[string]map[string]string
}

// NewAppInFile creates an App instance by parsing a .desktop file from the given path.
//
// Parameters:
//   - file: The file system path to the .desktop file.
//
// Returns:
//   - *App: The parsed application instance.
//   - error: Error if the file cannot be read or parsed, or if required sections are missing.
func (m *Manager) NewAppInFile(file string) (*App, error) {
	data, err := utils.GetMap(file, "Desktop Entry")
	if err != nil {
		return nil, err
	}

	app, err := m.NewApp(data)
	if err != nil {
		return nil, fmt.Errorf("error create app in file: %s: error: %v", file, err)
	}

	return app, nil
}

// NewApp creates an App instance from parsed .desktop file data.
//
// Parameters:
//   - data: A map representing the parsed .desktop file structure, with group names as keys
//     and key-value maps as values.
//
// The method parses the following standard .desktop entry fields:
//   - Basic fields: Name, Comment, Icon, Exec, Path
//   - Categories: Split by semicolon into a slice
//   - MimeTypes: Split by semicolon into a slice
//   - Desktop environments: OnlyShowIn, NotShowIn
//   - Keywords: Localized and split by semicolon
//   - Boolean flags: NoDisplay, Terminal, SingleMainWindow, StartupNotify, DBusActivatable
//   - Window management: StartupWMClass
//   - Desktop actions: Parsed via getActions()
//
// Returns:
//   - *App: The constructed application instance.
//   - error: Error if the required "Desktop Entry" section is missing.
func (m *Manager) NewApp(data map[string]map[string]string) (*App, error) {
	res := &App{
		raw: data,
	}

	general, exist := res.raw["Desktop Entry"]
	if !exist {
		return nil, errors.New("\"Desktop\" Entry section not found")
	}

	if name, ok := utils.GetAllLocales(general, "Name"); ok {
		res.Name = utils.GetLocalizedValue(name, m.locale)
	}

	if comment, ok := utils.GetAllLocales(general, "Comment"); ok {
		res.Comment = utils.GetLocalizedValue(comment, m.locale)
	}

	res.Icon = general["Icon"]
	res.Exec = general["Exec"]
	res.Path = general["Path"]

	res.StartupWMClass = general["StartupWMClass"]

	if categories, ok := general["Categories"]; ok {
		res.Categories = strings.Split(categories, ";")
		res.Categories = utils.FilterEmpty(res.Categories)
	}

	if mimeTypes, ok := general["MimeType"]; ok {
		res.MimeTypes = strings.Split(mimeTypes, ";")
		res.MimeTypes = utils.FilterEmpty(res.MimeTypes)
	}

	if onlyShowIn, ok := general["OnlyShowIn"]; ok {
		res.OnlyShowIn = strings.Split(onlyShowIn, ";")
		res.OnlyShowIn = utils.FilterEmpty(res.OnlyShowIn)
	}

	if notShowIn, ok := general["NotShowIn"]; ok {
		res.NotShowIn = strings.Split(notShowIn, ";")
		res.NotShowIn = utils.FilterEmpty(res.NotShowIn)
	}

	keywords, ok := utils.GetAllLocales(general, "Keywords")
	if ok {
		keywordsStr := utils.GetLocalizedValue(keywords, m.locale)
		res.Keywords = strings.Split(keywordsStr, ";")
		res.Keywords = utils.FilterEmpty(res.Keywords)
	}

	res.NoDisplay = general["NoDisplay"] == "true"
	res.Terminal = general["Terminal"] == "true"
	res.SingleWindow = general["SingleMainWindow"] == "true"
	res.StartupNotify = general["StartupNotify"] == "true"
	res.DBusActivatable = general["DBusActivatable"] == "true"

	res.actions = getActions(res.raw, m.locale)

	return res, nil
}

// GetActions returns all desktop actions associated with the application.
// Desktop actions are alternative launch modes (e.g., "Open in new window").
//
// Returns:
//   - []Action: A slice of Action instances representing the application's actions.
func (a *App) GetActions() []Action {
	return a.actions
}

// GetRaw returns the complete raw parsed data from the .desktop file.
// This provides access to all sections and keys, even those not exposed by the App struct.
//
// Returns:
//   - map[string]map[string]string: The complete parsed .desktop file data.
func (a *App) GetRaw() map[string]map[string]string {
	return a.raw
}

// Run launches the application using its Exec command.
// It cleans and validates the command before execution.
//
// Returns:
//   - error: Error if the Exec command is invalid or cannot be started.
func (a *App) Run() error {
	clean, err := utils.CleanExec(a.Exec)
	if err != nil {
		return err
	}

	return clean.Start()
}
