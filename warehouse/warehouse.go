package warehouse

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

	for ; cycles != 0; cycles -= 1 {
		// Pathfinding Algorithm
		ch <- currentWarehouse
	}
}

func copyMap[T Package | PalletJack | Truck](toClone map[Position]T) map[Position]T {
	ret := make(map[Position]T)
	for key, value := range toClone {
		ret[key] = value
	}
	return ret
}
