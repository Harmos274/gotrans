package main

import (
	"fmt"
	"github.com/Harmos274/gotrans/clean_warehouse"
)

func ShowWarehouse(warehouse clean_warehouse.Warehouse) string {
	return fmt.Sprintf("%X\n", warehouse)
}
