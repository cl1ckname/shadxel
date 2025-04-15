# 🌟 Shadxel

**Shadxel** is a 3D voxel visualizer and scripting playground written in Go, powered by OpenGL and Lua.  
It allows you to generate animated, procedural voxel scenes using simple Lua scripts — rendered with lighting, shading, and camera controls.

---

## 🛠 Installation

### Dependencies

- [Go](https://golang.org/dl/) 1.20+
- [SDL2](https://www.libsdl.org/)
- OpenGL 3.3 or later
- Lua 5.1+

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
./shadxel name
```

Or start without arguments to open a file picker:

```bash
./shadxel
```

---

## 📜 Lua Scripting

Each script should define a global `voxel(x, y, z, t)` function.  
It gets voxel coordinates and time (in seconds), and returns RGB values.

Return `(0, 0, 0)` for empty space.

### Example: Glowing Sphere

```lua
function voxel(x, y, z, t)
    local d = x*x + y*y + z*z
    if d < 40 + math.sin(t * 2) * 10 then
        return 255, 180, 100
    end
    return 0, 0, 0
end
```

Place your `.lua` files in any directory and open them via CLI or the file picker.

---

## 🎮 Controls

| Action            | Control     |
|-------------------|-------------|
| Rotate view       | Mouse drag  |
| Zoom in/out       | Mouse wheel |
| Reload script     | Press `R`   |
| Open file dialog  | Press `O`   |
| Quit              | Press `ESC` |

---

## 📸 Roadmap

- [ ] In-app Lua editor with live feedback  
- [ ] Export to `.vox` or `.png`  
- [ ] Camera keyframe animations  
- [ ] Modular shader styles  
- [ ] Audio-reactive voxel input  

---

## 🧪 Made with 💖 in Go, OpenGL, and Lua
```
