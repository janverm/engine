// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package renderer implements the scene renderer.
package renderer

import (
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
)

type Postprocessor struct {
	Width    int32
	Height   int32
	Fbo      uint32
	Tex      uint32
	Vao      uint32
	Prg      *gls.Program
	ppTime   gls.Uniform
	fTime    float32
	screen   []float32
	Renderer *Renderer
}

func (r *Renderer) CreatePostprocessor(width, height int32, vertexShaderSource, fragmentShaderSource string) *Postprocessor {
	pp := &Postprocessor{
		Width:    width,
		Height:   height,
		Renderer: r,
		screen: []float32{
			// xyz		color		texture coords
			-1, 1, 0, 1, 1, 1, 0, 1,
			-1, -1, 0, 1, 1, 1, 0, 0,
			1, -1, 0, 1, 1, 1, 1, 0,
			1, 1, 0, 1, 1, 1, 1, 1,
			-1, 1, 0, 1, 1, 1, 0, 1,
			1, -1, 0, 1, 1, 1, 1, 0,
		},
	}

	pp.Fbo = r.gs.GenFramebuffer()
	r.gs.BindFramebuffer(pp.Fbo)

	// set up a texture to render into
	pp.Tex = r.gs.GenTexture()
	r.gs.BindTexture(gls.TEXTURE_2D, pp.Tex)
	r.gs.TexImage2D(gls.TEXTURE_2D, 0, gls.RGB, width, height, gls.RGB, gls.UNSIGNED_BYTE, nil)
	//	r.gs.TexParameteri(gls.TEXTURE_2D, gls.TEXTURE_WRAP_S, gls.CLAMP_TO_EDGE)
	//	r.gs.TexParameteri(gls.TEXTURE_2D, gls.TEXTURE_WRAP_T, gls.CLAMP_TO_EDGE)
	r.gs.TexParameteri(gls.TEXTURE_2D, gls.TEXTURE_MIN_FILTER, gls.NEAREST)
	r.gs.TexParameteri(gls.TEXTURE_2D, gls.TEXTURE_MAG_FILTER, gls.NEAREST)
	r.gs.FramebufferTexture2D(gls.COLOR_ATTACHMENT0, gls.TEXTURE_2D, pp.Tex)
	r.gs.BindTexture(gls.TEXTURE_2D, 0)

	// attach depth and stencil buffers
	rbo := r.gs.GenRenderbuffer()
	r.gs.BindRenderbuffer(rbo)
	r.gs.RenderbufferStorage(gls.DEPTH24_STENCIL8, int(width), int(height))
	r.gs.FramebufferRenderbuffer(gls.DEPTH_STENCIL_ATTACHMENT, rbo)
	r.gs.BindRenderbuffer(0)

	// check the framebuffer status
	if r.gs.CheckFramebufferStatus() != gls.FRAMEBUFFER_COMPLETE {
		log.Fatal("Can't create frame buffer")
	}

	// create the "screen" quad
	vbo := r.gs.GenBuffer()
	r.gs.BindBuffer(gls.ARRAY_BUFFER, vbo)
	r.gs.BufferData(gls.ARRAY_BUFFER, 4*len(pp.screen), pp.screen, gls.STATIC_DRAW)

	pp.Vao = r.gs.GenVertexArray()
	r.gs.BindVertexArray(pp.Vao)
	r.gs.BindBuffer(gls.ARRAY_BUFFER, vbo)
	var offset uint32

	// position attribute
	r.gs.VertexAttribPointer(0, 3, gls.FLOAT, false, 8*4, offset)
	r.gs.EnableVertexAttribArray(0)
	offset += 3 * 4

	// color attribute
	r.gs.VertexAttribPointer(1, 3, gls.FLOAT, false, 8*4, offset)
	r.gs.EnableVertexAttribArray(1)
	offset += 3 * 4

	// texture coord attribute
	r.gs.VertexAttribPointer(2, 2, gls.FLOAT, false, 8*4, offset)
	r.gs.EnableVertexAttribArray(2)
	offset += 2 * 4

	// the screen shaders
	pp.Prg = r.gs.NewProgram()
	pp.Prg.AddShader(gls.VERTEX_SHADER, vertexShaderSource)
	pp.Prg.AddShader(gls.FRAGMENT_SHADER, fragmentShaderSource)
	pp.fTime = 0.0
	pp.ppTime.Init("ppTime")
	err := pp.Prg.Build()
	if err != nil {
		log.Fatal("can't create shader: %e", err)
	}

	r.gs.BindFramebuffer(0)

	return pp
}

func (pp *Postprocessor) Render(fbwidth, fbheight int32, scene core.INode, cam camera.ICamera, fTime float32) {
	// pass timing info to shader
	pp.fTime = fTime
	// render into the texture
	gs := pp.Renderer.gs
	gs.Viewport(0, 0, pp.Width, pp.Height)
	gs.BindFramebuffer(pp.Fbo)
	gs.ClearColor(0.1, 0.1, 0.1, 1.0)
	gs.Clear(gls.COLOR_BUFFER_BIT | gls.DEPTH_BUFFER_BIT)
	gs.Enable(gls.DEPTH_TEST)
	pp.Renderer.Render(scene, cam)

	// show texture on screen
	gs.Viewport(0, 0, int32(fbwidth), int32(fbheight))
	gs.BindFramebuffer(0)
	gs.ClearColor(1., 1., 1., 1.)
	gs.Clear(gls.COLOR_BUFFER_BIT)

	gs.UseProgram(pp.Prg)
	gs.BindVertexArray(pp.Vao)
	gs.Disable(gls.DEPTH_TEST)
	gs.BindTexture(gls.TEXTURE_2D, pp.Tex)
	gs.Uniform1fv(pp.ppTime.Location(gs), 1, &pp.fTime)
	gs.DrawArrays(gls.TRIANGLES, 0, 6)
}
