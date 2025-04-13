function voxel(x, y, z, t)
	x = x + math.floor(math.sin(t) * 3)
	y = y + math.floor(math.cos(t) * 3)
	local function isInMenger(x, y, z)
		while x > 0 or y > 0 or z > 0 do
			if x % 3 == 1 and y % 3 == 1 or
				x % 3 == 1 and z % 3 == 1 or
				y % 3 == 1 and z % 3 == 1 then
				return false
			end
			x = math.floor(x / 3)
			y = math.floor(y / 3)
			z = math.floor(z / 3)
		end
		return true
	end

	if isInMenger(x, y, z) then
		local r = 80 + x * 3
		local g = 80 + y * 3
		local b = 80 + z * 3
		return r % 255, g % 255, b % 255
	end

	return 0, 0, 0
end
