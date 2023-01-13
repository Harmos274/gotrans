# Gotrans

This is a [Golang](https://go.dev) project, therefore you'll need Go installed on your system.

## Goal

Giving a file that describes a warehouse with packages, forklifts and trucks inside it, the program
will have to optimise the distribution of packages to trucks using the forklifts.

To achieve that, **Gotrans** uses a pathfinding algorithm that orchestrate the movements of forklifts
around the map.

## Build

In order to build the project, you must have [gcc](https://gcc.gnu.org) and [OpenGL](https://www.opengl.org) libraries
installed on your device.

You can build the project with `go build`.

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

## Repository design

The sources are organised through 2 packages, the main package, `gotrans`, located at the root of the
directory, and its internal package `warehouse`.

The role of the `gotrans` package is to manage the communication with the user it contains the entry point
of the program in the `gotrans.go` file.

There's also the `parse_input_file.go` file, whose role is to handle the parsing of the input file given by
the users inputs, the `show_warehouse.go` that contains everything needed to print the warehouse on the
terminal and `graphical.go` that contains the functions needed to run the graphical UI.

In the `warehouse` package, the `warehouse.go` file contains the description of the Warehouse and the
functions to modify its data. The `event.go` file describes all the events occurring during the warehouse
cleaning execution cycles. Finally, the `al.go` file contains the pathfinding algorithm used in the cleaning
warehouse process.

## Pathfinding strategy

For each `forklift`, the algorithm finds the quickest path to go to every `package` in the `warehouse`,
if the `forklift` find a shorter path than another, its path is chosen, otherwise the `package` is
removed from the `forklift`'s targets and another `package` is selected.

Once the `forklift` arrives by the `package` it picks it up and search the quickest path to the `truck`.
When the `forklift` arrives by the `truck` it loads its `package` in the `truck`, if possible, otherwise it waits.

Once the `package` has been delivered the `forklift` go to another targets if there's one.