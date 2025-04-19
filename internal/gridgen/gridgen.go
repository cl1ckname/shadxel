package gridgen

import (
	"log"
	"shadxel/internal/luaengine"
	"shadxel/internal/voxel"
	"sync"
	"time"
)

type Gridgen struct {
	engine *luaengine.LuaEngine
	size   int
	period time.Duration

	mu      sync.Mutex
	grid    voxel.VoxelGrid
	running bool
	frame   int
	err     error
}

func New(engine *luaengine.LuaEngine, size int, period time.Duration) *Gridgen {
	return &Gridgen{
		engine: engine,
		period: period,
		size:   size,
	}
}

func (g *Gridgen) Start() error {
	g.running = true
	g.frame = 0
	g.err = nil
	g.updateGrid()
	if g.err != nil {
		return g.err
	}
	go g.gen()
	return nil
}

func (g *Gridgen) gen() {
	ticker := time.NewTicker(g.period)
	defer ticker.Stop()
	for g.running {
		g.genOrDelay(ticker)
	}
}

func (g *Gridgen) genOrDelay(t *time.Ticker) {
	select {
	case <-t.C:
		g.updateGrid()
		log.Println("Update!", g.frame)
	default:
		time.Sleep(time.Millisecond)
	}
}

func (g *Gridgen) updateGrid() {
	g.mu.Lock()
	defer g.mu.Unlock()

	grid, err := g.engine.GenerateGridParallel(g.size, g.frame)
	if err != nil {
		g.err = err
		g.running = false
		return
	}
	g.grid = grid
	g.frame++
}

func (g *Gridgen) Get() (voxel.VoxelGrid, error) {
	return g.grid, g.err
}

func (g *Gridgen) Load() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	// if err := g.engine.Load(); err != nil {
	// 	return err
	// }
	g.frame = 0
	g.err = nil
	return nil
}
