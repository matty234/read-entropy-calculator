package main

import (
	"fmt"
	"sort"
	"sync"
)

type FoundMers struct {
	Mer   string
	Count int
}

type MatchedMers struct {
	sync.Mutex

	Chastity float64

	OffsetMers map[int][]FoundMers
}

func CreateMatchedMers(chastity float64) MatchedMers {
	return MatchedMers{
		Mutex:      sync.Mutex{},
		Chastity:   chastity,
		OffsetMers: make(map[int][]FoundMers),
	}
}

func (m *MatchedMers) AddMer(offset int, mer string) {
	m.Mutex.Lock()

	offsetMers, ok := m.OffsetMers[offset]
	if !ok {
		offsetMers = make([]FoundMers, 0)
	}

	for k, v := range offsetMers {
		if v.Mer == mer {
			v.Count++
			m.OffsetMers[offset][k] = v
			m.Mutex.Unlock()
			return
		}
	}

	offsetMers = append(offsetMers, FoundMers{mer, 1})
	m.OffsetMers[offset] = offsetMers

	m.Mutex.Unlock()
}

func (m *MatchedMers) Summarise() {
	m.Mutex.Lock()

	for i := 0; i < len(m.OffsetMers); i++ {
		offsetMers := m.OffsetMers[i]
		fmt.Printf("Offset %d	count: %d\n", i, len(offsetMers))
	}

	// sort the offsets
	for i := 0; i < len(m.OffsetMers); i++ {
		unsortedMers, ok := m.OffsetMers[i]
		if !ok {
			continue
		}
		sort.Slice(unsortedMers, func(i, j int) bool {
			return unsortedMers[i].Mer > unsortedMers[j].Mer
		})

		for _, v := range unsortedMers {
			if float64(v.Count)/float64(len(unsortedMers)) > m.Chastity {
				fmt.Printf("%d\t%s\t%d\n", i, v.Mer, v.Count)
			}
		}
	}

	m.Mutex.Unlock()
}
