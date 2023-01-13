package warehouse

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
	truckName          string
	truckMaxWeight     Weight
	truckChargedWeight Weight
	position           Position
}

func (t TruckWait) EmitterName() string {
	return t.truckName
}

func (t TruckWait) AtPosition() Position {
	return t.position
}

func (t TruckWait) ChargedWeight() Weight {
	return t.truckChargedWeight
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
