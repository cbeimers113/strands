# Strands
An ecosystem simulator in early development written in Go using [g3n engine](https://github.com/g3n/engine).

## Install Dependencies
To install g3n engine's dependencies for your platform, see the [Dependencies](https://github.com/g3n/engine#dependencies) section of the README. (Don't worry about the Installation section)

## How to build and run:
```
git clone https://github.com/cbeimers113/strands.git
cd strands

go build . && ./strands
or
go run main.go
```

## History

Check out the [Releases](https://github.com/cbeimers113/strands/releases) page for the versions listed below!

<details>
    <summary>Changelog</summary>

### v0.2.1 - Bugfixes & GUI Enhancements
- Increase GUI scale
- Improve player icon in top bar
- Maximize by default
- Fullscreen/windowed option in config
- Fix memory leak bug between opening saves
- Reset position on new sim
- Fix autosave bug/persistent plants issue
- Toggle info panel for chemical quantities in top right
- Notifications appear in bottom left
- Topbar instead of info panel
- Space & control to go vertically up and down
- Limit how far from map we can go
- Tab to play/pause

### v0.2.0 - Application & Engine Features
- Spin the camera slowly when in menu
- GUI style improvements
- Notifications
  - Notifications fade out
  - Notify when save is saved or loaded
  - Notify when new sim is created
  - Notify when save is deleted
  - Notify when config saved
  - Notify when a seed is planted
- Add logo to menu
- Button to create a fresh sim
- Save player camera position and rotation in saves
- Icon
- Embed textures for static builds
- Config file is in %APPDATA%/Roaming
- Version number tracked in .version
- Optionally save the state of the simulation to the disk between games
- Keyboard controller for typing input
- Save and load game files
- Delete a save
- Popup system
  - Prevent overwriting save (exclude autosave)
  - Notify when save can't be loaded
  - Ask if user is sure before opening/deleting a save or exiting game
- Use a nerdfont

### v0.1.0 - General Enhancements:
- Settings menu to edit config
- Configurable day length and tick speed
- Configurable mouse sensitivity and player movement
- Controls section of info screen is togglable
- Right click tile to open context menu, plant a seed from context menu
- In-game clock & day-night cycle
- Refactored Views in the gui package
- Refactor entities as Entity interface

### v0.0.8
- Refactored project structure
- New JSON-based config system
- Smoother player controls
- Golang version update to v1.21.5

### v0.0.7
- Fixes issues with water simulation
- Adds a simulation pausing feature
- Updates README

### v0.0.6
- Plants spawn at random position on tiles with variations in colour
- Rewritten entity framework
- Better "looking at" highlighting (use wireframe instead of messing with texture color)
- Tiles have fertility value between 0 and 1 that determines how likely it is to support a plant
- Fixes a bug with the hexagon mesh that impacted relative entity positioning
- Adds water to the world and basic water physics (water spreading)
- Several small minor improvements

### v0.0.5
- Adds a WAILA (what am I looking at?) section to the HUD
- Adds a framework for dynamic elements in the atmosphere (gasses and liquids, may extend to powders in future)
- Improves the GUI framework
- Removes water as a tile type
- Improved cursor

### v0.0.4
- Refines player mouse movement
- Adds a menus system
- Adds a main/pause menu, togglable with esc
- Adds an infotext panel in the top left corner (may evolve into HUD)

### v0.0.3
- This patch replaces the default orbital controller with a more intuitive custom player controller framework.

### v0.0.2
- Changes tiles to hexagons, increases the map size and quality, and makes the map 3D

### v0.0.1
- The first working version that generates a simple tilemap and allows the user to move around the world and spawn plants on tiles.

</details>