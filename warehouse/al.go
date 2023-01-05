package warehouse

import "fmt"

const (
	UP = 1 << iota
	RIGHT
	DOWN
	LEFT
	ALL = UP | RIGHT | DOWN | LEFT
)

func Test(wh Warehouse) {
	paths := getPaths(wh)

	fmt.Println(paths)
}

func getPaths(wh Warehouse) []Path {
	distances := make([]Path, 0, min(len(wh.Packages), len(wh.PalletJacks)))

	for pos := range wh.PalletJacks {
		fmt.Println(pathToNearestPackage(wh, pos, distances))
	}

	return distances
}

func pathToNearestPackage(wh Warehouse, start Position, other_paths []Path) Path {
	best_path := Path{current: start, destination: start}
	current_path := []AttemptPosition{{Position: start, directions: ALL}}
	pos_set := make(map[Position]struct{})
	loops := 0
	path_max := 0

	for len(current_path) > 0 {
		loops += 1
		current_step := &current_path[len(current_path)-1]
		next_move := nextDirection(current_step.directions)

		if next_move != 0 && shouldKeepSearching(current_path, best_path) {
			current_step.directions ^= next_move
			move(wh, current_step.Position, next_move, &current_path, pos_set, &best_path,
				other_paths)
		} else {
			current_path = current_path[:len(current_path)-1]
		}

		if len(current_path) > path_max {
			path_max = len(current_path)
		}
	}

	fmt.Println(loops)
	fmt.Println(path_max)

	return best_path
}

func move(wh Warehouse, current_pos Position, direction int, path *[]AttemptPosition,
	pos_set positionSet, current_best *Path, other_paths []Path) {
	new_pos, possible := getNewPos(current_pos, direction, wh.Length, wh.Height)

	if !possible || alreadyVisited(pos_set, new_pos) {
		return
	}

	if wh.SomethingExistsAt(new_pos) {
		if wh.Packages.Exists(new_pos) && shouldTakePackage(new_pos, other_paths) {
			*current_best = attemptToPath(*path, new_pos)
		} else {
			return
		}
	} else {
		*path = append(*path, AttemptPosition{
			directions: ALL ^ inverseDirection(direction),
			Position:   new_pos,
		})
		pos_set[new_pos] = struct{}{}
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

func shouldTakePackage(pos Position, other_paths []Path) bool {
	return true
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

type positionSet = map[Position]struct{}

func alreadyVisited(set positionSet, pos Position) bool {
	_, has := set[pos]

	return has
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

func isUsed(pos Position, paths []Path) bool {
	for _, path := range paths {
		if pos == path.current || pos == path.destination {
			return true
		}
	}

	return false
}

func (self Path) isValid() bool {
	return self.current != self.destination
}

func distance(lhs Position, rhs Position) int {
	return abs(lhs.X-rhs.X) + abs(lhs.Y-rhs.Y)
}

// The 80s called, they want their programming language design back
func min(lhs int, rhs int) int {
	if lhs < rhs {
		return lhs
	} else {
		return rhs
	}
}

func abs(nbr int) int {
	if nbr < 0 {
		return -nbr
	} else {
		return nbr
	}
}
