package main

import (
	"fmt"
	"strings"

	"github.com/Harmos274/gotrans/warehouse"
)

type ShowableWarehouse warehouse.CycleState

func (sw ShowableWarehouse) WarehouseMap() string {
	wr := sw.Warehouse
	w := strings.Repeat("#", wr.Length*2+2)
	for y := 0; y < wr.Height; y++ {
		w += "#\n# "
		for x := 0; x < wr.Length; x++ {
			pos := warehouse.Position{X: x, Y: y}
			if wr.Packages.Exists(pos) {
				w += "📦"
			} else if wr.ForkLifts.Exists(pos) {
				w += "👷"
			} else if wr.Trucks.Exists(pos) {
				w += "🚚"
			} else {
				w += "  "
			}
		}
	}
	w += "#\n" + strings.Repeat("#", wr.Length*2+3) + "\n"
	return w
}

func (sw ShowableWarehouse) Output() string {
	var output string
	for _, e := range sw.Events {
		// fmt.Println(e)
		switch e.(type) {
		case warehouse.PickupPackage:
			output += fmt.Sprintf("%s is is taking the package %s at position [%d,%d]\n", e.EmitterName(), e.PackageName(), e.AtPosition().X, e.AtPosition().Y)
		case warehouse.ForkliftWait:
			output += fmt.Sprintf("%s is waiting at position [%d,%d]\n", e.EmitterName(), e.AtPosition().X, e.AtPosition().Y)
		case warehouse.ForkliftMove:
			output += fmt.Sprintf("%s move from [%d,%d] to [%d,%d]\n", e.EmitterName(), e.AtPosition().X, e.AtPosition().Y, e.ToPosition().X, e.ToPosition().Y)
		case warehouse.DeliverPackage:
			output += fmt.Sprintf("%s is delivering the package %s\n", e.EmitterName(), e.PackageName())
		case warehouse.TruckWait:
			output += fmt.Sprintf("%s is waiting. %d/%d\n", e.EmitterName(), e.ChargedWeight(), e.MaxWeight())
		case warehouse.TruckGone:
			output += fmt.Sprintf("%s is gone. %d/%d\n", e.EmitterName(), e.ChargedWeight(), e.MaxWeight())
		default:
			output += ""
		}

	}
	return output
}

func (sw ShowableWarehouse) String() string {
	return sw.Output() + sw.WarehouseMap()
}
