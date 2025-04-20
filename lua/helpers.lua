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

return {
	voxel = voxel,
	color = color,
	isin = isin,
	null = null,

	WHITE = color(255, 255, 255),
	BLACK = color(0, 0, 0),
}
