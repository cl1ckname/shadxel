h = require("helpers")

function Draw(x, y, z, t)
	x = math.abs(x)
	y = math.abs(y)
	z = math.abs(z)
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
		local r = 20 + x * 5
		local g = 30 + y * 2
		local b = 10 + z * 9
		return h.color(r % 255, g % 255, b % 255)
	end

	return h.null
end
