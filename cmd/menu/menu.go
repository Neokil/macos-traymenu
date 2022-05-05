//nolint: forbidigo
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/getlantern/systray"
)

var (
	iconCache map[string][]byte = map[string][]byte{} // nolint: gochecknoglobals,revive
	homeDir   string                                  // nolint: gochecknoglobals
)

type Config MenuItem

type MenuItem struct {
	Icon    string
	Title   string
	Tooltip string

	Items             *[]MenuItem // Submenu-Items (only Items or Action can be set at a time)
	Action            *string     // OnClick this will be executed (only Items or Action can be set at a time)
	CancellableAction bool        // Defines if the Action can be cancled. This will result in a Start/Stop Button
}

func (mi MenuItem) GetIcon() []byte {
	return GetIcon(mi.Icon)
}

func (mi Config) GetIcon() []byte {
	return GetIcon(mi.Icon)
}

func GetIcon(path string) []byte {
	if strings.HasPrefix(path, "~/") {
		path = homeDir + path[1:]
	}

	icon, exists := iconCache[path]
	if exists {
		return icon
	}

	icon, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("could not load icon %s: %w", path, err))
	}
	iconCache[path] = icon

	return icon
}

type app struct {
	config Config
}

func loadConfig() Config {
	configLocations := []string{
		"",
		homeDir + "/.traymenu/",
	}

	for _, cl := range configLocations {
		c, err := loadConfigFromPath(cl)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			panic(err)
		}

		return c
	}

	panic("No Config found")
}

func loadConfigFromPath(path string) (Config, error) {
	b, err := os.ReadFile(path + "config.json")
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config: %w", err)
	}
	config := Config{}
	err = json.Unmarshal(b, &config)
	if err != nil {
		panic(err)
	}

	return config, nil
}

func main() {
	var err error
	homeDir, err = os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	app := &app{
		config: loadConfig(),
	}

	systray.Run(app.ready, app.exit)
}

func setupActionDefault(mi *systray.MenuItem, config MenuItem) {
	go func(trigger chan struct{}, action string) {
		for range trigger {
			fmt.Println("Executing: " + action)
			cmd := exec.Command("bash", "-c", action)
			cmd.Stdout = os.Stdout
			err := cmd.Run()
			if err != nil {
				fmt.Println("Error: " + err.Error())
			}
		}
	}(mi.ClickedCh, *config.Action)
}

func setupActionCancelable(mi *systray.MenuItem, config MenuItem) {
	mi.SetTitle("Start: " + config.Title)

	go func(trigger chan struct{}, action string) {
		started := false

		var cmd *exec.Cmd

		for range trigger {
			if started {
				fmt.Printf("Killing process group %d\n", -cmd.Process.Pid)
				err := syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
				if err != nil {
					fmt.Printf("Failed to kill process: %s\n", err.Error())
				}
			} else {
				started = true
				mi.SetTitle("Stop: " + config.Title)

				fmt.Println("Executing: " + action)
				cmd = exec.Command("bash", "-c", action)
				cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				go func() {
					err := cmd.Run()
					if err != nil {
						fmt.Println("Error: " + err.Error())
					}
					started = false
					mi.SetTitle("Start: " + config.Title)
					fmt.Println("Execution finished")
				}()
			}
		}
	}(mi.ClickedCh, *config.Action)
}

func setupAction(mi *systray.MenuItem, config MenuItem) {
	mi.SetIcon(config.GetIcon())

	if (config.Action != nil) == (config.Items != nil) {
		panic(fmt.Sprintf("You need to define exactly one of either the Action or the Items for every Menu-Item (%s)", config.Title))
	}

	switch {
	case config.Action != nil && config.CancellableAction:
		setupActionCancelable(mi, config)
	case config.Action != nil && !config.CancellableAction:
		setupActionDefault(mi, config)
	case config.Items != nil:
		for _, si := range *config.Items {
			InitMenuItem(mi, si)
		}
	}
}

func (app app) ready() {
	fmt.Println("Ready...")
	systray.SetTemplateIcon(app.config.GetIcon(), app.config.GetIcon())
	systray.SetTitle(app.config.Title)
	systray.SetTooltip(app.config.Tooltip)

	for _, i := range *app.config.Items {
		mi := systray.AddMenuItem(i.Title, i.Tooltip)
		setupAction(mi, i)
	}

	systray.AddSeparator()
	quitButton := systray.AddMenuItem("Quit", "Quit application")
	go func() {
		<-quitButton.ClickedCh
		systray.Quit()
	}()
}

func InitMenuItem(root *systray.MenuItem, config MenuItem) {
	mi := root.AddSubMenuItem(config.Title, config.Tooltip)
	setupAction(mi, config)
}

func (app app) exit() {
	fmt.Println("Exiting...")
}
