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

func ComputeRootBound(faces [] Face) BoundingBox {
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

func faceIntersectsBox(face Face, box BoundingBox) bool {
  // Compute triangle AABB (inline min/max)

  minX := face.a.x
  if face.b.x < minX { minX = face.b.x }
  if face.c.x < minX { minX = face.c.x }

  maxX := face.a.x
  if face.b.x > maxX { maxX = face.b.x }
  if face.c.x > maxX { maxX = face.c.x }

  minY := face.a.y
  if face.b.y < minY { minY = face.b.y }
  if face.c.y < minY { minY = face.c.y }

  maxY := face.a.y
  if face.b.y > maxY { maxY = face.b.y }
  if face.c.y > maxY { maxY = face.c.y }

  minZ := face.a.z
  if face.b.z < minZ { minZ = face.b.z }
  if face.c.z < minZ { minZ = face.c.z }

  maxZ := face.a.z
  if face.b.z > maxZ { maxZ = face.b.z }
  if face.c.z > maxZ { maxZ = face.c.z }

  // AABB vs AABB overlap
  return (minX <= box.max.x && maxX >= box.min.x) &&
         (minY <= box.max.y && maxY >= box.min.y) &&
         (minZ <= box.max.z && maxZ >= box.min.z)
}

func filterFaces(faces []Face, box BoundingBox) []Face {
  filtered := make([]Face, 0, len(faces))

  for _, face := range faces {
    if faceIntersectsBox(face, box) {
      filtered = append(filtered, face)
    }
  }

  return filtered
}

func ComputeChildBound (parent BoundingBox, i int) BoundingBox {
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


func BuildOctree(faces []Face, bound BoundingBox, depth int) *OctreeNode {
  node := &OctreeNode{
      depth: depth,
      children: [8]*OctreeNode{},
      faces: faces, // remember to filter faces for each child node in a real implementation
      boundingBox: bound,
  }
  if depth > 0 {
    for i := range 8 {
      childBound := ComputeChildBound(node.boundingBox, i)
      // In a real implementation, you would filter faces that intersect with childBound
      filteredFaces := filterFaces(faces, childBound)
      node.children[i] = BuildOctree(filteredFaces, childBound, depth-1)
    }
  }
  
	return node
}

func PrintOctree(node *OctreeNode, indent string) {
  if node == nil {
    fmt.Println(indent + "nil")
    return
  }
  fmt.Printf("%sNode(depth=%d, faces=%d)\n", indent, node.depth, len(node.faces))
  for i, child := range node.children {
    if child != nil {
      fmt.Printf("%s├── Child %d:\n", indent, i)
      PrintOctree(child, indent+"│   ")
    }
  }
}

func main() {
	faces := []Face{
		{a: Vertice{x: 0, y: 0, z: 0}, b: Vertice{x: 1, y: 0, z: 0}, c: Vertice{x: 0, y: 1, z: 0}},
		{a: Vertice{x: 1, y: 0, z: 0}, b: Vertice{x: 1, y: 2, z: 0}, c: Vertice{x: 0, y: 1, z: 0}},
	}
	rootBound := ComputeRootBound(faces)
	octree := BuildOctree(faces, rootBound,2)
  PrintOctree(octree, "")

}