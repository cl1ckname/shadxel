package luaengine_test

import (
	"shadxel/internal/luaengine"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func BenchmarkFuncall(b *testing.B) {
	w, err := luaengine.NewWorker(funcallScript)
	if err != nil {
		b.Fatalf("init worker: %v", err)
	}
	for i := 0; i < b.N; i++ {
		_, err := w.GenerateRegion(0, 0, 0, 32, 32, 32, 0)
		if err != nil {
			b.Fatalf("generate: %v", err)
		}
	}
}

func BenchmarkFunccallGO(b *testing.B) {
	L := lua.NewState()
	defer L.Close()
	L.PreloadModule("h", Loader)

	w, err := luaengine.NewWorkerFromState(L, funcallGoScript)
	if err != nil {
		b.Fatalf("init worker: %v", err)
	}

	for i := 0; i < b.N; i++ {
		_, err := w.GenerateRegion(0, 0, 0, 32, 32, 32, 0)
		if err != nil {
			b.Fatalf("generate: %v", err)
		}
	}
}

func Loader(L *lua.LState) int {
	// register functions to the table
	mod := L.SetFuncs(L.NewTable(), exports)
	// register other stuff
	L.SetField(mod, "name", lua.LString("value"))

	// returns the module
	L.Push(mod)
	return 1
}

var exports = map[string]lua.LGFunction{
	"fun": fun,
}

func fun(L *lua.LState) int {
	arg := L.Get(-1)
	L.Pop(1)
	i := arg.(lua.LNumber)
	v := luafun(float64(i))
	L.Push(lua.LNumber(v))
	return 1
}

func luafun(i float64) float64 {
	if i < 16 {
		return 0x00ffaa
	}
	if i > 16 {
		return 0x00ffbb
	}
	return 0x00bb55
}

const funcallScript = `
function fun(x) 
	if x < 16 then 
		return 0x00ffaa
	end
	if x > 16 then 
		return 0x00ffbb
	end
	return 0x00bb55
end

function Draw(x,y,z,t)
	return fun(x) + fun(y) + fun(z) + fun(2*x) + fun(2*y) * fun(2 * z)
end
`

const funcallGoScript = `
h = require("h")
	
function Draw(x,y,z,t)
	return h.fun(x) + h.fun(y) + h.fun(z) + h.fun(2*x) + h.fun(2 * y) + h.fun(2 * z)
end
`

func BenchmarkBird(b *testing.B) {
	w, err := luaengine.NewWorker(birdScript)
	if err != nil {
		b.Fatalf("init worker: %v", err)
	}
	for i := 0; i < b.N; i++ {
		_, err := w.GenerateRegion(0, 0, 0, 32, 32, 32, 0)
		if err != nil {
			b.Fatalf("generate: %v", err)
		}
	}
}

const birdScript = `
---@class Point
---@field x integer
---@field y integer
---@field z integer


function voxel(r, g, b, visible)
	local v = r * 0x1000000 + g * 0x10000 + b * 0x100
	if visible then
		v = v + 1
	end
	return v
end

---@param r integer
---@param g integer
---@param b integer
---@return integer
function color(r, g, b)
	return voxel(r, g, b, true)
end

local null = voxel(0, 0, 0, false)

function isin(x, a, b)
	return a < x and x < b
end

---@param x integer
---@param y integer
---@param z integer
---@return Point
function point(x, y, z)
	return { x = x, y = y, z = z }
end

---@param p Point
---@param p1 Point
---@param p2 Point
---@return boolean
function isinbox(p, p1, p2)
	return p1.x < p.x and p.x < p2.x and
		p1.y < p.y and p.y < p2.y and
		p1.z < p.z and p.z < p2.z
end


local eye = color(64,32,64)
local nose = color(200, 160, 64)
local wing = color(170, 40, 50)
local body = color(220, 30, 80)
local hat = color(64, 32, 64)

local bw = 12
local bh = 10

local function iswing(x,y,z)
	return isin(x, -1, 1) and isin(z, -3, bw-3) and isin(y, -5, 1)
end

local function iswing2(x,y,z)
	return isin(x, -1, 1) and isin(z, -5, bw-5) and isin(y, -7, -1)
end

function Draw(x, y, z, t)
	local p = point(x,y,z)

	if isin(x, -bw, bw) and isin(y, 1, 8) and isin(z, 8, bw) then
		return eye
	end
	if isin(x, -bw, bw) and isin(y, 2, 8) and isin(z, 3, bw) then
		return eye
	end
	if isin(x, -3, 3) and isin(y, bh-2, bh+4) and isin(z, bw-5, bw+1) then
		return hat
	end
	if isin(x, -2, 2) and isin(y, bh-2, bh+5) and isin(z, bw-7, bw-1) then
		return hat
	end
	if isin(z, bw-1, bw+2) and isin(x, -2, 2) and isin(y, 0, 4) then
		return nose
	end
	if isin(x, bw-1,bw+1) and isin(z, 0, bw) and isin(y, -5, 1) then
		return wing
	end
	if iswing(x-bw, y,z) or iswing2(x-bw, y, z) or iswing(x+bw, y, z) or iswing2(x+bw, y, z) then
		return wing
	end
	if isinbox(p, point(-bw, -bh, -bw), point(bw,bh,bw)) then
		return body
	end
	if isin(y, -14, -bh+1) and (isin(x, 2, 8) or isin(x, -8, -2)) and (isin(z, 2, 8) or isin(z, -8, -2)) then
		return BLACK
	end
	if isin(x, -bw-1, bw+1) and y == 6 and z == 8 then
		return WHITE
	end
	if x < bw+1 and x > -bw-1 and y <= 6 and y >=4 and z <= 8 and z >=6 then
		return BLACK
	end
	return null
end

`
