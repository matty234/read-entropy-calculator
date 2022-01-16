package main

import "log"

type MerMatch struct {
	// The mer length
	MerLen int
	Offset int

	m *MatchedMers
}

func CreateMerMatch(len, offset int, m *MatchedMers) MerMatch {
	return MerMatch{len, offset, m}
}

func (m MerMatch) FindMers(reads chan string) {
	for line := range reads {
		if len(line) < m.MerLen+m.Offset {
			log.Fatal("Out of range: ", m.Offset, m.MerLen)
		}
		mer := line[m.Offset : m.Offset+m.MerLen]
		m.m.AddMer(m.Offset, mer)
	}
}
