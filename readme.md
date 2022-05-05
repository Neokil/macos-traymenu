# MacOS Traymenu
This is a small configurable tray menu for MacOS (not sure about other platforms)

# Installation
1. Clone the Project
2. Run `make` to build the project
3. Run `make create-default-config` to copy the default config to your home directory
4. Open the config at `~/.traymenu/config.json` and edit to your liking
5. To test if everything works simpy execute the `./traymenu` and look at the output
6. (Optional) Add it to your Auto-Start-Items by opening the System-Settings, going to "Users & Groups", selecting your User ans switching to "Login Items". There click on the plus-sign and select the "traymenu" file from the project-root and Tick the Box that says "hide".

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
