package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/cryptix/canvas"
)

type World struct {
	Canvas   *canvas.Canvas
	Cities   []*Node
	n, peers int
}

func NewWorld(n, peers int, cv *canvas.Canvas) *World {
	log.Println("Building new city.")
	world := new(World)

	wsize := cv.Bounds().Size()

	world.Canvas = cv
	cv.DrawRect(
		color.RGBA{0, 0, 0, 255},
		canvas.Vector{0, 0},
		canvas.Vector{float64(wsize.X), float64(wsize.Y)})
	rand.Seed(time.Now().UTC().UnixNano())

	world.Cities = make([]*Node, n)

	// log.Print("Init nodes...")
	for i := 0; i < n; i += 1 {
		world.Cities[i] = NewNode(peers, world.Canvas)
	}
	// log.Print("Done.")

	nodesCopy := make([]*Node, n)
	copy(nodesCopy, world.Cities)
	// log.Print("Sorting cities...")
	for _, node := range world.Cities {
		// sort the nodes by distance
		sorter := NodeSorter{nodesCopy, node}
		sort.Sort(sorter)
		node.Peers = append(node.Peers, nodesCopy[1:peers+1]...)
	}
	// log.Print("Done.")

	// log.Print("Drawing citites...")
	// Draw on circles representing nodes
	for _, node := range world.Cities {
		cv.DrawCircle(color.RGBA{22, 131, 201, 255}, node.Position, 5)
	}
	// blure is quite expensive
	// canvas.Blur(1, new(WeightFunctionDist))
	// log.Print("Done.")

	// log.Print("Drawing roads...")
	// Draw connections between nodes
	for _, node := range world.Cities {
		for _, peer := range node.Peers[:3] {
			cv.DrawLine(color.RGBA{0, 0, 0, 10}, node.Position, peer.Position)
		}
	}
	log.Print("Done.")

	return world
}

func (w *World) SendMessages(n int) {
	for i := 0; i <= n; i++ {
		city := w.Cities[i]
		city.Power = 255
		go city.Send()
	}
}

type Node struct {
	Position canvas.Vector
	Ch       chan *Node
	Peers    []*Node
	Canvas   *canvas.Canvas
	Power    uint8
}

func NewNode(peers int, cv *canvas.Canvas) *Node {
	node := new(Node)
	node.Peers = make([]*Node, 0, peers)
	size := cv.Bounds().Size()
	x := float64(size.X) * rand.Float64()
	y := float64(size.Y) * rand.Float64()
	node.Position = canvas.Vector{x, y}
	node.Canvas = cv
	node.Ch = make(chan *Node)
	node.Power = 0
	go node.Listen()
	return node
}

func (n *Node) Listen() {
	for {
		peer := <-n.Ch
		peer.Power -= 5
		n.Power = peer.Power
		n.Canvas.DrawLine(color.RGBA{255, n.Power, 0, 255}, n.Position, peer.Position)
		if n.Power > 0 {
			go n.Send()
		}
	}
}

func (n *Node) Send() {
	for _, target := range n.Peers {
		if target.Power == 0 {
			target.Ch <- n
			break
		}
	}
}

type NodeSorter struct {
	data   []*Node
	target *Node
}

func (sorter NodeSorter) Len() int { return len(sorter.data) }
func (sorter NodeSorter) Less(i, j int) bool {
	iDelta := sorter.data[i].Position.Sub(&sorter.target.Position)
	jDelta := sorter.data[j].Position.Sub(&sorter.target.Position)
	return iDelta.Length() < jDelta.Length()
}
func (sorter NodeSorter) Swap(i, j int) {
	sorter.data[i], sorter.data[j] = sorter.data[j], sorter.data[i]
}

type WeightFunctionDist struct{}

func (w WeightFunctionDist) Weight(x int, y int) float64 {
	d := math.Hypot(float64(x), float64(y))
	return 1 / (1 + d)
}
