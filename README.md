# 🌟 Shadxel

**Shadxel** is a 3D voxel visualizer and scripting playground written in Go, powered by OpenGL and Lua.  
It allows you to generate animated, procedural voxel scenes using simple Lua scripts — rendered with lighting, shading, and camera controls.

---

## 🛠 Installation

### Dependencies

- [Go](https://golang.org/dl/) 1.20+
- [SDL2](https://www.libsdl.org/)
- OpenGL 3.3 or later

> Note: macOS and Windows support are possible with tweaks and dependencies.

### Build

```bash
git clone https://github.com/yourname/shadxel.git
cd shadxel
go build
```

---

## 🚀 Usage

Run with a Lua script:

```bash
./bin/shadxel lua/script.lua
```

Or start without arguments to open a file picker:

```bash
./shadxel
```

---

## 📜 Lua Scripting

Each script should define a global `Draw(x, y, z, t)` function.  
It gets voxel coordinates and time (in seconds), and returns RGB values encoded in i32.

Return `0xrrggbb01` for visible voxel, where `rr`, `gg` and `bb` are the hex codes of color channel. To return empty voxel just return 0. There are some helper functions in `lua/helpers.lua` file, for e.g. `color(r, g, b)` that calculates voxel or `null` that actually just an empty voxel.

### Example: Glowing Sphere

```lua
h = require("h")

function Draw(x, y, z, t)
    local d = x*x + y*y + z*z
    if d < 40 + math.sin(t * 2) * 10 then
        return h.color(255, 255, 100)
    end
    return h.null
end
```

Place your `.lua` files in `lua/` directory and open them via CLI.

---

## 📸 Roadmap

- [ ] In-app Lua editor with live feedback  
- [ ] Export to `.vox` or `.png`  
- [ ] Camera keyframe animations  
- [ ] Modular shader styles  

