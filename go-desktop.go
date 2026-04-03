package godesktop

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/sujer-ux/go-desktop/utils"
)

type Manager struct {
	locale string

	defAppDirs    []string
	customAppDirs []string

	classNameHash map[string]string
}

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

func (m *Manager) AddCustomAppDirs(dir ...string) {
	m.customAppDirs = append(m.customAppDirs, dir...)
}

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

func (m *Manager) getClassNameHash() map[string]string {
	if m.classNameHash == nil {
		dirs := utils.ProcessDirectories(append(m.customAppDirs, m.defAppDirs...))
		m.classNameHash = utils.CreateClassNameHash(dirs)
	}
	return m.classNameHash
}
