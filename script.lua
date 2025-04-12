
-- Parameters
local width = 50
local height = 50
local box_size = 6

function voxel(x, y, t)
    -- Compute bouncing x position: goes 0 → 50 → 0 over time
	t = t * 3
    local speed = 1 -- cells per second
    local cycle = width - box_size
    local px = t % (2 * cycle)
    if px >= cycle then
        px = 2 * cycle - px
    end

    -- Draw a box
    if x >= px and x < px + box_size and y >= 22 and y < 22 + box_size then
        return 0, 255, 0 -- bright green box
    end

    return 20, 20, 40 -- dark background
end
