package warehouse

import (
	"log"
)

type Warehouse struct {
	Length, Height int
	Packages       EntityMap[Package]
	ForkLifts      EntityMap[ForkLift]
	Trucks         EntityMap[Truck]
}

type CycleState struct {
	Warehouse Warehouse
	Events    []Event
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

func CleanWarehouse(wh Warehouse, ch chan CycleState, cycles uint) {
	defer close(ch)
	paths := refreshPaths(wh, make([]Path, 0))
	var events []Event

	for ; cycles != 0 && !isOver(wh); cycles-- {
		paths, events = applyPaths(wh, paths)

		ch <- CycleState{Warehouse: wh.Clone(), Events: events}

		paths = refreshPaths(wh, paths)
	}
}

func applyPaths(wh Warehouse, paths []Path) ([]Path, []Event) {
	events := []Event{}
	waitingForklifts := mapToPositionSet(wh.ForkLifts)

	for index := 0; index < len(paths); {
		path := paths[index]

		if path.isValid() {
			forklift := wh.ForkLifts[path.current]
			delete(waitingForklifts, path.current)

			if len(path.steps) == 0 {
				if wh.Trucks.Exists(path.destination) {
					paths, index, events = dropPackage(path, forklift, index, wh.ForkLifts, wh.Trucks, paths, events)
				} else if wh.Packages.Exists(path.destination) {
					paths, events = takePackage(path, forklift, index, wh.ForkLifts, wh.Packages, paths, events)
				}
			} else {
				paths, events = moveForkLift(path, forklift, index, wh.ForkLifts, paths, events)
				index++
			}
		} else {
			log.Fatal("Invalid path was kept")
		}
	}

	for pos := range waitingForklifts {
		waiting := wh.ForkLifts[pos]

		events = append(events, ForkliftWait{forkliftName: waiting.Name, position: pos})
	}

	for pos, truck := range wh.Trucks {
		events = append(events, TruckWait{truckName: truck.Name, truckMaxWeight: truck.MaxWeight, truckLoadedWeight: truck.CurrentWeight, position: pos})
	}

	return paths, events
}

func moveForkLift(path Path, forklift ForkLift, index int, forkLifts EntityMap[ForkLift],
	paths []Path, events []Event,
) ([]Path, []Event) {
	if !forkLifts.Exists(path.steps[0]) {
		events = append(events, ForkliftMove{forkliftName: forklift.Name, eventPosition: path.current, target: path.steps[0]})

		delete(forkLifts, path.current)
		forkLifts[path.steps[0]] = forklift

		paths[index].current = path.steps[0]
		paths[index].steps = path.steps[1:]
	} else {
		events = append(events, ForkliftWait{forkliftName: forklift.Name, position: path.current})

		paths[index] = paths[len(paths)-1]
		paths = paths[:len(paths)-1]
	}

	return paths, events
}

func takePackage(path Path, forklift ForkLift, index int, forkLifts EntityMap[ForkLift],
	packages EntityMap[Package], paths []Path, events []Event,
) ([]Path, []Event) {
	events = append(events, PickupPackage{position: path.current, emitterName: forklift.Name, packName: packages[path.destination].Name})

	// Take package from map
	pack := packages[path.destination]
	delete(packages, path.destination)

	// Give package to forklift
	forklift.pack = &pack
	forkLifts[path.current] = forklift

	paths[index] = paths[len(paths)-1]
	return paths[:len(paths)-1], events
}

func dropPackage(path Path, forklift ForkLift, index int, forkLifts EntityMap[ForkLift],
	trucks EntityMap[Truck], paths []Path, events []Event,
) ([]Path, int, []Event) {
	truck := trucks[path.destination]

	if truck.CurrentWeight+forklift.pack.Weight <= truck.MaxWeight {
		events = append(events, DeliverPackage{position: path.current, emitterName: forklift.Name, packName: forklift.pack.Name})

		truck.CurrentWeight += forklift.pack.Weight
		forklift.pack = nil

		forkLifts[path.current] = forklift
		trucks[path.destination] = truck
		paths[index] = paths[len(paths)-1]
		paths = paths[:len(paths)-1]
	} else {
		events = append(events, ForkliftWait{forkliftName: forklift.Name, position: path.current})

		index++
	}

	return paths, index, events
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
