# MacOS Traymenu
This is a small configurable tray menu for MacOS (not sure about other platforms)

# Installation
1. Clone the Project
2. Run `make` to build the project
3. Run `make create-default-config` to copy the default config to your home directory
4. Open the config at `~/.traymenu/config.json` and edit to your liking
5. To test if everything works simpy execute the `./traymenu` and look at the output

# Autostart (MacOS)
1. Open the `Automator`-App and create a new Application
2. Add a new `Run Shell Script` action and give it the following command to execute: `cd /path/to/your/traymenu/repo && ./start-traymenu.sh`
3. Save it
4. Open the System-Settings, go to "Users & Groups", select your User and switch to "Login Items".
5. Click on the plus-sign and select the application you created in the previous steps.
6. To start it now, just double click on the entry. It should start whenever you log into your User-Account now.

# Config
The Configuration is pretty simple and essentially only as one Object-Type called Menu-Item that can be recursively be stacked through the `Items`-Property:
```
type MenuItem struct {
	Icon    string
	Title   string
	Tooltip string

	Items             *[]MenuItem // Submenu-Items (only Items or Action can be set at a time)
	Action            *string     // OnClick this will be executed (only Items or Action can be set at a time)
	CancellableAction bool        // Defines if the Action can be cancled. This will result in a Start/Stop Button
}
```
Important: Every Menu-Item can either have an Action or Subitems. Both is not possible!
