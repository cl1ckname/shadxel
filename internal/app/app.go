package app

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"shadxel/internal/camera"
	"shadxel/internal/config"
	"shadxel/internal/gridgen"
	"shadxel/internal/luaengine"
	"shadxel/internal/render"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	Version = "0.0.1"
	Period  = time.Second * 2
)

type App struct {
	window    *sdl.Window
	renderer  *render.Renderer
	camera    *camera.OrbitCamera
	mouseHeld bool
	lastX     int32
	engine    *gridgen.Gridgen
}

// Call this first — SDL requires main thread
func init() {
	runtime.LockOSThread()
}

func NewApp(c config.Config) (*App, error) {
	fmt.Printf("Running Shadxel %s!\n", Version)
	fmt.Printf("Size: %d (%d voxels on axis)\n", c.Size, c.Size*16)
	fmt.Printf("Script: %s\n", c.Script)
	fmt.Printf("Workers: %d\n", c.Workers)

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		return nil, err
	}

	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)

	window, err := sdl.CreateWindow("Shadxel", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, c.Width, c.Height, sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE)
	if err != nil {
		return nil, err
	}

	_, err = window.GLCreateContext()
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(c.Script)
	if err != nil {
		return nil, err
	}
	lua, err := luaengine.NewLuaEngine(string(content), c.Workers)
	if err != nil {
		return nil, err
	}
	gg := gridgen.New(lua, c.Size, Period)

	aspect := float32(c.Width) / float32(c.Height)
	renderer, err := render.NewRenderer(2./float32(c.Size)/16, aspect)
	if err != nil {
		return nil, err
	}

	return &App{
		window:   window,
		renderer: renderer,
		engine:   gg,
		camera:   camera.NewOrbitCamera(),
	}, nil
}

func (a *App) Run() {
	defer sdl.Quit()
	defer a.window.Destroy()

	running := true

	a.engine.Start()
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.MouseButtonEvent:
				if e.Button == sdl.BUTTON_LEFT {
					a.mouseHeld = e.State == sdl.PRESSED
					a.lastX = e.X
				}
			case *sdl.KeyboardEvent:
				if e.Type == sdl.KEYDOWN && e.Keysym.Sym == sdl.K_r {
					err := a.engine.Load()
					if err != nil {
						log.Println("Failed to reload script:", err)
					} else {
						log.Println("Script reloaded!")
					}
				}
				if e.Type == sdl.KEYDOWN && e.Keysym.Sym == sdl.K_F11 {
					flags := a.window.GetFlags()
					if flags&sdl.WINDOW_FULLSCREEN_DESKTOP != 0 {
						a.window.SetFullscreen(0)
					} else {
						a.window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
					}
				}
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_RESIZED {
					width := e.Data1
					height := e.Data2
					a.renderer.Resize(width, height)
				}
			case *sdl.MouseMotionEvent:
				if a.mouseHeld {
					dx := e.XRel
					dy := e.YRel
					a.camera.Rotate(float32(dx), float32(dy))
				}
			default:
				sdl.Delay(300)
			}
		}
		grid, err := a.engine.Get()
		if err != nil {
			log.Println("error", err)
			return
		}
		view := a.camera.ViewMatrix()
		a.renderer.Draw(grid, view)
		a.window.GLSwap()
	}
}
