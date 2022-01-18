package main

import (
	"fmt"
	"sync/atomic"
)

type merCount map[string]*uint64

type MatchedMers struct {
	mcs []merCount

	Chastity float64
}

func CreateMatchedMers(chastity float64, total int) MatchedMers {
	m := make([]merCount, total)

	for i := 0; i < total; i++ {
		m[i] = make(merCount)
	}

	return MatchedMers{
		mcs:      m,
		Chastity: chastity,
	}
}

func (m *MatchedMers) AddMer(offset int, mer string) {
	r := m.mcs[offset]

	if r[mer] == nil {
		r[mer] = new(uint64)
	}

	atomic.AddUint64(r[mer], 1)
}

func (m *MatchedMers) Summarise() {

	for i := 0; i < len(m.mcs); i++ {
		offset := i

		for k, v := range m.mcs[offset] {
			if float64(*v)/float64(len(m.mcs)) > m.Chastity {
				fmt.Printf("%d\t%s\t%d\n", i, k, *v)
			}
		}

	}

}
