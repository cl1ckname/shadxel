function voxel(x, y, z, t)
	local size = 25
	x = x - size
	y = y - size
	z = z - size
	local min = -size
	local max = size

	-- Check if we're on the surface of the cube
	local on_x_edge = (x == min or x == max)
	local on_y_edge = (y == min or y == max)
	local on_z_edge = (z == min or z == max)

	-- Count how many coordinates are on the border
	local edge_count = 0
	if on_x_edge then edge_count = edge_count + 1 end
	if on_y_edge then edge_count = edge_count + 1 end
	if on_z_edge then edge_count = edge_count + 1 end

	-- A voxel is part of an edge if exactly 2 coordinates are at border
	if edge_count == 2 then
		return 200, 255, 255 -- bright cyan for edges
	end

	return 0, 0, 0
end
