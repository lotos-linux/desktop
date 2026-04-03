package godesktop

import (
	"errors"
	"fmt"

	"github.com/sujer-ux/go-desktop/utils"
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

	name, ok := utils.GetAllLocales(general, "Name")
	if ok {
		res.Name = utils.GetLocalizedValue(name, m.locale)
	}

	comment, ok := utils.GetAllLocales(general, "Comment")
	if ok {
		res.Comment = utils.GetLocalizedValue(comment, m.locale)
	}

	icon, exist := general["Icon"]
	if exist {
		res.Icon = icon
	}

	exec, exist := general["Exec"]
	if exist {
		res.Exec = exec
	}

	singleWindowStr, exist := general["SingleMainWindow"]
	res.SingleWindow = exist && singleWindowStr == "true"

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
