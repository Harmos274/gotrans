package main

import (
	"fmt"
	"strings"

	"github.com/Harmos274/gotrans/warehouse"
)

type ShowableWarehouse warehouse.Warehouse

func (sw ShowableWarehouse) WarehouseMap() string {
	w := strings.Repeat("#", int(sw.Length)*2+2)
	for y := 0; y < int(sw.Height); y++ {
		w += "#\n# "
		for x := 0; x < int(sw.Length); x++ {
			pos := warehouse.Position{X: x, Y: y}
			if sw.Packages.Exists(pos) {
				w += "P "
			} else if sw.PalletJacks.Exists(pos) {
				w += "J "
			} else if sw.Trucks.Exists(pos) {
				w += "T "
			} else {
				w += "  "
			}
		}
	}
	w += "#\n" + strings.Repeat("#", int(sw.Length)*2+3) + "\n"
	return w
}

func (sw ShowableWarehouse) Output() string {
	var output string
	for pos, palletJack := range sw.PalletJacks {
		output += fmt.Sprintf("%s %s", palletJack.Name, palletJack.State)
		switch palletJack.State {
		case "GO":
			output += fmt.Sprintf(" [%d,%d]\n", pos.X, pos.Y)
		case "TAKE":
			// output += fmt.Sprintf(" %s %v\n", palletJack.pack.Name, palletJack.pack.Weight)
		case "LEAVE":
			// output += fmt.Sprintf(" %s %v\n", palletJack.pack.Name, palletJack.pack.Weight)
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
