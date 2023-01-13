package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/Harmos274/gotrans/warehouse"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

func runGraphic(initWr warehouse.Warehouse, cycles uint) {
	pixelgl.Run(generateMainLoop(initWr, cycles))
}

func generateMainLoop(initWr warehouse.Warehouse, cycles uint) func() {
	return func() {
		ch := make(chan warehouse.CycleState)

		go warehouse.CleanWarehouse(initWr, ch, cycles)

		var gr Graphical
		gr.CreateWindow()
		gr.CreateText("tour", 1, 0.5)
		for y := 1; y <= int(initWr.Height); y += 1 {
			for x := 1; x <= int(initWr.Length); x += 1 {
				gr.CreateRectangle(strconv.Itoa(y)+"/"+strconv.Itoa(x), x, y)
			}
		}

		currentCycle := 1
		for !gr.IsWindowClosed() {
			time.Sleep(time.Second)
			wr, ok := <-ch
			if ok {
				gr.ClearText("tour")
				placeEntities(wr.Warehouse, &gr)
				gr.ClearWindow()
				gr.ChangeText("tour", "tour "+strconv.Itoa(currentCycle))
				gr.DisplayRectangle("all")
				gr.DisplayText("all")
				gr.DisplayEntity("all")

				fmt.Println("tour", currentCycle)
				fmt.Println(ShowableWarehouse(wr))
				currentCycle += 1
			} else if currentCycle != 0 {
				fmt.Println("Terminé au tour", currentCycle-1)
				currentCycle = 0
			}
			gr.UpdateWindow()
		}
	}
}
func placeEntities(initWr warehouse.Warehouse, gr *Graphical) {
	gr.ClearEntities()
	for pos, trucks := range initWr.Trucks {
		gr.CreateEntity(trucks.Name, pos.X, int(initWr.Height)-pos.Y)
		gr.AddEntityInformation(trucks.Name, fmt.Sprintf("%d/%d\n", trucks.CurrentWeight, trucks.MaxWeight))
	}
	for pos, packages := range initWr.Packages {
		gr.CreateEntity(packages.Name, pos.X, int(initWr.Height)-pos.Y)
	}
	for pos, forklifts := range initWr.ForkLifts {
		gr.CreateEntity(forklifts.Name, pos.X, int(initWr.Height)-pos.Y)
	}
}

type GraphicalEntity struct {
	text        *text.Text
	information *text.Text
	rect        *imdraw.IMDraw
	color       pixel.RGBA
}

type Graphical struct {
	yRatio, xRatio float64
	win            *pixelgl.Window
	texts          map[string]*text.Text
	rects          map[string]*imdraw.IMDraw
	entities       map[string]GraphicalEntity
}

// ######################
// ####### WINDOW #######
// ######################
func (g *Graphical) CreateWindow() {
	cfg := pixelgl.WindowConfig{
		Title:  "GOTRANS",
		Bounds: pixel.R(0, 0, 1750, 700),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	g.win = win
	g.ClearWindow()
	// win.SetSmooth(true)
	g.texts = make(map[string]*text.Text)
	g.rects = make(map[string]*imdraw.IMDraw)
	g.entities = make(map[string]GraphicalEntity)
	g.xRatio = 250
	g.yRatio = 100
}

func (g *Graphical) ClearWindow() {
	g.win.Clear(pixel.RGB(0, 0.2, 0.2))
}

func (g *Graphical) UpdateWindow() {
	g.win.Update()
}

func (g *Graphical) IsWindowClosed() bool {
	return g.win.Closed()
}

// #####################
// ####### TEXTS #######
// #####################
func (g *Graphical) CreateText(id string, x float64, y float64) bool {
	_, exists := g.texts[id]
	if exists {
		return false
	}
	txt := text.New(pixel.V((x+0.5)*g.xRatio, y*g.yRatio), text.NewAtlas(basicfont.Face7x13, text.ASCII))
	txt.Color = pixel.RGB(1, 1, 1)
	_, _ = fmt.Fprintln(txt, id)
	g.texts[id] = txt
	return true
}

func (g *Graphical) ClearText(id string) bool {
	txt, exists := g.texts[id]
	if !exists {
		return false
	}
	txt.Clear()
	return true
}

func (g *Graphical) ChangeText(id string, newText string) bool {
	txt, exists := g.texts[id]
	if !exists {
		return false
	}
	_, _ = fmt.Fprintln(txt, newText)
	return true
}

func (g *Graphical) DisplayText(id string) bool {
	if id == "all" {
		for _, txt := range g.texts {
			txt.Draw(g.win, pixel.IM.Scaled(txt.Orig, 4))
		}
		return true
	}
	txt, exists := g.texts[id]
	if !exists {
		return false
	}
	txt.Draw(g.win, pixel.IM.Scaled(txt.Orig, 4))
	return true
}

// #########################
// ####### RECTANGLE #######
// #########################
func (g *Graphical) CreateRectangle(id string, x int, y int) bool {
	_, exists := g.rects[id]
	if exists {
		return false
	}
	real_x := float64(x) * g.xRatio
	real_y := float64(y) * g.yRatio
	rect := imdraw.New(nil)
	rect.Color = pixel.RGB(0.5, 0.5, 0.5)
	rect.Push(pixel.V(real_x, real_y))
	rect.Push(pixel.V(real_x+g.xRatio, real_y+g.yRatio))
	rect.Rectangle(3)
	rect.Color = pixel.RGB(0, 0, 0)
	rect.Push(pixel.V(real_x, real_y))
	rect.Push(pixel.V(real_x+g.xRatio, real_y+g.yRatio))
	rect.Rectangle(0)
	g.rects[id] = rect
	return true
}

func (g *Graphical) DisplayRectangle(id string) bool {
	if id == "all" {
		for _, rect := range g.rects {
			rect.Draw(g.win)
		}
		return true
	}
	rect, exists := g.rects[id]
	if !exists {
		return false
	}
	rect.Draw(g.win)
	return true
}

// ######################
// ####### ENTITY #######
// ######################
func (g *Graphical) CreateEntity(id string, x int, y int) bool {
	entity, exists := g.entities[id]
	var entityColor pixel.RGBA
	if exists {
		entityColor = entity.color
	} else {
		entityColor = pixel.RGB(rand.Float64(), rand.Float64(), rand.Float64())
	}
	txt := text.New(pixel.V((float64(x+1)+0.5)*g.xRatio, (float64(y)+0.5)*g.yRatio), text.NewAtlas(basicfont.Face7x13, text.ASCII))
	txt.Color = pixel.RGB(1, 1, 1)
	txt.Dot.X -= txt.BoundsOf(id).W() / 2
	_, _ = fmt.Fprintln(txt, id)

	real_x := float64(x+1) * g.xRatio
	real_y := float64(y) * g.yRatio
	rect := imdraw.New(nil)
	rect.Color = pixel.RGB(0.5, 0.5, 0.5)
	rect.Push(pixel.V(real_x, real_y))
	rect.Push(pixel.V(real_x+g.xRatio, real_y+g.yRatio))
	rect.Rectangle(3)
	rect.Color = entityColor
	rect.Push(pixel.V(real_x, real_y))
	rect.Push(pixel.V(real_x+g.xRatio, real_y+g.yRatio))
	rect.Rectangle(0)
	g.entities[id] = GraphicalEntity{text: txt, rect: rect, color: entityColor}
	return true
}

func (g *Graphical) ClearEntities() {
	for _, entity := range g.entities {
		entity.rect.Clear()
		entity.text.Clear()
		if entity.information != nil {
			entity.information.Clear()
		}
	}
}

func (g *Graphical) AddEntityInformation(id string, info string) bool {
	entity, exists := g.entities[id]
	if !exists {
		return false
	}
	entity.text.Dot.X -= entity.text.BoundsOf(info).W() / 2
	_, _ = fmt.Fprintln(entity.text, info)
	return true
}

func (g *Graphical) DisplayEntity(id string) bool {
	if id == "all" {
		for _, entity := range g.entities {
			entity.rect.Draw(g.win)
			entity.text.Draw(g.win, pixel.IM.Scaled(entity.text.Orig, 2))
			if entity.information != nil {
				entity.information.Draw(g.win, pixel.IM.Scaled(entity.text.Orig, 2))
			}
		}
		return true
	}
	entity, exists := g.entities[id]
	if !exists {
		return false
	}
	entity.rect.Draw(g.win)
	entity.text.Draw(g.win, pixel.IM.Scaled(entity.text.Orig, 2))
	return true
}
