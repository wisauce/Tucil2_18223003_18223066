package main

import (
	"fmt"
	"math"
)

type Vertice struct {
	x, y, z float64
}

type Face struct {
	a, b, c Vertice
}

type BoundingBox struct {
	min, max Vertice
}

type OctreeNode struct {
	depth int
	children [8]*OctreeNode
	faces []Face
	boundingBox BoundingBox
}

func (node *OctreeNode) ComputeRootBound(faces [] Face) BoundingBox {
  minX := faces[0].a.x
  minY := faces[0].a.y
  minZ := faces[0].a.z
  maxX := faces[0].a.x
  maxY := faces[0].a.y
  maxZ := faces[0].a.z

  for _, t := range faces {
    for _, v := range [] Vertice{t.a, t.b, t.c} {
      if v.x < minX { minX = v.x }
      if v.y < minY { minY = v.y }
      if v.z < minZ { minZ = v.z }

      if v.x > maxX { maxX = v.x }
      if v.y > maxY { maxY = v.y }
      if v.z > maxZ { maxZ = v.z }
    }
  }

  sizeX := maxX - minX
  sizeY := maxY - minY
  sizeZ := maxZ - minZ
  maxSize := math.Max(math.Max(sizeX, sizeY), sizeZ)

  return BoundingBox{
    min: Vertice{x: minX, y: minY, z: minZ},
    max: Vertice{x: minX + maxSize, y: minY + maxSize, z: minZ + maxSize},
  }
}

func (node *OctreeNode) ComputeChildBound (parent BoundingBox, i int) BoundingBox {
  min := parent.min
  max := parent.max

  midX := (min.x + max.x) / 2
  midY := (min.y + max.y) / 2
  midZ := (min.z + max.z) / 2

  var childMin, childMax Vertice

  if i&1 == 0 {
    childMin.x = min.x
    childMax.x = midX
  } else {
    childMin.x = midX
    childMax.x = max.x
  }

  if i&2 == 0 {
    childMin.y = min.y
    childMax.y = midY
  } else {
    childMin.y = midY
    childMax.y = max.y
  }

  if i&4 == 0 {
    childMin.z = min.z
    childMax.z = midZ
  } else {
    childMin.z = midZ
    childMax.z = max.z
  }

  return BoundingBox{min: childMin, max: childMax}
}


func BuildOctree(faces []Face, depth int) *OctreeNode {
	node := &OctreeNode{}
	node.depth = depth
	node.children = [8]*OctreeNode{}
	node.faces = faces
	node.boundingBox = node.ComputeRootBound(faces)
	return node
}


func main() {
	// parse object and convert to faces
	faces := []Face{
		{a: Vertice{x: 0, y: 0, z: 0}, b: Vertice{x: 1, y: 0, z: 0}, c: Vertice{x: 0, y: 1, z: 0}},
		{a: Vertice{x: 1, y: 0, z: 0}, b: Vertice{x: 1, y: 2, z: 0}, c: Vertice{x: 0, y: 1, z: 0}},
	}
	octree := BuildOctree(faces, 2)
	fmt.Println(octree.boundingBox)

}