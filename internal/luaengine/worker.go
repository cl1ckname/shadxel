package luaengine

import (
	"fmt"
	"shadxel/internal/voxel"

	lua "github.com/yuin/gopher-lua"
)

type Worker struct {
	L      *lua.LState
	script string
	fn     lua.LValue
	tasks  []chunk
}

func NewWorker(script string) (*Worker, error) {
	L := lua.NewState()

	// Preload helpers
	L.PreloadModule("helpers", func(L *lua.LState) int {
		if err := L.DoFile("lua/helpers.lua"); err != nil {
			panic(err)
		}
		return 1
	})

	return NewWorkerFromState(L, script)
}

func NewWorkerFromState(L *lua.LState, script string) (*Worker, error) {
	if err := L.DoString(script); err != nil {
		return nil, fmt.Errorf("worker failed to load lua script: %w", err)
	}

	if err := L.DoString(fallbackRegionWrap); err != nil {
		return nil, fmt.Errorf("failed to inject DrawRegion fallback: %w", err)
	}
	draw := L.GetGlobal("Draw")
	if draw.Type() != lua.LTFunction {
		return nil, fmt.Errorf("Draw(x,y,z,t) function not found in Lua")
	}

	fn := L.GetGlobal("DrawRegion")
	if fn.Type() != lua.LTFunction {
		return nil, fmt.Errorf("DrawRegion function not found in Lua")
	}

	return &Worker{
		L:      L,
		script: script,
		fn:     fn,
		tasks:  make([]chunk, 0, 32),
	}, nil
}

func (w *Worker) AddTask(task chunk) {
	w.tasks = append(w.tasks, task)
}

func (w *Worker) ProcessTasks(res chan<- result, t int) {
	for _, task := range w.tasks {
		grid, err := w.GenerateRegion(task.x0, task.y0, task.z0, task.x1, task.y1, task.z1, t)
		r := result{
			x:    task.x0,
			z:    task.z0,
			grid: grid,
			err:  err,
		}
		res <- r
	}
	w.tasks = nil
}

func (w *Worker) GenerateRegion(x0, y0, z0, x1, y1, z1, t int) (voxel.Grid, error) {
	L := w.L

	if err := L.CallByParam(lua.P{
		Fn:      w.fn,
		NRet:    1,
		Protect: true,
	}, lua.LNumber(x0), lua.LNumber(y0), lua.LNumber(z0),
		lua.LNumber(x1), lua.LNumber(y1), lua.LNumber(z1),
		lua.LNumber(t)); err != nil {
		return nil, fmt.Errorf("DrawRegion Lua error: %w", err)
	}

	result := L.Get(-1)
	L.Pop(1)

	table, ok := result.(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("expected table from DrawRegion")
	}

	sx := x1 - x0 + 1
	sy := y1 - y0 + 1
	sz := z1 - z0 + 1
	grid := make(voxel.Grid, sy)
	for y := 0; y < sy; y++ {
		grid[y] = make([][]voxel.Voxel, sx)
		for x := 0; x < sx; x++ {
			grid[y][x] = make([]voxel.Voxel, sz)
		}
	}
	// voxelCount := xSize * ySize * zSize
	// if table.Len() != voxelCount*4 {
	// 	return nil, fmt.Errorf("expected %d items, got %d", voxelCount*4, table.Len())
	// }

	i := 1
	for y := 0; y < sy; y++ {
		for x := 0; x < sx; x++ {
			for z := 0; z < sz; z++ {
				vox := uint32(lua.LVAsNumber(table.RawGetInt(i)))
				i++
				r := byte(vox >> 24)
				g := byte(vox >> 16)
				b := byte(vox >> 8)
				vis := vox&1 != 0

				grid[y][x][z] = voxel.NewVoxel(r, g, b, vis)
			}
		}
	}

	return grid, nil
}

func (w *Worker) close() {
	w.L.Close()
}

var fallbackRegionWrap = `
function DrawRegion(x0, y0, z0, x1, y1, z1, t)
	local out = {}
	for y = y0, y1 do
		for x = x0, x1 do
			for z = z0, z1 do
				local v = Draw(x, y, z, t)
				out[#out + 1] = v
			end
		end
	end
	return out
end
`
