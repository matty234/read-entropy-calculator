package main

import (
	"fmt"
	"log"
	"sort"
	"sync"
)

type FoundMers struct {
	Mer   string
	Count int
}

type MatchedMers struct {
	sync.Map

	Chastity float64

	//OffsetMers map[int][]FoundMers
}

func CreateMatchedMers(chastity float64) MatchedMers {
	return MatchedMers{
		Map:      sync.Map{},
		Chastity: chastity,
	}
}

func (m *MatchedMers) AddMer(offset int, mer string) {

	offsetMers, ok := m.Map.Load(offset)
	if !ok {
		offsetMers = make([]FoundMers, 0)
	}

	o := offsetMers.([]FoundMers)

	for k, v := range o {
		if v.Mer == mer {
			v.Count++
			o[k] = v
			return
		}
	}

	offsetMers = append(o, FoundMers{mer, 1})
	m.Map.Store(offset, offsetMers)

}

func (m *MatchedMers) Summarise() {

	length := 0
	m.Map.Range(func(key, value interface{}) bool {
		l := key.(int)
		if l > length {
			length = l
		}
		return true
	})

	for i := 0; i < length; i++ {
		offset := i
		me, ok := m.Map.Load(i)
		if !ok {
			log.Fatal("Overrun?")
		}

		mers := me.([]FoundMers)

		sort.Slice(mers, func(i, j int) bool {
			return mers[i].Count > mers[j].Count
		})

		fmt.Printf("Offset %d\n", offset)
		for _, v := range mers {
			if float64(v.Count)/float64(length) > m.Chastity {
				fmt.Printf("\t%s\t%d\n", v.Mer, v.Count)
			}
		}
	}

}
