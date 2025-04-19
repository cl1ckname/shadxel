-- @diagnostic disable: undefined-global
h = require("helpers")

function Draw(x, y, z, t)
	r = 20 + math.sin(t) * 5
	d = x * x + y * y + z * z
	if d <= r * r then
		return h.color(255, 255, 0)
	end
	return h.null
end
