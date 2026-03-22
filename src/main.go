package main

type Point struct {
	x, y, z float64
}

type Triangle struct {
	a, b, c Point
}

type BoundingBox struct {
	min, max Point
}

type OctreeNode struct {
	depth int
	children [8]*OctreeNode
	triangles []Triangle
	boundingBox BoundingBox
}