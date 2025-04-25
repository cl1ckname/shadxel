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

return {
	voxel = voxel,
	color = color,
	isin = isin,
	null = null,
	point = point,
	isinbox = isinbox,

	WHITE = color(255, 255, 255),
	BLACK = color(0, 0, 0),
}
