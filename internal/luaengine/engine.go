package luaengine

import (
	"fmt"
	"os"
	"shadxel/internal/voxel"

	lua "github.com/yuin/gopher-lua"
)

type LuaEngine struct {
	L      *lua.LState
	script string
	fn     lua.LValue
}

func NewLuaEngine(scriptPath string) (*LuaEngine, error) {
	L := lua.NewState()
	engine := LuaEngine{
		script: scriptPath,
		L:      L,
	}
	if err := engine.Load(); err != nil {
		return nil, err
	}
	return &engine, nil
}

func (le *LuaEngine) Load() error {
	le.L.PreloadModule("helpers", func(L *lua.LState) int {
		if err := L.DoFile("lua/helpers.lua"); err != nil {
			panic(err)
		}
		return 1
	})
	if err := le.L.DoFile(le.script); err != nil {
		return fmt.Errorf("failed to load lua script: %w", err)
	}

	if err := le.L.DoString(fallbackRegionWrap); err != nil {
		return fmt.Errorf("failed to inject DrawRegion fallback: %w", err)
	}
	fn := le.L.GetGlobal("DrawRegion")
	if fn.Type() != lua.LTFunction {
		return fmt.Errorf("wrap has wrong type")
	}

	le.fn = fn
	return nil
}

func (le *LuaEngine) GenerateGrid(size, t int) (voxel.VoxelGrid, error) {
	grid, err := le.generateGrid(size, t)
	vg := voxel.VoxelGrid{Data: grid, Size: size}
	return vg, err
}

func (le *LuaEngine) generateGrid(size, t int) (voxel.Grid, error) {
	half := size / 2
	x0, y0, z0 := -half, -half, -half
	x1, y1, z1 := half-1, half-1, half-1

	lua0X := lua.LNumber(x0)
	lua0Y := lua.LNumber(y0)
	lua0Z := lua.LNumber(z0)
	luaX := lua.LNumber(x1)
	luaY := lua.LNumber(y1)
	luaZ := lua.LNumber(z1)
	luaT := lua.LNumber(t)
	if err := le.L.CallByParam(lua.P{
		Fn:      le.fn,
		NRet:    1,
		Protect: true,
	}, lua0X, lua0Y, lua0Z, luaX, luaY, luaZ, luaT); err != nil {
		fmt.Fprintln(os.Stderr, "Lua error:", err)
		return nil, err
	}

	tbl, ok := le.L.Get(-1).(*lua.LTable)
	le.L.Pop(1)
	if !ok {
		return nil, fmt.Errorf("returns not a table")
	}

	grid := make(voxel.Grid, size)
	for yi := 1; yi <= size; yi++ {
		rowY := tbl.RawGetInt(yi)
		grid[y0+yi-1+half] = make([][]voxel.Voxel, size)

		rowYTable, ok := rowY.(*lua.LTable)
		if !ok {
			continue
		}

		for xi := 1; xi <= size; xi++ {
			rowX := rowYTable.RawGetInt(xi)
			grid[y0+yi-1+half][x0+xi-1+half] = make([]voxel.Voxel, size)

			rowXTable, ok := rowX.(*lua.LTable)
			if !ok {
				continue
			}

			for zi := 1; zi <= size; zi++ {
				vobj := rowXTable.RawGetInt(zi)
				vtable, ok := vobj.(*lua.LTable)
				if !ok {
					continue
				}

				v := voxel.Voxel{
					Color: voxel.Color{
						R: uint8(lua.LVAsNumber(vtable.RawGetString("r"))),
						G: uint8(lua.LVAsNumber(vtable.RawGetString("g"))),
						B: uint8(lua.LVAsNumber(vtable.RawGetString("b"))),
					},
					V: lua.LVAsBool(vtable.RawGetString("visible")),
				}
				grid[y0+yi-1+half][x0+xi-1+half][z0+zi-1+half] = v
			}
		}
	}

	return grid, nil
}

func (le *LuaEngine) Close() {
	le.L.Close()
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
