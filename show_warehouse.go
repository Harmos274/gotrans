package main

import (
	"fmt"
	"strings"

	"github.com/Harmos274/gotrans/warehouse"
)

type ShowableWarehouse warehouse.Warehouse

func (sw ShowableWarehouse) WarehouseMap() string {
	w := strings.Repeat("#", sw.Length*2+2)
	for y := 0; y < sw.Height; y++ {
		w += "#\n# "
		for x := 0; x < sw.Length; x++ {
			pos := warehouse.Position{X: x, Y: y}
			if sw.Packages.Exists(pos) {
				w += "P "
			} else if sw.ForkLifts.Exists(pos) {
				w += "J "
			} else if sw.Trucks.Exists(pos) {
				w += "T "
			} else {
				w += "  "
			}
		}
	}
	w += "#\n" + strings.Repeat("#", sw.Length*2+3) + "\n"
	return w
}

func (sw ShowableWarehouse) Output() string {
	var output string
	for pos, forkLift := range sw.ForkLifts {
		output += fmt.Sprintf("%s %s", forkLift.Name, forkLift.State)
		switch forkLift.State {
		case "GO":
			output += fmt.Sprintf(" [%d,%d]\n", pos.X, pos.Y)
		case "TAKE":
			// output += fmt.Sprintf(" %s %v\n", forkLift.pack.Name, forkLift.pack.Weight)
		case "LEAVE":
			// output += fmt.Sprintf(" %s %v\n", forkLift.pack.Name, forkLift.pack.Weight)
		default:
			output += "\n"
		}
	}
	for _, truck := range sw.Trucks {
		output += fmt.Sprintf("%s %s %d/%d\n", truck.Name, truck.State, truck.CurrentWeight, truck.MaxWeight)
	}
	return output
}

func (sw ShowableWarehouse) String() string {
	return sw.Output()
}
