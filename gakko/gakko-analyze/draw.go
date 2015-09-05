package main

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/ajstarks/svgo"
	"github.com/ledyba/gakko-analyzer/nico/client"
)

func drawGroup(s *svg.SVG, width int, bottom int, color string, grp *Group) {
	offX := grp.pt.X
	offY := grp.pt.Y - bottom
	s.Rect(offX-20, offY+grp.minY-20, grp.maxX+40, grp.maxY-grp.minY+40, fmt.Sprintf("%s;", color))
	for i, log := range grp.logs {
		x := grp.layout[i].X
		y := grp.layout[i].Y
		s.Text(x+offX-10, y+offY-20, fmt.Sprintf("No. %d", log.No))
		s.Text(x+offX, y+offY, log.Content)
	}
	if grp.parent != nil {
		drawGroup(s, width, bottom, color, grp.parent)
		s.Line(grp.pt.X-20, offY+20, grp.parent.pt.X+grp.parent.maxX+20, grp.parent.pt.Y-bottom+20, "stroke: black;stroke-width:2")
	}
}

func drawGraph(w io.Writer, logs []*client.Chat, grps []*Group) {
	for _, grp := range grps {
		grp.DoLayout(logs[0].Date)
	}
	top, bottom := Layout(grps)
	width := int((logs[len(logs)-1].Date-logs[0].Date)/36) + 100
	height := top - bottom + 300
	s := svg.New(w)
	s.Start(width, height)
	{ //time
		it := time.Unix(logs[0].Date, 0)
		it = it.Truncate(time.Hour * 24)
		it = it.Add(time.Hour * 24)
		end := time.Unix(logs[len(logs)-1].Date, 0)
		end = end.Truncate(time.Hour * 24)
		for it.Sub(end) < time.Hour*12 {
			x := int((it.Unix() - logs[0].Date) / 36)
			y := 30
			s.Text(x, y, it.Format("2006/01/02"), "font-color: rgb(0,0,0); opacity: 0.5; font-size: 2em;")
			s.Line(x, 0, x, height, "stroke: rgb(0,0,0); opacity: 0.5;stroke-width:2; font-size: 2em;stroke-dasharray:  5,5;")
			it = it.Add(24 * time.Hour)
		}
	}
	for _, grp := range grps {
		color := s.RGBA(192+rand.Intn(64), 192+rand.Intn(64), 192+rand.Intn(64), 0.5)
		drawGroup(s, width, bottom-200, color, grp)
	}
	s.End()

}
