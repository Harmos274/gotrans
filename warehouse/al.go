package warehouse

import "fmt"

const (
	UP = 1 << iota
	RIGHT
	DOWN
	LEFT
	ALL = UP | RIGHT | DOWN | LEFT
)

func refreshPaths(wh Warehouse, current_paths []Path) []Path {
	targeted_packages := countTargetedPackages(wh, current_paths)
	idle := getIdleForklifts(wh.PalletJacks, current_paths, true)

	for pos := range idle {
		path := pathToObject(wh, pos, current_paths, shouldGoToTruck)

		if path.isValid() {
			current_paths = append(current_paths, path)
		}
	}

	idle = getIdleForklifts(wh.PalletJacks, current_paths, false)

	fmt.Println("Idle empty forklifts", idle)
	fmt.Println(len(idle), "> 0 &&", targeted_packages, "<", len(wh.Packages))
	for len(idle) > 0 && targeted_packages < len(wh.Packages) {
		pos := idle.randomElem()

		path := pathToObject(wh, pos, current_paths, shouldGoToPackage)
		replaced := false

		if path.isValid() {
			current_paths, replaced = insertPath(path, current_paths)
		} else {
			fmt.Println("Could not find path for", pos, "returned", path)
		}

		if replaced {
			idle = getIdleForklifts(wh.PalletJacks, current_paths, false)
		} else {
			delete(idle, pos)
		}
	}

	return current_paths
}

func pathToObject(wh Warehouse, start Position, other_paths []Path, validator validator) Path {
	best_path := Path{current: start, destination: start}
	current_path := []AttemptPosition{{Position: start, directions: ALL}}
	pos_set := make(map[Position]struct{})

	for len(current_path) > 0 {
		current_step := &current_path[len(current_path)-1]
		next_move := nextDirection(current_step.directions)

		if next_move != 0 && shouldKeepSearching(current_path, best_path) {
			current_step.directions ^= next_move
			move(wh, current_step.Position, next_move, &current_path, pos_set, &best_path,
				other_paths, validator)
		} else {
			//fmt.Println("Walking back from", current_step.Position.X, current_step.Position.Y, "best path is", best_path)
			current_path = current_path[:len(current_path)-1]
			delete(pos_set, current_step.Position)
		}
	}

	fmt.Println("Computed best path for", start, best_path)

	return best_path
}

func move(wh Warehouse, current_pos Position, direction int, path *[]AttemptPosition,
	pos_set positionSet, current_best *Path, other_paths []Path, validator validator) {
	new_pos, possible := getNewPos(current_pos, direction, wh.Length, wh.Height)

	if !possible || alreadyVisited(pos_set, new_pos) ||
		isSomeoneOnThisTile(new_pos, other_paths, len(*path)-1) {
		return
	}

	if wh.SomethingExistsAt(new_pos) {
		if validator(wh, *path, new_pos, other_paths) {
			*current_best = attemptToPath(*path, new_pos)
			//fmt.Println("Found best", *current_best)
		} else {
			//fmt.Println("Collided with invalid object at", new_pos)
			return
		}
	} else {
		*path = append(*path, AttemptPosition{
			directions: ALL ^ inverseDirection(direction),
			Position:   new_pos,
		})
		pos_set[new_pos] = struct{}{}
		//fmt.Println("Moved to", new_pos)
	}
}

func nextDirection(directions int) int {
	if directions&UP != 0 {
		return UP
	} else if directions&RIGHT != 0 {
		return RIGHT
	} else if directions&DOWN != 0 {
		return DOWN
	} else if directions&LEFT != 0 {
		return LEFT
	} else {
		return 0
	}
}

func getNewPos(pos Position, direction int, size_x int, size_y int) (Position, bool) {
	if (direction == UP && pos.Y == 0) || (direction == DOWN && pos.Y == size_y-1) ||
		(direction == LEFT && pos.X == 0) || (direction == RIGHT && pos.X == size_x-1) {
		return Position{}, false
	}

	if direction == UP {
		pos.Y -= 1
	} else if direction == DOWN {
		pos.Y += 1
	} else if direction == LEFT {
		pos.X -= 1
	} else if direction == RIGHT {
		pos.X += 1
	}

	return pos, true
}

func attemptToPath(attempt []AttemptPosition, destination Position) Path {
	path := Path{current: attempt[0].Position, destination: destination}

	for pos := 1; pos < len(attempt); pos += 1 {
		path.steps = append(path.steps, attempt[pos].Position)
	}

	return path
}

func inverseDirection(direction int) int {
	if direction == UP {
		return DOWN
	} else if direction == DOWN {
		return UP
	} else if direction == LEFT {
		return RIGHT
	} else if direction == RIGHT {
		return LEFT
	} else {
		return 0
	}
}

func shouldGoToTruck(wh Warehouse, _ []AttemptPosition, pos Position, _ []Path) bool {
	return wh.Trucks.Exists(pos)
}

func shouldGoToPackage(wh Warehouse, attempt []AttemptPosition, pos Position, other_paths []Path) bool {
	if !wh.Packages.Exists(pos) {
		return false
	}

	existing := getExistingPathTo(pos, other_paths)

	return !existing.isValid() || len(attempt) < len(existing.steps)+1
}

func shouldKeepSearching(attempt []AttemptPosition, current_best Path) bool {
	return (!current_best.isValid() || len(attempt) < len(current_best.steps)+1)
}

type Path struct {
	current     Position
	destination Position
	steps       []Position
}

type AttemptPosition struct {
	Position
	directions int
}

type positionSet map[Position]struct{}

type validator = func(Warehouse, []AttemptPosition, Position, []Path) bool

func (self positionSet) randomElem() Position {
	for pos := range self {
		return pos
	}

	return Position{}
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

func (self Path) contains(pos Position) bool {
	if self.current == pos || self.destination == pos {
		return true
	}

	for _, step := range self.steps {
		if pos == step {
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
			packages += 1
		}
	}

	return packages
}

func getIdleForklifts(forklifs EntityMap[PalletJack], paths []Path, loaded bool) positionSet {
	idle_set := make(map[Position]struct{})

	for pos, forklift := range forklifs {

		if (forklift.pack != nil) == loaded {
			idle_set[pos] = struct{}{}
		}
	}

	for _, path := range paths {
		delete(idle_set, path.current)
	}

	return idle_set
}

func (self Path) isValid() bool {
	return self.current != self.destination
}
