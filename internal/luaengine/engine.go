package luaengine

import (
	"fmt"
	"os"
	"shadxel/internal/voxel"

	lua "github.com/yuin/gopher-lua"
)

type LuaEngine struct {
	L *lua.LState
}

func NewLuaEngine(scriptPath string) (*LuaEngine, error) {
	L := lua.NewState()
	if err := L.DoFile(scriptPath); err != nil {
		return nil, fmt.Errorf("failed to load lua script: %w", err)
	}
	return &LuaEngine{L: L}, nil
}

func (le *LuaEngine) GenerateGrid(width, height, t int) voxel.Grid {
	grid := make(voxel.Grid, height)
	for y := 0; y < height; y++ {
		grid[y] = make([]voxel.Color, width)
		for x := 0; x < width; x++ {
			r, g, b := le.callVoxelFunc(x, y, t)
			grid[y][x] = voxel.Color{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
			}
		}
	}
	return grid
}

func (le *LuaEngine) callVoxelFunc(x, y, t int) (int, int, int) {
	fn := le.L.GetGlobal("voxel")
	if fn.Type() != lua.LTFunction {
		return 255, 0, 255 // fallback: magenta if no voxel() found
	}

	if err := le.L.CallByParam(lua.P{
		Fn:      fn,
		NRet:    3,
		Protect: true,
	}, lua.LNumber(x), lua.LNumber(y), lua.LNumber(t)); err != nil {
		fmt.Fprintln(os.Stderr, "Lua error:", err)
		return 255, 0, 0
	}

	r := int(lua.LVAsNumber(le.L.Get(-3)))
	g := int(lua.LVAsNumber(le.L.Get(-2)))
	b := int(lua.LVAsNumber(le.L.Get(-1)))
	le.L.Pop(3)

	return r, g, b
}

func (le *LuaEngine) Close() {
	le.L.Close()
}
