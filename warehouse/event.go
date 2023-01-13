package warehouse

import "fmt"

// Event events occurring while cleaning the warehouse
// EmitterName Returns event's emitter name
// AtPosition Returns event's emission position
type Event interface {
	EmitterName() string
	AtPosition() Position
}

type PickupPackage struct {
	position    Position
	emitterName string
	packName    string
}

func (pickPack PickupPackage) String() string {
	return fmt.Sprintf("%s is taking the package %s at position [%d,%d]", pickPack.emitterName, pickPack.packName, pickPack.position.X, pickPack.position.Y)
}

func (pickPack PickupPackage) EmitterName() string {
	return pickPack.emitterName
}

func (pickPack PickupPackage) AtPosition() Position {
	return pickPack.position
}

func (pickPack PickupPackage) PackageName() string {
	return pickPack.packName
}

type ForkliftWait struct {
	forkliftName string
	position     Position
}

func (fw ForkliftWait) String() string {
	return fmt.Sprintf("%s is waiting at position [%d,%d]", fw.forkliftName, fw.position.X, fw.position.Y)
}

func (fw ForkliftWait) EmitterName() string {
	return fw.forkliftName
}

func (fw ForkliftWait) AtPosition() Position {
	return fw.position
}

type ForkliftMove struct {
	forkliftName  string
	eventPosition Position
	target        Position
}

func (f ForkliftMove) String() string {
	return fmt.Sprintf("%s move from [%d,%d] to [%d,%d]", f.forkliftName, f.eventPosition.X, f.eventPosition.Y, f.target.X, f.target.Y)
}

func (f ForkliftMove) EmitterName() string {
	return f.forkliftName
}

func (f ForkliftMove) AtPosition() Position {
	return f.eventPosition
}

func (f ForkliftMove) ToPosition() Position {
	return f.target
}

type DeliverPackage struct {
	position    Position
	emitterName string
	packName    string
}

func (d DeliverPackage) String() string {
	return fmt.Sprintf("%s is delivering the package %s", d.emitterName, d.packName)
}

func (d DeliverPackage) EmitterName() string {
	return d.emitterName
}

func (d DeliverPackage) AtPosition() Position {
	return d.position
}

func (d DeliverPackage) PackageName() string {
	return d.packName
}

type TruckWait struct {
	truckName         string
	truckMaxWeight    Weight
	truckLoadedWeight Weight
	position          Position
}

func (t TruckWait) String() string {
	return fmt.Sprintf("%s is waiting. %d/%d", t.truckName, t.truckLoadedWeight, t.truckMaxWeight)
}

func (t TruckWait) EmitterName() string {
	return t.truckName
}

func (t TruckWait) AtPosition() Position {
	return t.position
}

func (t TruckWait) ChargedWeight() Weight {
	return t.truckLoadedWeight
}

func (t TruckWait) MaxWeight() Weight {
	return t.truckMaxWeight
}

type TruckGone struct {
	truckName          string
	truckMaxWeight     Weight
	truckChargedWeight Weight
	position           Position
}

func (t TruckGone) String() string {
	return fmt.Sprintf("%s is gone. %d/%d", t.truckName, t.truckChargedWeight, t.truckMaxWeight)
}

func (t TruckGone) EmitterName() string {
	return t.truckName
}

func (t TruckGone) AtPosition() Position {
	return t.position
}

func (t TruckGone) ChargedWeight() Weight {
	return t.truckChargedWeight
}

func (t TruckGone) MaxWeight() Weight {
	return t.truckMaxWeight
}
