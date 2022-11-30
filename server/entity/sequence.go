package entity

import (
	"fmt"
	"regexp"
)

type Sequences struct {
	Letters []string `json:"letters"`
}

func (s Sequences) Validate() (bool, error) {
	graph := make([][]string, len(s.Letters))
	rg, _ := regexp.Compile(`^[BUDH]+$`)

	for i, val := range s.Letters {
		// checks if sequences contain only valid letters and length
		if match := rg.MatchString(val); !match || len(val) != len(s.Letters) {
			return false, fmt.Errorf("invalid input: '%s'", val)
		}
		graph[i] = []string{val}
	}

	if countValidSequences(graph) >= 2 {
		return true, nil
	}

	return false, nil
}

func countValidSequences(graph [][]string) int {
	validSeqs := 0

	seen := make([][]bool, len(graph))
	for row := range seen {
		seen[row] = make([]bool, len(graph[0]))
	}

	// run a DFS in each graph cell
	for x := 0; x < len(graph); x++ {
		for y := 0; y < len(graph[0]); y++ {
			if seen[x][y] {
				continue
			}
			if dfs(graph, seen, x, y, 1) >= 4 {
				validSeqs++
			}
		}
	}

	return validSeqs
}

// TODO: check for direction change. If it changes, return counter and reset last seen
func dfs(graph [][]string, seen [][]bool, x, y, counter int) int {
	if seen[x][y] {
		return counter
	}
	seen[x][y] = true

	// traverse in all eight directions (top, bottom, left, right and 4 diagonals)
	for _, direction := range getDirections() {
		dx, dy := direction[0], direction[1]
		if isValidBoundary(x+dx, y+dy, len(graph)) {
			if graph[x][y] == graph[x+dx][y+dy] {
				counter++
				dfs(graph, seen, x+dx, y+dy, counter)
			}
		}
	}
	return counter
}

func getDirections() [][]int {
	return [][]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
		{-1, -1},
		{1, 1},
		{-1, 1},
		{1, -1},
	}
}

func isValidBoundary(x, y, graphLen int) bool {
	return x < 0 || y < 0 || x >= graphLen || y >= graphLen
}
