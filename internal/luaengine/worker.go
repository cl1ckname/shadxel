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

	sizeX := x1 - x0 + 1
	sizeY := y1 - y0 + 1
	sizeZ := z1 - z0 + 1

	grid := make(voxel.Grid, sizeY)
	for yi := 1; yi <= sizeY; yi++ {
		rowY := table.RawGetInt(yi)
		rowYTable, ok := rowY.(*lua.LTable)
		if !ok {
			continue
		}
		grid[yi-1] = make([][]voxel.Voxel, sizeX)
		for xi := 1; xi <= sizeX; xi++ {
			rowX := rowYTable.RawGetInt(xi)
			rowXTable, ok := rowX.(*lua.LTable)
			if !ok {
				continue
			}
			grid[yi-1][xi-1] = make([]voxel.Voxel, sizeZ)
			for zi := 1; zi <= sizeZ; zi++ {
				vobj := rowXTable.RawGetInt(zi)
				vtable, ok := vobj.(*lua.LTable)
				if !ok {
					continue
				}
				grid[yi-1][xi-1][zi-1] = voxel.Voxel{
					Color: voxel.Color{
						R: uint8(lua.LVAsNumber(vtable.RawGetString("r"))),
						G: uint8(lua.LVAsNumber(vtable.RawGetString("g"))),
						B: uint8(lua.LVAsNumber(vtable.RawGetString("b"))),
					},
					V: lua.LVAsBool(vtable.RawGetString("visible")),
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
		local rowY = {}
		for x = x0, x1 do
			local rowX = {}
			for z = z0, z1 do
				rowX[#rowX + 1] = Draw(x, y, z, t)
			end
			rowY[#rowY + 1] = rowX
		end
		out[#out + 1] = rowY
	end
	return out
end
`
