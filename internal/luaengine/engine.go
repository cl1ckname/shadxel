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
}

func NewLuaEngine(scriptPath string) (*LuaEngine, error) {
	L := lua.NewState()
	if err := L.DoFile(scriptPath); err != nil {
		return nil, fmt.Errorf("failed to load lua script: %w", err)
	}
	engine := LuaEngine{
		script: scriptPath,
	}
	if err := engine.Load(); err != nil {
		return nil, err
	}
	return &engine, nil
}

func (le *LuaEngine) Load() error {
	L := lua.NewState()
	if err := L.DoFile(le.script); err != nil {
		return fmt.Errorf("failed to load lua script: %w", err)
	}
	le.L = L
	return nil
}

func (le *LuaEngine) GenerateGrid(size, t int) voxel.VoxelGrid {
	grid := make(voxel.Grid, size)
	for y := 0; y < size; y++ {
		grid[y] = make([][]voxel.Color, size)
		for x := 0; x < size; x++ {
			grid[y][x] = make([]voxel.Color, size)
			for z := 0; z < size; z++ {
				r, g, b := le.callVoxelFunc(x, y, z, t, size)
				grid[y][x][z] = voxel.Color{
					R: uint8(r),
					G: uint8(g),
					B: uint8(b),
				}
			}
		}
	}
	return voxel.VoxelGrid{
		Data: grid,
		Size: size,
	}
}

func (le *LuaEngine) callVoxelFunc(x, y, z, t, size int) (int, int, int) {
	fn := le.L.GetGlobal("voxel")
	if fn.Type() != lua.LTFunction {
		return 255, 0, 255 // fallback: magenta if no voxel() found
	}

	luaX := lua.LNumber(x - size/2)
	luaY := lua.LNumber(y - size/2)
	luaZ := lua.LNumber(z - size/2)
	if err := le.L.CallByParam(lua.P{
		Fn:      fn,
		NRet:    3,
		Protect: true,
	}, luaX, luaY, luaZ, lua.LNumber(t)); err != nil {
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
