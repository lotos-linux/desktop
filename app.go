package godesktop

import (
	"errors"
	"fmt"

	"github.com/sujer-ux/go-desktop/utils"
)

type App struct {
	Name         string
	Comment      string
	Icon         string
	Exec         string
	SingleWindow bool

	actions []Action
	raw     map[string]map[string]string
}

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

	res.actions = GetActions(res.raw, m.locale)

	return res, nil
}

func (a *App) GetActions() []Action {
	return a.actions
}

func (a *App) GetRaw() map[string]map[string]string {
	return a.raw
}

func (a *App) Run() error {
	clean, err := utils.CleanExec(a.Exec)
	if err != nil {
		return err
	}

	return clean.Start()
}
