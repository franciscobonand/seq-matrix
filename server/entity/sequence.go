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
		if match := rg.MatchString(val); !match {
			return false, fmt.Errorf("invalid input: '%s'", val)
		}
		graph[i] = []string{val}
	}

	return true, nil
}
