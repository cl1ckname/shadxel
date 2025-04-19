function Draw(x, y, t)
	-- Background
	local r, g, b = 150, 200, 255

	-- Body
	if y >= 20 and y <= 35 and x >= 15 and x <= 35 then
		r, g, b = 80, 200, 120
	end

	-- Eyes
	if y >= 18 and y <= 20 and (x >= 18 and x <= 20 or x >= 30 and x <= 32) then
		r, g, b = 255, 255, 255
	end

	if y >= 18 and y <= 20 and (x >= 18 and x <= 20 or x >= 30 and x <= 32) and t % 2 == 1 then
		r, g, b = 30, 150, 60
	end
	-- Pupils
	if y == 19 and (x == 19 or x == 31) and t % 2 == 0 then
		r, g, b = 0, 0, 0
	end

	-- Mouth
	if y == 33 and x >= 22 and x <= 28 then
		r, g, b = 0, 0, 0
	end

	return r, g, b
end
