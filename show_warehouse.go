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
		fmt.Println(e)
	}
	return output
}

func (sw ShowableWarehouse) String() string {
	return sw.Output() + sw.WarehouseMap()
}
