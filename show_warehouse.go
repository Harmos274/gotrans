package main

import (
	"fmt"
	"github.com/Harmos274/gotrans/warehouse"
)

type ShowableWarehouse warehouse.Warehouse

func (sw ShowableWarehouse) String() string {
	// Faut changer ça c'est pour la compilation !!!!!!
	return fmt.Sprintf("%X\n", (warehouse.Warehouse)(sw))
}
