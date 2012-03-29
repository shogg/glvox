package main

import (
	"github.com/banthar/Go-SDL/sdl"
	"github.com/shogg/gl"
	"github.com/shogg/glfps"
	"github.com/shogg/glvox"
	"io/ioutil"
	"fmt"
	"os"
	"time"
)

type vec3 [3]float32
type vec2 [2]float32

const (
	MaxTexBufferSize = 134217728	// 2^27
)

var (
	width, height = 640, 480
	prg0 = gl.Program(0)
	prg gl.Program
	lastModified time.Time
	shadowOff bool

	cam *glvox.Cam = glvox.NewCam()

	vertices = []vec3 {
		{-1.3,-1.0,-1.0}, // 0      2--3
		{ 1.3,-1.0,-1.0}, // 1      |\ |
		{-1.3, 1.0,-1.0}, // 2      | \|
		{ 1.3, 1.0,-1.0}, // 3      0--1
	}

	texCoords = []vec2 {
		{ 0.0, 0.0}, // 0      2--3
		{ 1.3, 0.0}, // 1      |\ |
		{ 0.0, 1.0}, // 2      | \|
		{ 1.3, 1.0}, // 3      0--1
	}

	vertexShaderSrc = `
		varying vec2 texCoord;
		uniform float time;

		void main()
		{
			texCoord = vec2(gl_MultiTexCoord0) * 2.0 - 1.0;
			gl_Position = ftransform();
		}`
)

func main() {
	sdl.Init(sdl.INIT_VIDEO)

	defer sdl.Quit()

	screen := sdl.SetVideoMode(width, height, 32, sdl.OPENGL)
	if screen == nil {
		panic("Couldn't set video mode: " + sdl.GetError() + "\n")
	}

	sdl.EnableKeyRepeat(200, 20)

	if err := gl.Init(); err != 0 {
		panic("gl error")
	}

	cam.Yaw(3.14)

	initGl()
	initShaders()
	initVoxels()
	mainLoop()
}

func initVoxels() {

	s := int32(164096)

	voxels := glvox.NewOctree()
	voxels.Dim(s, s, s)

	glvox.ReadBinvox("skull256.binvox", voxels, 896, 896, 896)

	data := voxels.Index
	buf := gl.GenBuffer()
	buf.Bind(gl.TEXTURE_BUFFER)
	gl.BufferData(gl.TEXTURE_BUFFER, len(data)*4, data, gl.STATIC_DRAW)

	gl.ActiveTexture(gl.TEXTURE0)
	tex := gl.GenTexture()
	tex.Bind(gl.TEXTURE_BUFFER)
	gl.TexBuffer(gl.TEXTURE_BUFFER, gl.R32I, buf)

	voxelsLoc := prg.GetUniformLocation("voxels")
	voxelsLoc.Uniform1i(0)

	var value [1]int32
	gl.GetIntegerv(gl.MAX_TEXTURE_BUFFER_SIZE, value[:])
	fmt.Println("max texture buffer size:", value[0]/1024/1024, "MiB")

	sizeLoc := prg.GetUniformLocation("size")
	sizeLoc.Uniform1i(int(voxels.WHD))
	fmt.Println("voxel data uploaded:", len(voxels.Index)*4/1024/1024, "MiB")
}

func initShaders() {

	// Compile vertex shader
	vshader := gl.CreateShader(gl.VERTEX_SHADER)
	vshader.Source(vertexShaderSrc)
	vshader.Compile()
	if vshader.Get(gl.COMPILE_STATUS) != gl.TRUE {
		panic("vertex shader compile error: " + vshader.GetInfoLog())
	}

	// Compile fragment shader
	fshader := gl.CreateShader(gl.FRAGMENT_SHADER)
	fragSrc, err := readShader("vox.glsl")
	if err != nil {
		panic(err)
	}

	fshader.Source(fragSrc)
	fshader.Compile()
	if fshader.Get(gl.COMPILE_STATUS) != gl.TRUE {
		panic("fragment shader compile error: " + fshader.GetInfoLog())
	}

	// Link program
	prg = gl.CreateProgram()
	prg.AttachShader(vshader)
	prg.AttachShader(fshader)
	prg.Link()
	if prg.Get(gl.LINK_STATUS) != gl.TRUE {
		panic("linker error: " + prg.GetInfoLog())
	}

	fmt.Println(prg.GetInfoLog())

	prg.Use()
}

func checkShaders() {

	finfo, err := os.Stat("vox.glsl")
	if err != nil { return }

	t := finfo.ModTime()

	if lastModified.IsZero() { lastModified = t }

	if t.After(lastModified) {
		lastModified = t
		prg.Delete()
		initShaders()
	}
}

func readShader(filename string) (s string, err error) {

	data, err := ioutil.ReadFile(filename)
	if err != nil { return }

	s = string(data)
	return
}

func initGl() {

	h := float64(height) / float64(width)

	gl.ShadeModel(gl.FLAT)

	gl.Viewport(0, 0, width, height)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Frustum(-1.0, 1.0, -h, h, 5.0, 60.0)

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Translatef(0.0, 0.0, -6.0)

	gl.EnableClientState(gl.VERTEX_ARRAY)
	gl.EnableClientState(gl.TEXTURE_COORD_ARRAY)
	gl.VertexPointer(3, gl.FLOAT, 0, vertices)
	gl.TexCoordPointer(2, gl.FLOAT, 0, texCoords)
}

func draw() {

	t := float32(sdl.GetTicks()) / 500.0
	time := prg.GetUniformLocation("time")
	time.Uniform1f(t)

	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
	if gl.GetError() != gl.NO_ERROR {
		panic("draw error")
	}
}

func drawOverlay() {
	prg0.Use()
	glfps.Draw(10, 10)
	prg.Use()
}

func mainLoop() {
	done := false
	for !done {
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch e.(type) {
			case *sdl.QuitEvent:
				done = true
				break
			case *sdl.KeyboardEvent:
				keyEvent := e.(*sdl.KeyboardEvent)
				if keyEvent.Type == sdl.KEYDOWN {
					done = handleKey(keyEvent)
				}
			}
		}

		updateShaderParams()
		draw()
		drawOverlay()
		sdl.GL_SwapBuffers()
		checkShaders()
	}
}

func updateShaderParams() {

	camPos := prg.GetUniformLocation("cam.pos")
	camPos.Uniform3f(cam.Pos.X, cam.Pos.Y, cam.Pos.Z)

	camDir := prg.GetUniformLocation("cam.dir")
	camDir.Uniform3f(cam.Dir.X, cam.Dir.Y, cam.Dir.Z)

	camUp := prg.GetUniformLocation("cam.up")
	camUp.Uniform3f(cam.Up.X, cam.Up.Y, cam.Up.Z)

	camLeft := prg.GetUniformLocation("cam.left")
	camLeft.Uniform3f(cam.Left.X, cam.Left.Y, cam.Left.Z)

	shadow := prg.GetUniformLocation("shadowOff")
	soff := 0; if shadowOff { soff = 1 }
	shadow.Uniform1i(soff)
}

func handleKey(e *sdl.KeyboardEvent) (done bool) {

	key := e.Keysym.Sym
	mod := e.Keysym.Mod

	if key == sdl.K_RETURN { done = true }

	var d float32 = .1
	if mod & sdl.KMOD_LSHIFT != 0 {
		d = 10.0
	}

	switch key {
	case sdl.K_s:
		shadowOff = !shadowOff
	case sdl.K_LEFT:
		if mod & sdl.KMOD_LCTRL != 0 {
			cam.Yaw(-.05)
		} else {
			cam.Strafe(d)
		}
	case sdl.K_RIGHT:
		if mod & sdl.KMOD_LCTRL != 0 {
			cam.Yaw( .05)
		} else {
			cam.Strafe(-1.0*d)
		}
	case sdl.K_UP:
		if mod & sdl.KMOD_LCTRL != 0 {
			cam.Pitch( .05)
		} else {
			cam.Move(d)
		}
	case sdl.K_DOWN:
		if mod & sdl.KMOD_LCTRL != 0 {
			cam.Pitch(-.05)
		} else {
			cam.Move(-1.0*d)
		}
	}

	return done
}
