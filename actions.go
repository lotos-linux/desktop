package godesktop

import (
	"strings"

	"github.com/lotos-linux/desktop/utils"
)

// Action represents a desktop action from a .desktop file.
// Desktop actions provide alternative ways to launch an application
// (e.g., "New Window", "Open Profile Manager", "Compose Message").
type Action struct {
	// Name is the localized display name of the action
	Name string
	// Exec is the command line to execute this action
	Exec string
	// Icon is the icon name or path for this action
	Icon string
}

// getActions parses and extracts all desktop actions from a .desktop file's raw data.
//
// Parameters:
//   - raw: The complete parsed .desktop file data.
//   - locale: The locale code to use for localized action names.
//
// Returns:
//   - []Action: A slice of Action instances found in the .desktop file.
func getActions(raw map[string]map[string]string, locale string) []Action {
	var actionsRes []Action

	general, exist := raw["Desktop Entry"]
	if !exist {
		return actionsRes
	}

	actionsStr, exist := general["Actions"]
	if !exist {
		return actionsRes
	}

	actionsList := strings.Split(actionsStr, ";")

	for _, actionName := range actionsList {
		action := Action{}

		actionName = strings.TrimSpace(actionName)
		if actionName == "" {
			continue
		}

		key := "Desktop Action " + actionName
		actionGroup, exist := raw[key]
		if !exist {
			continue
		}

		name, ok := utils.GetAllLocales(actionGroup, "Name")
		if ok {
			action.Name = utils.GetLocalizedValue(name, locale)
		}

		exec, exist := actionGroup["Exec"]
		if exist {
			action.Exec = exec
		}

		icon, exist := actionGroup["Icon"]
		if exist {
			action.Icon = icon
		}

		actionsRes = append(actionsRes, action)
	}

	return actionsRes
}

// Run executes the action using its Exec command.
// It cleans and validates the command before launching.
//
// Returns:
//   - error: Error if the Exec command is invalid or cannot be started.
func (a *Action) Run() error {
	clean, err := utils.CleanExec(a.Exec)
	if err != nil {
		return err
	}

	return clean.Start()
}
