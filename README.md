# Context

This small project was created to explore game development using Go.\
This is not meant to be a *production-ready* project, instead it aims to enable me to explore concepts in Go, mathematics, 2D graphics, SFX, and game design.

# Features
- Built using Go, with [Ebiten](https://ebitengine.org/) as a 2D graphics game engine
- Cross-platform compatibility
- Pseudo-randomness using seeds\
  For clarification on using seeds for random number generators: [Random Seeds and Reproducibility](https://medium.com/towards-data-science/random-seeds-and-reproducibility-933da79446e3) by Daniel Godoy.

## To Do
- [X] Add UI (e.g. Player Health; Death Screen; etc.)
- [X] Save the High Score to disk for score permanence
- [X] Add a main menu, and "Retry" button
- [X] Add SFX for player and enemy damage
- [X] HP bars on top of enemies

### If time allows
- [X] Add Power-ups
  - [X] Health Pickup
  - [X] Damage Boosters
  - [X] Critical Chance
  - [X] Critical Modifier
- [X] Add statistics to UI (Current Damage; Current Critical Chance & Modifier; etc.)
- [ ] Add Damage Numbers
- [ ] Add stronger enemies based on wave progression
- [ ] Add boss enemy based on wave progression
- [ ] Add "Play" button to main menu
- [ ] Add "Quit" button to main menu
- [ ] Add "Settings" menu, with button in main menu

# Installation

### Prerequisites
- Install [Golang](https://go.dev/dl/). This project was developed in version `1.24.0`, on `windows/amd64`.

### Clone the repository
```sh
git clone https://github.com/spsgomes/go-game-space-shooter.git
cd go-game-space-shooter
```

### Build & Run
If you prefer not using `make`, here are the commands necessary to run and build the project:

#### Run
```sh
go run cmd/go-game-space-shooter/main.go
```

#### Build
```sh
go build ./cmd/go-game-space-shooter/
./go-game-space-shooter
```
**Please note:** when building in different operating systems, the extension will be different (ex: Windows generates a .exe file)

If you prefer, you can specify the name of the executable by using `-o` flag, such as:
```sh
go build -o game.exe ./cmd/go-game-space-shooter/
./game.exe
```

# How to Play

### Controls
- `W Key`: Up
- `S Key`: Down
- `A Key`: Left
- `D Key`: Right
- `Space`/`Left Click`: Shoot
- `Esc`: Pause/Continue

### Objective
Destroy the enemy ships, and earn a High Score!


# Configuration
You can specify certain configurations (such as window size, enable fullscreen, etc.) by changing the file in `./configs.env`.

### Examples

To enable Fullscreen:
```go
FULLSCREEN_ENABLED: 1
```

To change the enemy wave spawn time:
```go
ENEMY_SPAWN_TIME: 10
```

## Notes
The game saves to the user configuration folder.\
This can typically be found in:
- Windows: `%AppData%`
- MacOS: `~/Library/application Support`
- Unix Systems: `$XDG_CONFIG_HOME` as specified by https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html

## Dependencies
- [Ebiten](https://ebitengine.org/) for 2D graphics game engine.\
  If you're using macOS or Linux, please visit [Ebiten Install page](https://ebitengine.org/en/documents/install.html) as the package requires some dependencies.

## Assets
- Art: [Space Shooter Redux by Kenney](https://kenney.nl/assets/space-shooter-redux)
- Background music: [EXAGGERATE | TEMPTATION by Rhapsody](https://freemusicarchive.org/music/rhapsody/single/exaggerate-temptation/)

## License
This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.
