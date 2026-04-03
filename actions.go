package godesktop

import (
	"strings"

	"github.com/lotos-linux/go-desktop/go-desktop/utils"
)

type Action struct {
	Name string
	Exec string
	Icon string
}

func GetActions(raw map[string]map[string]string, locale string) []Action {
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

func (a *Action) Run() error {
	clean, err := utils.CleanExec(a.Exec)
	if err != nil {
		return err
	}

	return clean.Start()
}
