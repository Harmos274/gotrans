package warehouse

import (
	"fmt"
)

type Warehouse struct {
	Length, Height int
	Packages       EntityMap[Package]
	ForkLifts      EntityMap[ForkLift]
	Trucks         EntityMap[Truck]
}

func (wh Warehouse) SomethingExistsAt(pos Position) bool {
	return wh.Packages.Exists(pos) || wh.ForkLifts.Exists(pos) || wh.Trucks.Exists(pos)
}

func (wh Warehouse) Clone() Warehouse {
	var cloned Warehouse

	cloned.Height = wh.Height
	cloned.Length = wh.Length
	cloned.Packages = copyMap(wh.Packages)
	cloned.ForkLifts = copyMap(wh.ForkLifts)
	cloned.Trucks = copyMap(wh.Trucks)

	return cloned
}

type EntityMap[T Package | ForkLift | Truck] map[Position]T

type Position struct {
	X, Y int
}

type Package struct {
	Weight Weight
	Name   string
}
type Weight int

type ForkLift struct {
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

func CleanWarehouse(wh Warehouse, ch chan Warehouse, cycles uint) {
	defer close(ch)
	paths := refreshPaths(wh, make([]Path, 0))

	for ; cycles != 0 && !isOver(wh); cycles -= 1 {
		paths = applyPaths(wh, paths)

		ch <- wh.Clone()

		paths = refreshPaths(wh, paths)
	}
}

func applyPaths(wh Warehouse, paths []Path) []Path {
	for index := 0; index < len(paths); {
		path := paths[index]

		if path.isValid() {
			forklift := wh.ForkLifts[path.current]

			if len(path.steps) == 0 {
				if wh.Trucks.Exists(path.destination) {
					paths, index = dropPackage(path, forklift, index, wh.ForkLifts, wh.Trucks, paths)
				} else if wh.Packages.Exists(path.destination) {
					paths = takePackage(path, forklift, index, wh.ForkLifts, wh.Packages, paths)
				}
			} else {
				paths = moveForkLift(path, forklift, index, wh.ForkLifts, paths)
				index += 1
			}
		} else {
			fmt.Println("Invalid path was kept")
		}
	}

	return paths
}

func moveForkLift(path Path, forklift ForkLift, index int, forkLifts EntityMap[ForkLift],
	paths []Path) []Path {
	delete(forkLifts, path.current)
	forkLifts[path.steps[0]] = forklift

	paths[index].current = path.steps[0]
	paths[index].steps = path.steps[1:]

	return paths
}

func takePackage(path Path, forklift ForkLift, index int, forkLifts EntityMap[ForkLift],
	packages EntityMap[Package], paths []Path) []Path {
	// Take package from map
	pack := packages[path.destination]
	delete(packages, path.destination)

	// Give package to forklift
	forklift.pack = &pack
	forkLifts[path.current] = forklift

	paths[index] = paths[len(paths)-1]
	return paths[:len(paths)-1]
}

func dropPackage(path Path, forklift ForkLift, index int, forkLifts EntityMap[ForkLift],
	trucks EntityMap[Truck], paths []Path) ([]Path, int) {
	truck := trucks[path.destination]

	if truck.CurrentWeight+forklift.pack.Weight <= truck.MaxWeight {
		truck.CurrentWeight += forklift.pack.Weight
		forklift.pack = nil

		forkLifts[path.current] = forklift
		trucks[path.destination] = truck
		paths[index] = paths[len(paths)-1]
		paths = paths[:len(paths)-1]
	} else {
		index += 1
	}

	return paths, index
}

func copyMap[T Package | ForkLift | Truck](toClone map[Position]T) map[Position]T {
	ret := make(map[Position]T)
	for key, value := range toClone {
		ret[key] = value
	}
	return ret
}

func isOver(wh Warehouse) bool {
	if len(wh.ForkLifts) == 0 {
		return true
	}

	if len(wh.Packages) > 0 {
		return false
	}

	for _, forklift := range wh.ForkLifts {
		if forklift.pack != nil {
			return false
		}
	}

	return true
}
