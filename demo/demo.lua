local function voxel(r, g, b, visible)
	return { r = r, g = g, b = b, visible = visible }
end

local function color(r, g, b)
	return voxel(r, g, b, true)
end

local null = { r = 0, g = 0, b = 0, visible = false }

function Draw(x, y, z, t)
	if x > -20 and x < 20 and y > -20 and y < 20 and z > -20 and z < 20 then
		return color((20 + x) * 5, (y + 20) * 5, (z + 20) * 5)
	end
	return null
end
