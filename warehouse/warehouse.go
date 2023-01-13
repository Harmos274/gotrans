package warehouse

import (
	"log"
)

// Warehouse description of the Warehouse
// Length length of the Warehouse
// Height height of the Warehouse
// ForkLifts map of every ForkLift associated to their Position in the Warehouse
// Packages map of every Package associated to their Position in the Warehouse
// Trucks map of every Truck associated to their Position in the Warehouse
type Warehouse struct {
	Length, Height int
	Packages       EntityMap[Package]
	ForkLifts      EntityMap[ForkLift]
	Trucks         EntityMap[Truck]
}

// CycleState association of the Warehouse and its associated events at a specific cycle
// Warehouse the warehouse
// Events the events
type CycleState struct {
	Warehouse Warehouse
	Events    []Event
}

// SomethingExistsAt checks if something exists at a position in Warehouse
func (wh Warehouse) SomethingExistsAt(pos Position) bool {
	return wh.Packages.Exists(pos) || wh.ForkLifts.Exists(pos) || wh.Trucks.Exists(pos)
}

// Clone clone a Warehouse
func (wh Warehouse) Clone() Warehouse {
	var cloned Warehouse

	cloned.Height = wh.Height
	cloned.Length = wh.Length
	cloned.Packages = copyMap(wh.Packages)
	cloned.ForkLifts = copyMap(wh.ForkLifts)
	cloned.Trucks = copyMap(wh.Trucks)

	return cloned
}

// EntityMap a map of entities
type EntityMap[T Package | ForkLift | Truck] map[Position]T

// Position a position in a 2D plane
// X the position in the X axis
// Y the position in the Y axis
type Position struct {
	X, Y int
}

// Package description of a Package
// Weight weight of the Package
// Name name of the Package
type Package struct {
	Weight Weight
	Name   string
}

// Weight a weight
type Weight int

// ForkLift description of a ForkLift
// Name name of the ForkLift
type ForkLift struct {
	Name string
	pack *Package
}

// Truck description of a Truck
// Name name of the Truck
// MaxWeight maximum Weight of the Truck
// CurrentWeight actual loaded Weight of the Truck
// ElapseDischargingTime how many cycles are needed for the Truck to return
// TimeUntilReturn the actual cycles left for the Truck to return
type Truck struct {
	Name                  string
	MaxWeight             Weight
	CurrentWeight         Weight
	ElapseDischargingTime int
	TimeUntilReturn       int
}

// Exists check if something exists at this Position on the EntityMap
func (ettMap EntityMap[T]) Exists(pos Position) bool {
	_, exists := ettMap[pos]
	return exists
}

// CleanWarehouse clean the Warehouse and populates the CycleState channel
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
	fullTrucks := make(map[Position]struct{})

	for index := 0; index < len(paths); {
		path := paths[index]

		if path.isValid() {
			forklift := wh.ForkLifts[path.current]
			delete(waitingForklifts, path.current)

			if len(path.steps) == 0 {
				if wh.Trucks.Exists(path.destination) {
					paths, index, events = dropPackage(path, forklift, index, wh.ForkLifts,
						wh.Trucks, paths, fullTrucks, events)
				} else if wh.Packages.Exists(path.destination) {
					paths, events = takePackage(path, forklift, index, wh.ForkLifts, wh.Packages,
						paths, events)
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

	events = processTrucks(wh, fullTrucks, events)

	return paths, events
}

func moveForkLift(path Path, forklift ForkLift, index int, forkLifts EntityMap[ForkLift],
	paths []Path, events []Event,
) ([]Path, []Event) {
	if !forkLifts.Exists(path.steps[0]) {
		events = append(events, ForkliftMove{
			forkliftName:  forklift.Name,
			eventPosition: path.current, target: path.steps[0],
		})

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
	events = append(events, PickupPackage{
		position: path.current, emitterName: forklift.Name,
		packName: packages[path.destination].Name,
	})

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
	trucks EntityMap[Truck], paths []Path, fulltrucks positionSet, events []Event,
) ([]Path, int, []Event) {
	truck := trucks[path.destination]

	if truck.CurrentWeight+forklift.pack.Weight <= truck.MaxWeight && truck.TimeUntilReturn == 0 {
		events = append(events, DeliverPackage{
			position: path.current, emitterName: forklift.Name,
			packName: forklift.pack.Name,
		})

		truck.CurrentWeight += forklift.pack.Weight
		forklift.pack = nil

		forkLifts[path.current] = forklift
		trucks[path.destination] = truck
		paths[index] = paths[len(paths)-1]
		paths = paths[:len(paths)-1]
	} else {
		events = append(events, ForkliftWait{forkliftName: forklift.Name, position: path.current})
		fulltrucks[path.destination] = struct{}{}

		index++
	}

	return paths, index, events
}

func processTrucks(wh Warehouse, fullTrucks positionSet, events []Event) []Event {
	for pos, truck := range wh.Trucks {
		if truck.TimeUntilReturn == 0 && truck.MaxWeight <= truck.CurrentWeight {
			fullTrucks[pos] = struct{}{}
		}
	}

	for pos := range fullTrucks {
		sendTruck(pos, wh.Trucks)
	}

	for pos, truck := range wh.Trucks {
		if truck.TimeUntilReturn == 0 {
			events = append(events, createTruckWait(truck, pos))
		} else {
			truck.TimeUntilReturn--
			if truck.TimeUntilReturn == 0 {
				truck.CurrentWeight = 0
			}

			wh.Trucks[pos] = truck

			if truck.TimeUntilReturn == truck.ElapseDischargingTime+1 || truck.TimeUntilReturn == 0 {
				events = append(events, createTruckWait(truck, pos))
			} else {
				events = append(events, createTruckGone(truck, pos))
			}
		}
	}

	return events
}

func sendTruck(pos Position, trucks EntityMap[Truck]) {
	truck := trucks[pos]

	if truck.TimeUntilReturn == 0 {
		truck.TimeUntilReturn = truck.ElapseDischargingTime + 1
		trucks[pos] = truck
	}
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
