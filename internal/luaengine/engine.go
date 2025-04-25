package luaengine

import (
	"shadxel/internal/voxel"
)

const chunkSize = 32

type chunk struct {
	x0, y0, z0 int
	x1, y1, z1 int
}

type result struct {
	x    int
	z    int
	grid voxel.Grid
	err  error
}

type LuaEngine struct {
	workers []*Worker
}

func NewLuaEngine(script string) (*LuaEngine, error) {
	engine := LuaEngine{}
	for i := 0; i < 4; i++ {
		worker, err := NewWorker(script)
		if err != nil {
			return nil, err
		}
		engine.workers = append(engine.workers, worker)
	}
	return &engine, nil
}

func (le *LuaEngine) GenerateGridParallel(s, t int) (voxel.VoxelGrid, error) {
	size := s * 32
	half := size / 2

	for _, worker := range le.workers {
		worker.tasks = make([]chunk, 0, 32)
	}

	var count int
	for cz := -half; cz < half; cz += chunkSize {
		for cx := -half; cx < half; cx += chunkSize {
			worker := le.workers[count%len(le.workers)]
			worker.AddTask(chunk{
				x0: cx,
				y0: -half,
				z0: cz,
				x1: cx + chunkSize - 1,
				y1: half - 1,
				z1: cz + chunkSize - 1,
			})
			count++
		}
	}

	results := make(chan result, count)

	for _, worker := range le.workers {
		go func(worker *Worker) {
			worker.ProcessTasks(results, t)
		}(worker)
	}

	// Merge all results into a single voxel grid
	final := make(voxel.Grid, size)
	for y := range final {
		final[y] = make([][]voxel.Voxel, size)
		for x := range final[y] {
			final[y][x] = make([]voxel.Voxel, size)
		}
	}

	for res := range results {
		if res.err != nil {
			return voxel.VoxelGrid{}, res.err
		}
		for z := 0; z < chunkSize; z++ {
			for x := 0; x < chunkSize; x++ {
				for y := -half; y < half; y++ {
					final[half+y][half+res.x+x][half+res.z+z] = res.grid[half+y][x][z]
				}
			}
		}
		count--
		if count == 0 {
			close(results)
		}
	}

	return voxel.VoxelGrid{
		Data: final,
		Size: size,
	}, nil
}

func (le *LuaEngine) Close() {
	for _, w := range le.workers {
		w.close()
	}
}
