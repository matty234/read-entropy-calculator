package main

import (
	"log"
	"sync"
)

type MerMatchCollection struct {
	mcs []*MerMatch

	wg *sync.WaitGroup
}

func CreateMerMatchCollection(chastity float64, k, consideredBases int, m *MatchedMers) MerMatchCollection {
	var mcs []*MerMatch = make([]*MerMatch, 0)
	wait := &sync.WaitGroup{}
	for i := 0; i+k < consideredBases; i++ {
		wait.Add(1)
		log.Println("Creating new MerMatch for offset", i)
		m := CreateMerMatch(k, i, m, wait)
		mcs = append(mcs, &m)
	}
	return MerMatchCollection{
		mcs: mcs,
		wg:  wait,
	}
}

func (m *MerMatchCollection) Start() {
	for _, v := range m.mcs {
		go v.FindMers()
	}
}

func (m *MerMatchCollection) Broadcast(mer string) {
	for _, v := range m.mcs {
		v.c <- mer
	}
}

func (m *MerMatchCollection) Wait() {
	m.wg.Wait()
}

func (m *MerMatchCollection) Done() {
	for _, v := range m.mcs {
		close(v.c)
	}
}

type MerMatch struct {
	// The mer length
	MerLen int
	Offset int

	c chan string

	wg *sync.WaitGroup
	m  *MatchedMers
}

func CreateMerMatch(len, offset int, m *MatchedMers, wg *sync.WaitGroup) MerMatch {
	c := make(chan string, 10)
	return MerMatch{len, offset, c, wg, m}
}

func (m *MerMatch) FindMers() {
	defer m.wg.Done()
	for line := range m.c {
		if len(line) < m.MerLen+m.Offset {
			log.Fatal("Out of range: ", m.Offset, m.MerLen, len(line))
		}
		mer := line[m.Offset : m.Offset+m.MerLen]
		m.m.AddMer(m.Offset, mer)
	}
}
