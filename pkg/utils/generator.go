package utils

import "sync"

type IDGenerator struct {
	nextID int
	mu     sync.Mutex
}

func NewIDGenerator(startID int) *IDGenerator {
	return &IDGenerator{nextID: startID}
}

func (g *IDGenerator) GetID() int {
	g.mu.Lock()
	defer g.mu.Unlock()

	id := g.nextID
	g.nextID++
	return id
}
