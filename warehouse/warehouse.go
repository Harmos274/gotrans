package warehouse

import (
	"fmt"
	"time"
)

type Warehouse struct {
	Length, Height int
	Packages       EntityMap[Package]
	PalletJacks    EntityMap[PalletJack]
	Trucks         EntityMap[Truck]
}

func (wh Warehouse) SomethingExistsAt(pos Position) bool {
	return wh.Packages.Exists(pos) || wh.PalletJacks.Exists(pos) || wh.Trucks.Exists(pos)
}

func (wh Warehouse) Clone() Warehouse {
	var cloned Warehouse

	cloned.Height = wh.Height
	cloned.Length = wh.Length
	cloned.Packages = copyMap(wh.Packages)
	cloned.PalletJacks = copyMap(wh.PalletJacks)
	cloned.Trucks = copyMap(wh.Trucks)

	return cloned
}

type EntityMap[T Package | PalletJack | Truck] map[Position]T

type Position struct {
	X, Y int
}

type Package struct {
	Weight Weight
	Name   string
}
type Weight int

type PalletJack struct {
	Name  string
	State string
	pack  *Package
}

type Truck struct {
	Name                  string
	State                 string
	MaxWeight             Weight
	CurrentWeight         Weight
	ElapseDischargingTime int
}

func (ettMap EntityMap[T]) Exists(pos Position) bool {
	_, exists := ettMap[pos]
	return exists
}

func CleanWarehouse(initialWarehouse Warehouse, ch chan Warehouse, cycles uint) {
	defer close(ch)
	currentWarehouse := initialWarehouse
	paths := refreshPaths(currentWarehouse, make([]Path, 0))

	fmt.Println("Initial paths", paths)
	for ; cycles != 0; cycles -= 1 {
		paths = applyPaths(initialWarehouse, paths)
		fmt.Println("Paths after application", paths)

		ch <- currentWarehouse.Clone()
		time.Sleep(100 * time.Millisecond)

		paths = refreshPaths(currentWarehouse, paths)
		fmt.Println("Paths after refresh", paths)
	}
}

func applyPaths(wh Warehouse, paths []Path) []Path {
	for index := 0; index < len(paths); {
		path := paths[index]

		fmt.Println("Tring to apply", path)

		if path.isValid() {
			forklift := wh.PalletJacks[path.current]

			if len(path.steps) == 0 {
				if wh.Trucks.Exists(path.destination) {
					truck := wh.Trucks[path.destination]

					if truck.CurrentWeight+forklift.pack.Weight <= truck.MaxWeight {
						truck.CurrentWeight += forklift.pack.Weight
						forklift.pack = nil

						wh.PalletJacks[path.current] = forklift
						wh.Trucks[path.destination] = truck
						paths[index] = paths[len(paths)-1]
						paths = paths[:len(paths)-1]
						fmt.Println(path.current, "dropped package to", path.destination)
					}
				} else if wh.Packages.Exists(path.destination) {
					// Take package from map
					pack := wh.Packages[path.destination]
					delete(wh.Packages, path.destination)

					// Give package to forklift
					forklift.pack = &pack
					wh.PalletJacks[path.current] = forklift

					paths[index] = paths[len(paths)-1]
					paths = paths[:len(paths)-1]
					fmt.Println(path.current, "took", path.destination)
				}
			} else {
				delete(wh.PalletJacks, path.current)
				wh.PalletJacks[path.steps[0]] = forklift

				fmt.Println("Moved", path.current, "to", path.steps[0])
				paths[index].current = path.steps[0]
				paths[index].steps = path.steps[1:]
				index += 1
			}
		} else {
			fmt.Println("Invalid path was kept")
		}
	}

	return paths
}

func copyMap[T Package | PalletJack | Truck](toClone map[Position]T) map[Position]T {
	ret := make(map[Position]T)
	for key, value := range toClone {
		ret[key] = value
	}
	return ret
}
