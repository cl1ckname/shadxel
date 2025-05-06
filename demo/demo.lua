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

function Draw(x, y, z, t)
	if x > -20 and x < 20 and y > -20 and y < 20 and z > -20 and z < 20 then
		return color((20 + x) * 5, (y + 20) * 5, (z + 20) * 5)
	end
	return null
end
