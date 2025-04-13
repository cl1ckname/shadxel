function voxel(x, y, z, t)
    if x > 10 and x < 20 and y > 10 and y < 20 and z > 10 and z < 20 then
        return 100 - ((x-10) * 10), (y-10)*10, (z-10) * 10
    end
    return 0, 0, 0
end
