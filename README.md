# Gotrans

This is a [Golang](https://go.dev) project, therefore you'll need Go installed on your system.

## Goal

Giving a file that describes a warehouse with packages, forklifts and trucks inside it, the program
will have to optimise the distribution of packages to trucks using the forklifts.

To achieve that, **Gotrans** uses a pathfinding algorithm that orchestrate the movements of forklifts
around the map.

## Build

You can build the project with `go build`.

## Test

Run tests with `go test`.

## Run

You can run the project after build with the `gotrans` executable.

```
$> gotrans <file>
```

### Gotrans setup file

The file passed to **gotrans** executable describes the warehouse and its entities.

It is formatted as follows:

**First line**: Warehouse length, height and the number of execution cycles.
**X next lines**: Package name, X and Y position and color (yellow = 100Kg, green = 200Kg, blue = 500Kg).
**Y next lines**: Forklift name and X and Y position.
**Z next lines**: Truck name, X and Y position, max weight and cycle and cooldown after loading.

Example:
```
5 5 1000 -- Warehouse length, height and the number of execution cycles. 
colis_a_livrer 2 1 green -- Package name, X and Y position and color.
paquet 2 2 BLUE
deadpool 0 3 yellow
colère_DU_dragon 4 1 green
transpalette_1 0 0 -- Forklift name and X and Y position.
camion_b 3 4 4000 5 -- Truck name, X and Y position, max weight and cycle and cooldown after loading.
```

## Design

## Pathfinding strategy