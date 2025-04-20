package luaengine

import (
	"shadxel/internal/voxel"
)

type chunk struct {
	x0, y0, z0 int
	x1, y1, z1 int
}

type LuaEngine struct {
	script  string
	workers []*Worker
}

func NewLuaEngine(scriptPath string) (*LuaEngine, error) {
	engine := LuaEngine{
		script: scriptPath,
	}
	for i := 0; i < 1; i++ {
		worker, err := NewWorker(scriptPath)
		if err != nil {
			return nil, err
		}
		engine.workers = append(engine.workers, worker)
	}
	return &engine, nil
}

func (le *LuaEngine) GenerateGridParallel(size, t int) (voxel.VoxelGrid, error) {
	workerCount := len(le.workers)
	regionSize := size / workerCount
	half := size / 2

	tasks := make([]chunk, 0, workerCount)
	for i := 0; i < workerCount; i++ {
		x0 := -half
		y0 := -half + i*regionSize
		z0 := -half
		x1 := x0 + size - 1
		y1 := y0 + regionSize - 1
		z1 := z0 + size - 1
		tasks = append(tasks, chunk{x0, y0, z0, x1, y1, z1})
	}

	type result struct {
		index int
		grid  voxel.Grid
		err   error
	}

	results := make(chan result, workerCount)

	for i, task := range tasks {
		go func(i int, task chunk) {
			worker := le.workers[i]

			grid, err := worker.GenerateRegion(task.x0, task.y0, task.z0, task.x1, task.y1, task.z1, t)
			results <- result{i, grid, err}
		}(i, task)
	}

	// Merge all results into a single voxel grid
	finalGrid := make(voxel.Grid, size)
	for i := 0; i < workerCount; i++ {
		res := <-results
		if res.err != nil {
			return voxel.VoxelGrid{}, res.err
		}
		offset := res.index * regionSize
		copy(finalGrid[offset:], res.grid)
	}

	return voxel.VoxelGrid{
		Data: finalGrid,
		Size: size,
	}, nil
}

func (le *LuaEngine) Close() {
	for _, w := range le.workers {
		w.close()
	}
}
