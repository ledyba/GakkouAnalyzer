package main

import (
	lg "log"
	"math"
	"sort"
	"strings"

	"github.com/antzucaro/matchr"
	"github.com/ledyba/gakko-analyzer/nico/client"
)

const limit = 0.5

type point struct {
	X int
	Y int
}

type Group struct {
	pt        point
	totalMinY int
	totalMaxY int
	parent    *Group
	logs      []*client.Chat
	layout    []*point
	maxX      int
	minY      int
	maxY      int

	fromMainDist float32
}

func NewGroup(p *Group) *Group {
	return &Group{
		parent: p,
		logs:   []*client.Chat{},
	}
}

func Layout(grps []*Group) (int, int) {
	mx := 0
	height := 0
	for _, grp := range grps {
		l := grp.totalMaxY - grp.totalMinY
		height += l
	}
	var mg *Group
	for _, grp := range grps {
		l := grp.totalMaxY - grp.totalMinY
		if l > mx {
			mx = l
			mg = grp
		}
	}
	mg.LayoutY(0)
	bottom := mg.totalMinY
	top := mg.totalMaxY
	var others []*Group
	for _, grp := range grps {
		if grp != mg {
			grp.fromMainDist = grp.AvgDistance(mg)
			others = append(others, grp)
		}
	}
	sort.Sort(distgrp(others))
	for _, grp := range others {
		if -bottom < top {
			grp.LayoutY(bottom - grp.totalMaxY)
			bottom = bottom - (grp.totalMaxY - grp.totalMinY)
		} else {
			grp.LayoutY(top - grp.totalMinY)
			top = top + (grp.totalMaxY - grp.totalMinY)
		}
	}
	return top, bottom
}

type distgrp []*Group

func (p distgrp) Len() int {
	return len(p)
}

func (p distgrp) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p distgrp) Less(i, j int) bool {
	return p[i].fromMainDist < p[j].fromMainDist
}

func (grp *Group) LayoutY(y int) {
	if grp.parent != nil {
		grp.parent.LayoutY(y)
	}
	grp.pt.Y = y
}

func (grp *Group) DoLayout(beg int64) {
	if grp.parent != nil {
		grp.parent.DoLayout(beg)
	}
	begDate := grp.logs[0].Date
	grp.pt.X = int((begDate - beg) / 36)

	allMaxX := 0
	allMaxY := 0
	allMinY := 0
	for _, v := range grp.logs {
		date := v.Date - begDate
		pt := &point{}
		pt.X = int(date / 36)
		allMaxX = pt.X
		minY := 0
		maxY := 0
		for _, lpt := range grp.layout {
			if lpt.X+30 >= pt.X {
				if lpt.Y < minY {
					minY = lpt.Y
				}
				if lpt.Y > maxY {
					maxY = lpt.Y
				}
			}
		}
		if -minY < maxY {
			pt.Y = minY - 20
			if pt.Y < allMinY {
				allMinY = pt.Y
			}
		} else {
			pt.Y = maxY + 20
			if pt.Y > allMaxY {
				allMaxY = pt.Y
			}
		}
		grp.layout = append(grp.layout, pt)
	}
	grp.maxY = allMaxY
	grp.minY = allMinY
	grp.maxX = allMaxX + 300
	if grp.parent != nil {
		if grp.parent.totalMaxY < grp.maxY {
			grp.totalMaxY = grp.maxY
		} else {
			grp.totalMaxY = grp.parent.totalMaxY
		}
		if grp.parent.totalMinY < grp.minY {
			grp.totalMinY = grp.parent.totalMinY
		} else {
			grp.totalMinY = grp.minY
		}
	} else {
		grp.totalMinY = grp.minY
		grp.totalMaxY = grp.maxY
	}
}

func (grp *Group) Add(log *client.Chat) {
	grp.logs = append(grp.logs, log)
}
func (grp *Group) Distance(log *client.Chat) (float32, float32) {
	nsum := float32(0)
	ncnt := 0
	fsum := float32(0)
	fcnt := 0
	for _, v := range grp.logs {
		if v.No >= log.No-1000 {
			nsum += dist(v, log)
			ncnt++
		} else {
			fsum += dist(v, log)
			fcnt++
		}
	}
	return (nsum / float32(ncnt)), (fsum / float32(fcnt))
}
func (grp *Group) AvgDistance(g2 *Group) float32 {
	sum := float32(0)
	cnt := 0
	for _, a := range grp.logs {
		for _, b := range g2.logs {
			cnt++
			sum += dist(a, b)
		}
	}
	return (sum / float32(cnt))
}

func dist(a *client.Chat, b *client.Chat) float32 {
	al := len([]rune(a.Content))
	bl := len([]rune(b.Content))
	return float32(matchr.Levenshtein(a.Content, b.Content)) / float32(math.Max(float64(al), float64(bl)))
}

func matchLog(log *client.Chat, words []string) bool {
	for _, v := range words {
		if strings.Contains(log.Content, v) {
			return true
		}
	}
	return false
}

func makeGraph(logs []*client.Chat, words []string) []*Group {
	var groups []*Group
	for _, log := range logs {
		neargrp := -1
		nearmin := float32(0.5)
		fargrp := -1
		farmin := float32(0.5)
		if matchLog(log, words) {
			for i, grp := range groups {
				near, far := grp.Distance(log)
				if near < nearmin {
					nearmin = near
					neargrp = i
				}
				if far < farmin {
					farmin = far
					fargrp = i
				}
			}
			if neargrp >= 0 {
				groups[neargrp].Add(log)
				lg.Printf("%d group (%s) <-- %s", neargrp, groups[neargrp].logs[0].Content, log.Content)
			} else if fargrp >= 0 {
				grp := NewGroup(groups[fargrp])
				groups[fargrp] = grp
				grp.Add(log)
				lg.Printf("%d group (%s) <== %s", fargrp, groups[fargrp].logs[0].Content, log.Content)
			} else {
				grp := NewGroup(nil)
				groups = append(groups, grp)
				grp.Add(log)
				lg.Printf("new: %s", log.Content)
			}
		}
	}
	return groups
}
