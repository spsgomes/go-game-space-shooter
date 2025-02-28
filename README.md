# Context

This small project was created to explore game development using Go.\
This is not meant to be a *production-ready* project, instead it aims to enable me to explore concepts in Go, mathematics, 2D graphics, SFX, and game design.

## Features
- Built using Go, with [Ebiten](https://ebitengine.org/) as a 2D graphics game engine
- Cross-platform compatibility
- Pseudo-randomness using seeds\
  For clarification on using seeds for random number generators: [Random Seeds and Reproducibility](https://medium.com/towards-data-science/random-seeds-and-reproducibility-933da79446e3) by Daniel Godoy.

## To Do
- [X] Add UI (e.g. Player Health; Death Screen; etc.)
- [ ] Save the High Score to disk for score permanence
- [ ] Add a main menu, and "Retry" button
- [ ] Add SFX for player and enemy damage
- [ ] Add Power-ups (Health Pickup; Damage Boosters; Critical Chance; Critical Modifier; etc.)
- [ ] Add statistics to UI (Current Damage; Current Critical Chance & Modifier; etc.)
- [ ] HP bars on top of enemies

## Installation

### Prerequisites
- Install [Golang](https://go.dev/dl/)

### Clone the repository
```sh
git clone https://github.com/spsgomes/go-game-space-shooter.git
cd go-game-space-shooter
```

### Run the game
```sh
go run cmd/go-game-space-shooter/main.go
```

### Build the game
To build for your current OS:
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

## How to Play

### Controls
- `W Key`: Up
- `S Key`: Down
- `A Key`: Left
- `D Key`: Right
- `Space` / `Left Click`: Shoot

### Objective
Destroy the enemy ships, and earn a High Score!

## Configurations
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

## Dependencies
- [Ebiten](https://ebitengine.org/) for 2D graphics game engine

## Assets
- Art: [Space Shooter Redux by Kenney](https://kenney.nl/assets/space-shooter-redux)
- Background music: [EXAGGERATE | TEMPTATION by Rhapsody](https://freemusicarchive.org/music/rhapsody/single/exaggerate-temptation/)

## License
This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.
