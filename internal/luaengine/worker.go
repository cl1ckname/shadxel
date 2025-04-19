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
}

func NewWorker(scriptPath string) (*Worker, error) {
	L := lua.NewState()

	// Preload helpers
	L.PreloadModule("helpers", func(L *lua.LState) int {
		if err := L.DoFile("lua/helpers.lua"); err != nil {
			panic(err)
		}
		return 1
	})

	if err := L.DoFile(scriptPath); err != nil {
		return nil, fmt.Errorf("worker failed to load lua script: %w", err)
	}

	if err := L.DoString(fallbackRegionWrap); err != nil {
		return nil, fmt.Errorf("failed to inject DrawRegion fallback: %w", err)
	}

	fn := L.GetGlobal("DrawRegion")
	if fn.Type() != lua.LTFunction {
		return nil, fmt.Errorf("DrawRegion function not found in Lua")
	}

	return &Worker{
		L:      L,
		script: scriptPath,
		fn:     fn,
	}, nil
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
				r := uint8(lua.LVAsNumber(table.RawGetInt(i)))
				i++
				g := uint8(lua.LVAsNumber(table.RawGetInt(i)))
				i++
				b := uint8(lua.LVAsNumber(table.RawGetInt(i)))
				i++
				vis := lua.LVAsNumber(table.RawGetInt(i)) != 0
				i++

				grid[y][x][z] = voxel.Voxel{
					Color: voxel.Color{R: r, G: g, B: b},
					V:     vis,
				}
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
				out[#out + 1] = v.r
				out[#out + 1] = v.g
				out[#out + 1] = v.b
				out[#out + 1] = v.visible and 1 or 0
			end
		end
	end
	return out
end
`
