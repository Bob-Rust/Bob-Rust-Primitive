package primitive

import (
	"fmt"
	"github.com/fogleman/gg"
)

type Rectangle struct {
	Worker *Worker
	X1, Y1 int
	X2, Y2 int
}

func NewRandomRectangle(worker *Worker) *Rectangle {
	rnd := worker.Rnd
	x1 := rnd.Intn(worker.W)
	y1 := rnd.Intn(worker.H)
	var size = possibleSizes[rnd.Intn(6)*paintingToolScale]
	x2 := clampInt(x1+size+1, 0, worker.W-1)
	y2 := clampInt(y1+size+1, 0, worker.H-1)
	return &Rectangle{worker, x1, y1, x2, y2}
}

func (r *Rectangle) bounds() (x1, y1, x2, y2 int) {
	x1, y1 = r.X1, r.Y1
	x2, y2 = r.X2, r.Y2
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	return
}

func (r *Rectangle) Draw(dc *gg.Context, scale float64) {
	x1, y1, x2, y2 := r.bounds()
	dc.DrawRectangle(float64(x1), float64(y1), float64(x2-x1+1), float64(y2-y1+1))
	dc.Fill()
}

func (r *Rectangle) BORST(attrs string) string {
	panic("implement me")
}

func (r *Rectangle) SVG(attrs string) string {
	x1, y1, x2, y2 := r.bounds()
	w := x2 - x1 + 1
	h := y2 - y1 + 1
	return fmt.Sprintf(
		"<rect %s x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" />",
		attrs, x1, y1, w, h)
}

func (r *Rectangle) Copy() Shape {
	a := *r
	return &a
}

func (r *Rectangle) Mutate() {
	w := r.Worker.W
	h := r.Worker.H
	rnd := r.Worker.Rnd
	var size = possibleSizes[rnd.Intn(6)*paintingToolScale]
	switch rnd.Intn(2) {
	case 0:
		var offsetX = int(rnd.NormFloat64() * 16)
		var offsetY = int(rnd.NormFloat64() * 16)
		r.X1 = clampInt(r.X1+offsetX, 0, w-1)
		r.Y1 = clampInt(r.Y1+offsetY, 0, h-1)
		r.X2 = clampInt(r.X2+offsetX+size, 0, w-1)
		r.Y2 = clampInt(r.Y2+offsetY+size, 0, h-1)
	case 1:
		r.X2 = clampInt(r.X1+size, 0, w-1)
		r.Y2 = clampInt(r.Y1+size, 0, h-1)
	}
}

func (r *Rectangle) Rasterize() []Scanline {
	x1, y1, x2, y2 := r.bounds()
	lines := r.Worker.Lines[:0]
	for y := y1; y <= y2; y++ {
		lines = append(lines, Scanline{y, x1, x2, 0xffff})
	}
	return lines
}
