# godesktop

A Go library for parsing and interacting with Linux desktop entries (.desktop files). Provides application discovery, management, and launching capabilities from standard and custom directories.

## Installation

```bash
go get github.com/sujer-ux/go-desktop
```

## API Reference

### Manager

The main entry point for application discovery and management.

#### `New(locale ...string) *Manager`

Creates a new Manager instance.

**Parameters:**
- `locale` (optional): Locale code (e.g., "en_US", "ru_RU") for localized application names and comments

**Returns:** Configured Manager instance

#### `AddCustomAppDirs(dir ...string)`

Adds custom directories to search for .desktop files. These directories are searched before default system directories.

#### `GetAllApps() []*App`

Returns all discovered applications from both default and custom directories. Duplicate applications are filtered out.

#### `FindByClass(className string) (*App, error)`

Finds an application by its class name (base name of the .desktop file).

**Parameters:**
- `className`: Application class name (e.g., "firefox", "org.gnome.Nautilus")

**Returns:** Application instance or error if not found

#### `NewAppInFile(file string) (*App, error)`

Creates an App instance by parsing a .desktop file from the given path.

### App

Represents a desktop application with its metadata and capabilities.

#### Fields

| Field | Type | Description |
|-------|------|-------------|
| `Name` | `string` | Localized application name |
| `Comment` | `string` | Short description (localized) |
| `Icon` | `string` | Icon name or path |
| `Exec` | `string` | Command line to launch application |
| `SingleWindow` | `bool` | Single window instance flag |
| `Categories` | `[]string` | Application categories |
| `MimeTypes` | `[]string` | Supported MIME types |
| `NoDisplay` | `bool` | Hide from menus flag |
| `Terminal` | `bool` | Run in terminal window flag |
| `OnlyShowIn` | `[]string` | Desktop environments to show |
| `NotShowIn` | `[]string` | Desktop environments to hide |
| `Keywords` | `[]string` | Search keywords (localized) |
| `StartupWMClass` | `string` | Window manager class |
| `StartupNotify` | `bool` | Startup notification flag |
| `Path` | `string` | Working directory |
| `DBusActivatable` | `bool` | D-Bus activation flag |

#### Methods

##### `GetActions() []Action`

Returns all desktop actions associated with the application.

##### `GetRaw() map[string]map[string]string`

Returns complete raw parsed data from the .desktop file.

##### `Run() error`

Launches the application using its Exec command.

### Action

Represents a desktop action providing alternative launch modes.

#### Fields

| Field | Type | Description |
|-------|------|-------------|
| `Name` | `string` | Localized display name |
| `Exec` | `string` | Command to execute |
| `Icon` | `string` | Icon name or path |

#### Methods

##### `Run() error`

Executes the action command.

## Example Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/sujer-ux/go-desktop"
)

func main() {
    // Initialize manager with Russian locale
    manager := godesktop.New("ru_RU")
    
    // Add custom application directory
    manager.AddCustomAppDirs("/opt/custom-apps")
    
    // Get all applications
    apps := manager.GetAllApps()
    fmt.Printf("Found %d applications\n", len(apps))
    
    // Find application by class name
    firefox, err := manager.FindByClass("firefox")
    if err != nil {
        log.Printf("Firefox not found: %v", err)
    } else {
        fmt.Printf("Application: %s\n", firefox.Name)
        fmt.Printf("Description: %s\n", firefox.Comment)
        fmt.Printf("Command: %s\n", firefox.Exec)
        fmt.Printf("Categories: %v\n", firefox.Categories)
        fmt.Printf("MIME Types: %v\n", firefox.MimeTypes)
        fmt.Printf("Terminal required: %t\n", firefox.Terminal)
        
        // Launch application
        if err := firefox.Run(); err != nil {
            log.Printf("Failed to launch: %v", err)
        }
    }
    
    // Parse specific .desktop file
    app, err := manager.NewAppInFile("/usr/share/applications/org.gnome.Nautilus.desktop")
    if err != nil {
        log.Printf("Failed to parse: %v", err)
    } else {
        fmt.Printf("\nNautilus Details:\n")
        fmt.Printf("Name: %s\n", app.Name)
        fmt.Printf("Icon: %s\n", app.Icon)
        fmt.Printf("Startup WM Class: %s\n", app.StartupWMClass)
        
        // Work with desktop actions
        actions := app.GetActions()
        for _, action := range actions {
            fmt.Printf("Action: %s -> %s\n", action.Name, action.Exec)
            
            // Execute specific action
            // if action.Name == "New Window" {
            //     action.Run()
            // }
        }
        
        // Access raw data
        raw := app.GetRaw()
        if desktopEntry, exists := raw["Desktop Entry"]; exists {
            fmt.Printf("Raw Version: %s\n", desktopEntry["Version"])
        }
    }
    
    // Iterate through all applications with filtering
    for _, app := range apps {
        if !app.NoDisplay && len(app.Categories) > 0 {
            fmt.Printf("\n%s (%s)\n", app.Name, app.Exec)
            fmt.Printf("  Comment: %s\n", app.Comment)
            fmt.Printf("  Categories: %v\n", app.Categories)
            
            if len(app.Keywords) > 0 {
                fmt.Printf("  Keywords: %v\n", app.Keywords)
            }
            
            // Launch application based on condition
            // if app.Name == "Terminal" {
            //     app.Run()
            // }
        }
    }
    
    // Working with actions
    if firefox != nil {
        actions := firefox.GetActions()
        for _, action := range actions {
            fmt.Printf("\nFirefox Action: %s\n", action.Name)
            fmt.Printf("  Exec: %s\n", action.Exec)
            if action.Icon != "" {
                fmt.Printf("  Icon: %s\n", action.Icon)
            }
            
            // Execute action
            // if action.Name == "New Private Window" {
            //     if err := action.Run(); err != nil {
            //         log.Printf("Action failed: %v", err)
            //     }
            // }
        }
    }
    
    // Handle launch errors
    unknownApp, err := manager.FindByClass("nonexistent.app")
    if err != nil {
        fmt.Printf("Expected error: %v\n", err)
    } else {
        if err := unknownApp.Run(); err != nil {
            fmt.Printf("Launch error: %v\n", err)
        }
    }
}
```

## Default Directories

The library searches for .desktop files in the following default locations:
- `/usr/share/applications/`
- `/usr/local/share/applications/`
- `~/.local/share/applications/`

Custom directories added via `AddCustomAppDirs()` take precedence over default directories.

## Error Handling

The library returns errors for:
- Missing or unreadable .desktop files
- Invalid .desktop file format
- Missing required "Desktop Entry" sections
- Application not found by class name
- Invalid or malformed Exec commands
- Failed command execution