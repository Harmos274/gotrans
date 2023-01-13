// Package warehouse Contains the warehouse data structure and related utilities
package warehouse

import (
	"fmt"
	"os"
)

const (
	nONE = iota
	uP
	rIGHT
	dOWN
	lEFT
)

// Path represent a ForkLift's Path with its current Position its target and all the steps
type Path struct {
	current     Position
	destination Position
	steps       []Position
}

type direction = int

func refreshPaths(wh Warehouse, currentPaths []Path) []Path {
	targetedPackages := countTargetedPackages(wh, currentPaths)
	idle := getIdleForklifts(wh.ForkLifts, currentPaths, true)
	trucks := mapToPositionSet(wh.Trucks)
	for pos := range idle {
		truckValidator := func(_ []attemptPosition, targetPos Position) bool {
			return wh.Trucks[targetPos].MaxWeight >= wh.ForkLifts[pos].pack.Weight
		}
		path := pathToObject(wh, pos, trucks, currentPaths, truckValidator)

		if path.isValid() {
			currentPaths = append(currentPaths, path)
		}
	}

	idle = getIdleForklifts(wh.ForkLifts, currentPaths, false)
	packages := mapToPositionSet(wh.Packages)

	for len(idle) > 0 && targetedPackages < len(wh.Packages) {
		pos := idle.randomElem()
		packageValidator := func(path []attemptPosition, pos Position) bool {
			return shouldGoToPackage(wh, path, pos, currentPaths)
		}

		path := pathToObject(wh, pos, packages, currentPaths, packageValidator)
		replaced := false

		if path.isValid() {
			currentPaths, replaced = insertPath(path, currentPaths)
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "Could not find path for %v, returned %v", pos, path)
		}

		if replaced {
			idle = getIdleForklifts(wh.ForkLifts, currentPaths, false)
		} else {
			delete(idle, pos)
		}
	}

	return currentPaths
}

func pathToObject(wh Warehouse, start Position, targets positionSet, otherPaths []Path, validator validator) Path {
	bestPath := Path{current: start, destination: start}
	nearestTarget := getNearestEntityPos(start, targets)
	currentPath := []attemptPosition{{
		Position:   start,
		directions: getDirectionsPriority(start, nearestTarget, nONE), nextDirection: 0,
	}}
	posSet := make(map[Position]struct{})

	for len(currentPath) > 0 {
		currentStep := &currentPath[len(currentPath)-1]

		for currentStep.nextDirection < 4 &&
			currentStep.directions[currentStep.nextDirection] == nONE {
			currentStep.nextDirection++
		}

		if currentStep.nextDirection < 4 && shouldKeepSearching(currentPath, bestPath) {
			nextMove := currentStep.directions[currentStep.nextDirection]
			move(wh, currentStep.Position, nextMove, &currentPath, posSet, &bestPath,
				otherPaths, validator, targets)
			currentStep.nextDirection++
		} else {
			currentPath = currentPath[:len(currentPath)-1]
			delete(posSet, currentStep.Position)
		}
	}

	return bestPath
}

func move(wh Warehouse, currentPos Position, direction int, path *[]attemptPosition,
	posSet positionSet, currentBest *Path, otherPaths []Path, validator validator, targets positionSet,
) {
	newPos, possible := getNewPos(currentPos, direction, wh.Length, wh.Height)

	if !possible || alreadyVisited(posSet, newPos) ||
		isSomeoneOnThisTile(newPos, otherPaths, len(*path)-1) {
		return
	}

	if targets.has(newPos) {
		if validator(*path, newPos) {
			*currentBest = attemptToPath(*path, newPos)
		}
		delete(targets, newPos)
	} else if wh.SomethingExistsAt(newPos) {
		return
	} else {
		*path = append(*path, attemptPosition{
			directions: getDirectionsPriority(newPos, getNearestEntityPos(newPos, targets), inverseDirection(direction)),
			Position:   newPos,
		})
		posSet[newPos] = struct{}{}
	}
}

func getDirectionsPriority(current Position, target Position,
	comingFrom direction,
) [4]direction {
	distanceX := target.X - current.X
	distanceY := target.Y - current.Y
	priorityX := [2]direction{rIGHT, lEFT}
	priorityY := [2]direction{dOWN, uP}

	switch comingFrom {
	case rIGHT:
		priorityX[0] = nONE
	case lEFT:
		priorityX[1] = nONE
	case dOWN:
		priorityY[0] = nONE
	case uP:
		priorityY[1] = nONE
	}

	if distanceX < 0 {
		priorityX = [2]direction{priorityX[1], priorityX[0]}
	}
	if distanceY < 0 {
		priorityY = [2]direction{priorityY[1], priorityY[0]}
	}

	positions := [4]direction{priorityX[0], priorityY[0], priorityX[1], priorityY[1]}

	if abs(distanceY) > abs(distanceX) {
		positions = [4]direction{priorityY[0], priorityX[0], priorityY[1], priorityX[1]}
	}

	return positions
}

func getNewPos(pos Position, direction int, sizeX int, sizeY int) (Position, bool) {
	switch direction {
	case uP:
		pos.Y--
	case dOWN:
		pos.Y++
	case lEFT:
		pos.X--
	case rIGHT:
		pos.X++
	}

	valid := pos.X >= 0 && pos.X < sizeX && pos.Y >= 0 && pos.Y < sizeY

	return pos, valid
}

func getNearestEntityPos(pos Position, entitiesPos positionSet) Position {
	nearest := entitiesPos.randomElem()
	nearestDistance := euclideanDistance(pos, nearest)

	for entityPos := range entitiesPos {
		currentDistance := euclideanDistance(pos, entityPos)

		if currentDistance < nearestDistance {
			nearest = entityPos
			nearestDistance = currentDistance
		}
	}

	return nearest
}

func euclideanDistance(lhs Position, rhs Position) int {
	return abs(lhs.X-rhs.X) + abs(lhs.Y-rhs.Y)
}

func attemptToPath(attempt []attemptPosition, destination Position) Path {
	path := Path{current: attempt[0].Position, destination: destination}

	for pos := 1; pos < len(attempt); pos++ {
		path.steps = append(path.steps, attempt[pos].Position)
	}

	return path
}

func inverseDirection(direction int) int {
	switch direction {
	case uP:
		return dOWN
	case dOWN:
		return uP
	case lEFT:
		return rIGHT
	case rIGHT:
		return lEFT
	default:
		return 0
	}
}

func shouldGoToPackage(wh Warehouse, attempt []attemptPosition, pos Position, otherPaths []Path) bool {
	if !wh.Packages.Exists(pos) {
		return false
	}

	existing := getExistingPathTo(pos, otherPaths)

	return !existing.isValid() || len(attempt) <= len(existing.steps)
}

func shouldKeepSearching(attempt []attemptPosition, currentBest Path) bool {
	return !currentBest.isValid() || len(attempt) <= len(currentBest.steps)
}

type attemptPosition struct {
	Position
	directions    [4]direction
	nextDirection int
}

type positionSet map[Position]struct{}

type validator = func([]attemptPosition, Position) bool

func (set positionSet) randomElem() Position {
	for pos := range set {
		return pos
	}

	return Position{}
}

func (set positionSet) has(pos Position) bool {
	_, ok := set[pos]

	return ok
}

func alreadyVisited(set positionSet, pos Position) bool {
	_, has := set[pos]

	return has
}

func isSomeoneOnThisTile(pos Position, paths []Path, turn int) bool {
	for _, path := range paths {
		if len(path.steps) > turn && path.steps[turn] == pos {
			return true
		}
	}

	return false
}

func getExistingPathTo(pos Position, paths []Path) Path {
	for _, path := range paths {
		if pos == path.destination {
			return path
		}
	}

	return Path{}
}

func insertPath(path Path, paths []Path) ([]Path, bool) {
	for index, existing := range paths {
		if existing.current == path.current {
			paths[index] = path
			return paths, true
		}
	}

	return append(paths, path), false
}

func countTargetedPackages(wh Warehouse, paths []Path) int {
	packages := 0

	for _, path := range paths {
		if wh.Packages.Exists(path.destination) {
			packages++
		}
	}

	return packages
}

func getIdleForklifts(forklifts EntityMap[ForkLift], paths []Path, loaded bool) positionSet {
	idleSet := make(map[Position]struct{})

	for pos, forklift := range forklifts {
		if (forklift.pack != nil) == loaded {
			idleSet[pos] = struct{}{}
		}
	}

	for _, path := range paths {
		delete(idleSet, path.current)
	}

	return idleSet
}

func mapToPositionSet[T Package | Truck | ForkLift](entities EntityMap[T]) positionSet {
	set := make(map[Position]struct{}, len(entities))

	for pos := range entities {
		set[pos] = struct{}{}
	}

	return set
}

func (path Path) isValid() bool {
	return path.current != path.destination
}

func abs(value int) int {
	if value >= 0 {
		return value
	}
	return -value
}
