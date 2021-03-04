package primitive

import (
	"fmt"
	"math"

	"github.com/fogleman/gg"
)

type Ellipse struct {
	Worker *Worker
	X, Y   int
	Rx, Ry int
	Circle bool
}

func closestSize(size int) int {
	c := possibleSizes[0]
	closestSoFar := diffInt(size, c)
	for _, p := range possibleSizes[1:] {
		distance := diffInt(size, p)
		if distance < closestSoFar {
			// Set the return
			c = p
			// Record closest distance
			closestSoFar = distance
		}
	}
	return c
}

var possibleSizes = [6]int{1, 2, 4, 6, 10, 13}

func getSizeIndex(size int) int {
	for i, s := range possibleSizes {
		if s == size {
			return i
		}
	}
	return 0
}

func NewRandomCircle(worker *Worker) *Ellipse {
	rnd := worker.Rnd
	x := rnd.Intn(worker.W)
	y := rnd.Intn(worker.H)
	r := possibleSizes[rnd.Intn(len(possibleSizes))*paintingToolScale]
	return &Ellipse{worker, x, y, r, r, true}
}

func (c *Ellipse) Draw(dc *gg.Context, scale float64) {
	dc.DrawEllipse(float64(c.X), float64(c.Y), float64(c.Rx), float64(c.Ry))
	dc.Fill()
}

func (c *Ellipse) SVG(attrs string) string {
	return fmt.Sprintf(
		"<ellipse %s cx=\"%d\" cy=\"%d\" rx=\"%d\" ry=\"%d\" />",
		attrs, c.X, c.Y, c.Rx, c.Ry)
}

/**
Is actually a circle >.> not an ellipse
*/
func (c *Ellipse) BORST(attrs string) string {
	return fmt.Sprintf(
		"%d,%d,%d,%s",
		c.X, c.Y, getSizeIndex(c.Rx), attrs)
}

func (c *Ellipse) Copy() Shape {
	a := *c
	return &a
}

func (c *Ellipse) Mutate() {
	w := c.Worker.W
	h := c.Worker.H
	rnd := c.Worker.Rnd
	switch rnd.Intn(3) {
	case 0:
		c.X = clampInt(c.X+int(rnd.NormFloat64()*16), 0, w-1)
		c.Y = clampInt(c.Y+int(rnd.NormFloat64()*16), 0, h-1)
	case 1:
	case 2:
		c.Rx = clampInt(closestSize(c.Rx+int(rnd.NormFloat64()*16)), 1, w-1)
		if c.Circle {
			c.Ry = c.Rx
		}
	}
}

func (c *Ellipse) Rasterize() []Scanline {
	w := c.Worker.W
	h := c.Worker.H
	lines := c.Worker.Lines[:0]
	aspect := float64(c.Rx) / float64(c.Ry)
	for dy := 0; dy < c.Ry; dy++ {
		y1 := c.Y - dy
		y2 := c.Y + dy
		if (y1 < 0 || y1 >= h) && (y2 < 0 || y2 >= h) {
			continue
		}
		s := int(math.Sqrt(float64(c.Ry*c.Ry-dy*dy)) * aspect)
		x1 := c.X - s
		x2 := c.X + s
		if x1 < 0 {
			x1 = 0
		}
		if x2 >= w {
			x2 = w - 1
		}
		if y1 >= 0 && y1 < h {
			lines = append(lines, Scanline{y1, x1, x2, 0xffff})
		}
		if y2 >= 0 && y2 < h && dy > 0 {
			lines = append(lines, Scanline{y2, x1, x2, 0xffff})
		}
	}
	return lines
}
