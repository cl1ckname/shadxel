function voxel(x, y, z, t)
    -- Center coordinates
    local cx, cy, cz = 16, 8, 16
    local lx = x - cx
    local ly = y - cy
    local lz = z - cz

    -- Body: big yellow blob
    if lx*lx + (ly*1.5)^2 + lz*lz < 100 then
        return 255, 230, 0 -- Yellow
    end

    -- Ears (black cones on top)
    if ly > 12 and lx*lx + (ly - 14)^2 + lz*lz < 6 then
        return 20, 20, 20 -- Black tips
    end

    -- Cheeks (red spheres)
    if ly >= 6 and ly <= 8 and
       ((lx + 6)^2 + (lz + 4)^2 < 4 or (lx - 6)^2 + (lz + 4)^2 < 4) then
        return 255, 50, 50 -- Red cheeks
    end

    -- Eyes (black dots above cheeks)
    if ly == 9 and ((lx == -3 and lz == -2) or (lx == 3 and lz == -2)) then
        return 0, 0, 0 -- Eyes
    end

    -- Thunder tail (basic zigzag to the back)
    if (x >= 24 and x <= 26 and z >= 15 and z <= 17 and y >= 5 and y <= 7) or
       (x >= 26 and x <= 28 and z >= 13 and z <= 15 and y >= 7 and y <= 9) or
       (x >= 24 and x <= 26 and z >= 11 and z <= 13 and y >= 9 and y <= 11) then
        return 255, 230, 0 -- Same yellow as body
    end

    return 0, 0, 0 -- Transparent
end
