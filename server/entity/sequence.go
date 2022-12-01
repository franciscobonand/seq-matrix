package entity

import (
	"fmt"
	"regexp"
	"strings"
)

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

func countValidSequences(seqs []string) int {
	validSeqs := 0

	// rows
	for _, seq := range seqs {
		for _, rep := range getPossibleRepetitions() {
			if strings.Contains(seq, rep) {
				validSeqs++
			}
		}
	}

	// transpose for columns
	// secondaryDiag for secondary diagonal values
	// TODO: create principalDiag and add parallelism

	return validSeqs
}

func getPossibleRepetitions() []string {
	return []string{
		"BBBB",
		"UUUU",
		"DDDD",
		"HHHH",
	}
}

func transpose(seqs []string) []string {
	transposed := []string{}
	str := ""
	for i := 0; i < len(seqs); i++ {
		for _, seq := range seqs {
			str = str + string(seq[i])
		}
		transposed = append(transposed, str)
		str = ""
	}
	return transposed
}

func secondaryDiag(seqs []string) []string {
	secDiag := []string{}
	arrSize := len(seqs)
	str := ""
	for k := arrSize - (arrSize - 4) - 1; k < arrSize-1; k++ {
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
		secDiag = append(secDiag, str)
		str = ""
	}

	str = ""
	for i := 0; i < arrSize; i++ {
		str = str + (string(seqs[i][i]))
	}
	secDiag = append(secDiag, str)
	str = ""

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
		secDiag = append(secDiag, str)
		str = ""
	}

	return secDiag
}
