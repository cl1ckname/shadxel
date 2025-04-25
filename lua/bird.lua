-- @diagnostic disable: undefined-global
h = require("helpers")

local eye = h.color(64,32,64)
local nose = h.color(200, 160, 64)
local wing = h.color(170, 40, 50)
local body = h.color(220, 30, 80)
local hat = h.color(64, 32, 64)

local bw = 12
local bh = 10

local function iswing(x,y,z)
	return h.isin(x, -1, 1) and h.isin(z, -3, bw-3) and h.isin(y, -5, 1)
end

local function iswing2(x,y,z)
	return h.isin(x, -1, 1) and h.isin(z, -5, bw-5) and h.isin(y, -7, -1)
end

function Draw(x, y, z, t)
	local p = h.point(x,y,z)

	if h.isin(x, -bw, bw) and h.isin(y, 1, 8) and h.isin(z, 8, bw) then
		return eye
	end
	if h.isin(x, -bw, bw) and h.isin(y, 2, 8) and h.isin(z, 3, bw) then
		return eye
	end
	if h.isin(x, -3, 3) and h.isin(y, bh-2, bh+4) and h.isin(z, bw-5, bw+1) then
		return hat
	end
	if h.isin(x, -2, 2) and h.isin(y, bh-2, bh+5) and h.isin(z, bw-7, bw-1) then
		return hat
	end
	if h.isin(z, bw-1, bw+2) and h.isin(x, -2, 2) and h.isin(y, 0, 4) then
		return nose
	end
	if h.isin(x, bw-1,bw+1) and h.isin(z, 0, bw) and h.isin(y, -5, 1) then
		return wing
	end
	if iswing(x-bw, y,z) or iswing2(x-bw, y, z) or iswing(x+bw, y, z) or iswing2(x+bw, y, z) then
		return wing
	end
	if h.isinbox(p, h.point(-bw, -bh, -bw), h.point(bw,bh,bw)) then
		return body
	end
	if h.isin(y, -14, -bh+1) and (h.isin(x, 2, 8) or h.isin(x, -8, -2)) and (h.isin(z, 2, 8) or h.isin(z, -8, -2)) then
		return h.BLACK
	end
	if h.isin(x, -bw-1, bw+1) and y == 6 and z == 8 then
		return h.WHITE
	end
	if x < bw+1 and x > -bw-1 and y <= 6 and y >=4 and z <= 8 and z >=6 then
		return h.BLACK
	end
	return h.null
end
