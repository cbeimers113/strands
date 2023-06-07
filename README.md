# strands
An ecosystem simulator in early development written in Go using [g3n engine](https://github.com/g3n/engine).

### Install Dependencies
To install g3n engine's dependencies for your platform, see the [Dependencies](https://github.com/g3n/engine#dependencies) section of the README. (Don't worry about the Installation section)

### How to build and run:
```
git clone https://github.com/cbeimers113/strands.git
cd strands

go build . && ./strands
or
go run main.go
```

### Current Features:

    * Controls and info section on HUD
    * Randomly generated hextile map
    * Water that flows down the map's topography
    * Plants which will grow on tiles over time
