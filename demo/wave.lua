
function voxel(x, y, z, t)
    local cx, cy, cz = 25, 25, 25
    local dx = x - cx
    local dy = y - cy
    local dz = z - cz

    local dist = math.sqrt(dx*dx + dy*dy + dz*dz)
    local wave = math.sin(dist - t * 4)

    if wave > 0.2 then
        local r = 100 + math.floor(80 * math.sin(t + x * 0.2))
        local g = 100 + math.floor(80 * math.sin(t + y * 0.2))
        local b = 100 + math.floor(80 * math.sin(t + z * 0.2))
        return r, g, b
    end

    return 0, 0, 0
end
