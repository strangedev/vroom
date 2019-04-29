package physics

import (
	"errors"
	"github.com/strangedev/vroom/algebra"
	"github.com/strangedev/vroom/gfx"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/imdraw"
)

const (
	splitThreshold = 2
	minBoxWidth    = 5
	maxItems = 1000
)

type QuadKey struct {
	Pnt    interface{}
	Bounds algebra.Rectangle
}

type QuadNode struct {
	bounds   algebra.Rectangle
	keys     []*QuadKey
	Children []*QuadNode
}

type QuadTree interface {
	insert(k *QuadKey)
	InsertAt(pnt interface{}, bounds algebra.Rectangle) *QuadKey
	Remove(k *QuadKey) error
	ClippingCandidates(k *QuadKey) <-chan *QuadKey
	ClippingCandidatesAt(bounds algebra.Rectangle) <-chan *QuadKey
	gfx.Drawable
}

func NewQuadTree(width, height float64) QuadTree {
	root := newQuadNode(
		algebra.Rectangle{
			Ul: algebra.Vector2{0, height},
			Ur: algebra.Vector2{width, height},
			Dl: algebra.Vector2{0, 0},
			Dr: algebra.Vector2{width, 0},
		},
	)
	return &root
}

func newQuadNode(bounds algebra.Rectangle) QuadNode {
	return QuadNode{
		bounds,
		make([]*QuadKey, 0, splitThreshold),
		make([]*QuadNode, 0, 4),
	}
}

func (n *QuadNode) isLeaf() bool {
	return len(n.Children) == 0
}

func (k *QuadKey) pickNode(ul, ur, dl, dr *QuadNode) (n *QuadNode, e error) {
	clipCount := 0
	if ul.bounds.Clips(k.Bounds) && clipCount <= 1 {
		n = ul
		clipCount++
	}
	if ur.bounds.Clips(k.Bounds) && clipCount <= 1 {
		n = ur
		clipCount++
	}
	if dl.bounds.Clips(k.Bounds) && clipCount <= 1 {
		n = dl
		clipCount++
	}
	if dr.bounds.Clips(k.Bounds) && clipCount <= 1 {
		n = dr
		clipCount++
	}
	if clipCount == 0 {
		panic("Key did not clip any quadrant! If this happens, beat up the programmer ... gosh that's me")
	} else if clipCount > 1 {
		e = errors.New("More than one candidate")
	}
	return
}

func (n *QuadNode) split() {
	width := 0.5 * (n.bounds.Ur[0] - n.bounds.Ul[0])
	height := 0.5 * (n.bounds.Ul[1] - n.bounds.Dl[1])
	boundsUl := algebra.Rectangle{
		Ul: n.bounds.Ul,
		Ur: algebra.Vector2{n.bounds.Ul[0] + width, n.bounds.Ul[1]},
		Dl: algebra.Vector2{n.bounds.Ul[0], n.bounds.Ul[1] - height},
		Dr: algebra.Vector2{n.bounds.Ul[0] + width, n.bounds.Ul[1] - height},
	}
	boundsUr := algebra.Rectangle{
		Ul: boundsUl.Ur,
		Ur: n.bounds.Ur,
		Dl: boundsUl.Dr,
		Dr: algebra.Vector2{n.bounds.Ur[0], n.bounds.Ur[1] - height},
	}
	boundsDl := algebra.Rectangle{
		Ul: boundsUl.Dl,
		Ur: boundsUl.Dr,
		Dl: n.bounds.Dl,
		Dr: algebra.Vector2{n.bounds.Dl[0] + width, n.bounds.Dl[1]},
	}
	boundsDr := algebra.Rectangle{
		Ul: boundsUl.Dr,
		Ur: boundsUr.Dr,
		Dl: boundsDl.Dr,
		Dr: n.bounds.Dr,
	}
	childUl := newQuadNode(boundsUl)
	childUr := newQuadNode(boundsUr)
	childDl := newQuadNode(boundsDl)
	childDr := newQuadNode(boundsDr)

	parentKeys := make([]*QuadKey, 0, len(n.keys))
	for _, key := range n.keys {
		node, err := key.pickNode(&childUl, &childUr, &childDl, &childDr)
		if err != nil {
			parentKeys = append(parentKeys, key)
		} else {
			node.keys = append(node.keys, key)
		}
	}

	if len(parentKeys) != len(n.keys) {
		n.keys = parentKeys
		n.Children = []*QuadNode{&childUl, &childUr, &childDl, &childDr}
	}
	// otherwise, no keys could be placed in children
	// do not modify QuadNode n and let the new nodes
	// be garbage collected
}

func (root *QuadNode) insert(k *QuadKey) {
	current := root
	for !current.isLeaf() {
		node, err := k.pickNode(
			current.Children[0],
			current.Children[1],
			current.Children[2],
			current.Children[3],
		)
		if err != nil {
			current.keys = append(current.keys, k)
			return
		}
		current = node
	}
	current.keys = append(current.keys, k)
	boxWidth, _ := current.bounds.Size()
	if boxWidth >= minBoxWidth && len(current.keys) > splitThreshold {
		current.split()
	}
}

func (root *QuadNode) InsertAt(pnt interface{}, bounds algebra.Rectangle) (k *QuadKey) {
	k = &QuadKey{pnt, bounds}
	root.insert(k)
	return
}

func (root *QuadNode) Remove(k *QuadKey) error {
	return nil
}

func (root *QuadNode) Draw(win *pixelgl.Window) {
	imd := imdraw.New(nil)
	imd.Color = pixel.RGB(1, 0, 0)

	queue := make(chan *QuadNode, 100)
	queue <- root
	for node := range queue{
		for _, child := range node.Children {
			queue <- child
		}
		if len(queue) == 0 {
			close(queue)
		}
		imd.Push(
			node.bounds.Ul.ToPixelVec(),
			node.bounds.Ur.ToPixelVec(),
		)
		imd.Line(2)
		imd.Push(
			node.bounds.Ul.ToPixelVec(),
			node.bounds.Dl.ToPixelVec(),
		)
		imd.Line(2)
		imd.Push(
			node.bounds.Ur.ToPixelVec(),
			node.bounds.Dr.ToPixelVec(),
		)
		imd.Line(2)
		imd.Push(
			node.bounds.Dl.ToPixelVec(),
			node.bounds.Dr.ToPixelVec(),
		)
		imd.Line(2)		
	}
	imd.Draw(win)
}

func (root *QuadNode) ClippingCandidates(k *QuadKey) <-chan *QuadKey {
	ch := root.ClippingCandidatesAt((*k).Bounds)
	candidates := make(chan *QuadKey, maxItems)
	go func () {
		for c := range ch {
			if c.Pnt != k.Pnt {
				candidates <- c
			}
		}
		close(candidates)
	}()
	return candidates
}

func (root *QuadNode) ClippingCandidatesAt(bounds algebra.Rectangle) <-chan *QuadKey {
	candidates := make(chan *QuadKey)
	go func() {
		queue := make(chan *QuadNode, maxItems)
		queue <- root
		for node := range queue{
			for _, key := range node.keys {
				if key.Bounds.Clips(bounds) {
					candidates <- key
				}
			}
			for _, child := range node.Children {
				if child.bounds.Clips(bounds) {
					queue <- child
				}
			}
			if len(queue) == 0 {
				close(queue)
			}
		}
		close(candidates)
	}()
	return candidates
}
