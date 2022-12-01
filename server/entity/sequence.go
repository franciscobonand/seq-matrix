package entity

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
)

// ! corner case: "DDDDBDDDD" only counts as one (in every direction this happens)

type Sequences struct {
	Letters []string `json:"letters"`
}

func (s Sequences) Validate() (bool, error) {
	rg, _ := regexp.Compile(`^[BUDH]+$`)

	for _, val := range s.Letters {
		// checks if sequences contain only valid letters and length
		if match := rg.MatchString(val); !match || len(val) != len(s.Letters) {
			return false, fmt.Errorf("invalid input: '%s'", val)
		}
	}

	return countValidSequences(s.Letters) >= 2, nil
}

func countValidSequences(seqs []string) uint32 {
	var wg sync.WaitGroup
	var validSeqs atomic.Uint32

	wg.Add(4)
	checkRows(&wg, &validSeqs, seqs)
	checkColumns(&wg, &validSeqs, seqs)
	checkPrimaryDiag(&wg, &validSeqs, seqs)
	checkSecondaryDiag(&wg, &validSeqs, seqs)
	wg.Wait()

	return validSeqs.Load()
}

// checkRows returns number of repetitions in each line (rows)
func checkRows(wg *sync.WaitGroup, reps *atomic.Uint32, seqs []string) {
	defer wg.Done()

	for _, seq := range seqs {
		if checkForRepetitions(seq) {
			reps.Add(1)
		}
	}
}

// checkColumns returns number of repetitions in transposed matrix (columns)
func checkColumns(wg *sync.WaitGroup, reps *atomic.Uint32, seqs []string) {
	defer wg.Done()
	str := ""

	for i := 0; i < len(seqs); i++ {
		for _, seq := range seqs {
			str = str + string(seq[i])
		}
		if checkForRepetitions(str) {
			reps.Add(1)
		}
		str = ""
	}
}

// checkSecondaryDiag returns number of repetitions in the secondary diagonals with length >= 4
func checkSecondaryDiag(wg *sync.WaitGroup, reps *atomic.Uint32, seqs []string) {
	defer wg.Done()
	arrSize := len(seqs)
	str := ""

	for k := arrSize - (arrSize - 4) - 1; k <= arrSize-1; k++ {
		i := k
		j := 0
		for {
			if i < 0 {
				break
			}
			str = str + string(seqs[i][j])
			i = i - 1
			j = j + 1
		}
		if checkForRepetitions(str) {
			reps.Add(1)
		}
		str = ""
	}

	for k := 1; k <= arrSize-4; k++ {
		i := arrSize - 1
		j := k
		for {
			if j > arrSize-1 {
				break
			}
			str = str + (string(seqs[i][j]))
			i = i - 1
			j = j + 1
		}
		if checkForRepetitions(str) {
			reps.Add(1)
		}
		str = ""
	}
}

// checkPrimaryDiag returns number of repetitions in the primary diagonals with length >= 4
func checkPrimaryDiag(wg *sync.WaitGroup, reps *atomic.Uint32, seqs []string) {
	defer wg.Done()
	arrSize := len(seqs)
	str := ""

	for k := 1; k <= arrSize-4; k++ {
		i := k
		j := 0
		for {
			if i > arrSize-1 {
				break
			}
			str = str + string(seqs[i][j])
			i = i + 1
			j = j + 1
		}
		if checkForRepetitions(str) {
			reps.Add(1)
		}
		str = ""
	}

	for k := 0; k <= arrSize-4; k++ {
		i := 0
		j := k
		for {
			if j > arrSize-1 {
				break
			}
			str = str + (string(seqs[i][j]))
			i = i + 1
			j = j + 1
		}
		if checkForRepetitions(str) {
			reps.Add(1)
		}
		str = ""
	}
}

func getPossibleRepetitions() []string {
	return []string{
		"BBBB",
		"UUUU",
		"DDDD",
		"HHHH",
	}
}

func checkForRepetitions(str string) bool {
	for _, rep := range getPossibleRepetitions() {
		if strings.Contains(str, rep) {
			return true
		}
	}
	return false
}
