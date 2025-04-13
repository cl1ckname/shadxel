function voxel(x, y, z, t)
	local d = x *x + y*y + z*z
	if d < (75 + math.sin(2 * math.pi / 10 * t)*10) then
		local sp = 100 * math.sin(2*math.pi / 100 * t)
		return 100 + sp, 150 - sp, 200
	end
    return 0, 0, 0
end
