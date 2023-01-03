package warehouse

type Warehouse struct {
	Length, Height uint
	Packages       EntityMap[Package]
	PalletJacks    EntityMap[PalletJack]
	Trucks         EntityMap[Truck]
}

func (wh Warehouse) SomethingExistsAtThisPosition(pos Position) bool {
	return wh.Packages.Exists(pos) || wh.PalletJacks.Exists(pos) || wh.Trucks.Exists(pos)
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
