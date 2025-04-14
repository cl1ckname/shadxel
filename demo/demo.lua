function voxel(x, y, z, t)
	if x > -20 and x < 20 and y > -20 and y < 20 and z > -20 and z < 20 then
		return (20 + x) * 5, (y + 20) * 5, (z + 20) * 5
	end
	return 0, 0, 0
end
