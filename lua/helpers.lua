---@class Voxel
---@field r integer
---@field g integer
---@field b integer
---@field visible boolean


function voxel(r, g, b, visible)
	return { r = r, g = g, b = b, visible = visible }
end

---@param r integer
---@param g integer
---@param b integer
---@return Voxel
function color(r, g, b)
	return voxel(r, g, b, true)
end

null = { r = 0, g = 0, b = 0, visible = false }

return {
	voxel = voxel,
	color = color,
	null = null,
}
