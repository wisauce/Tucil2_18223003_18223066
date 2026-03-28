package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Vertex struct {
	x, y, z float64
}

type Face struct {
	a, b, c Vertex
}

type BoundingBox struct {
	min, max Vertex
}

type OctreeNode struct {
	depth int
	children [8]*OctreeNode
	faces []Face
	boundingBox BoundingBox
  skipped bool
}

func ComputeRootBound(faces [] Face) BoundingBox {
  minX := faces[0].a.x
  minY := faces[0].a.y
  minZ := faces[0].a.z
  maxX := faces[0].a.x
  maxY := faces[0].a.y
  maxZ := faces[0].a.z
  for _, t := range faces {
    for _, v := range [] Vertex{t.a, t.b, t.c} {
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
    min: Vertex{x: minX, y: minY, z: minZ},
    max: Vertex{x: minX + maxSize, y: minY + maxSize, z: minZ + maxSize},
  }
}

func faceIntersectsBox(face Face, box BoundingBox) bool {
  // move triangle so box center at origin
  cx := (box.min.x + box.max.x) * 0.5
  cy := (box.min.y + box.max.y) * 0.5
  cz := (box.min.z + box.max.z) * 0.5

  hx := (box.max.x - box.min.x) * 0.5
  hy := (box.max.y - box.min.y) * 0.5
  hz := (box.max.z - box.min.z) * 0.5

  v0x := face.a.x - cx
  v0y := face.a.y - cy
  v0z := face.a.z - cz

  v1x := face.b.x - cx
  v1y := face.b.y - cy
  v1z := face.b.z - cz

  v2x := face.c.x - cx
  v2y := face.c.y - cy
  v2z := face.c.z - cz

  // compute edges
  e0x := v1x - v0x
  e0y := v1y - v0y
  e0z := v1z - v0z

  e1x := v2x - v1x
  e1y := v2y - v1y
  e1z := v2z - v1z

  e2x := v0x - v2x
  e2y := v0y - v2y
  e2z := v0z - v2z

  // 1 - test axes L = edge cross axes
  // e0
  p0 := v0z*e0y - v0y*e0z
  p2 := v2z*e0y - v2y*e0z
  minP, maxP := p0, p2
  if minP > maxP { minP, maxP = maxP, minP }
  rad := hz*math.Abs(e0y) + hy*math.Abs(e0z)
  if minP > rad || maxP < -rad { return false }

  p0 = v0x*e0z - v0z*e0x
  p2 = v2x*e0z - v2z*e0x
  minP, maxP = p0, p2
  if minP > maxP { minP, maxP = maxP, minP }
  rad = hx*math.Abs(e0z) + hz*math.Abs(e0x)
  if minP > rad || maxP < -rad { return false }

  p1 := v1y*e0x - v1x*e0y
  p2 = v2y*e0x - v2x*e0y
  minP, maxP = p1, p2
  if minP > maxP { minP, maxP = maxP, minP }
  rad = hx*math.Abs(e0y) + hy*math.Abs(e0x)
  if minP > rad || maxP < -rad { return false }

  // e1
  p0 = v0z*e1y - v0y*e1z
  p2 = v2z*e1y - v2y*e1z
  minP, maxP = p0, p2
  if minP > maxP { minP, maxP = maxP, minP }
  rad = hz*math.Abs(e1y) + hy*math.Abs(e1z)
  if minP > rad || maxP < -rad { return false }

  p0 = v0x*e1z - v0z*e1x
  p2 = v2x*e1z - v2z*e1x
  minP, maxP = p0, p2
  if minP > maxP { minP, maxP = maxP, minP }
  rad = hx*math.Abs(e1z) + hz*math.Abs(e1x)
  if minP > rad || maxP < -rad { return false }

  p0 = v0y*e1x - v0x*e1y
  p1 = v1y*e1x - v1x*e1y
  minP, maxP = p0, p1
  if minP > maxP { minP, maxP = maxP, minP }
  rad = hx*math.Abs(e1y) + hy*math.Abs(e1x)
  if minP > rad || maxP < -rad { return false }

  // e2
  p0 = v0z*e2y - v0y*e2z
  p1 = v1z*e2y - v1y*e2z
  minP, maxP = p0, p1
  if minP > maxP { minP, maxP = maxP, minP }
  rad = hz*math.Abs(e2y) + hy*math.Abs(e2z)
  if minP > rad || maxP < -rad { return false }

  p0 = v0x*e2z - v0z*e2x
  p1 = v1x*e2z - v1z*e2x
  minP, maxP = p0, p1
  if minP > maxP { minP, maxP = maxP, minP }
  rad = hx*math.Abs(e2z) + hz*math.Abs(e2x)
  if minP > rad || maxP < -rad { return false }

  p1 = v1y*e2x - v1x*e2y
  p2 = v2y*e2x - v2x*e2y
  minP, maxP = p1, p2
  if minP > maxP { minP, maxP = maxP, minP }
  rad = hx*math.Abs(e2y) + hy*math.Abs(e2x)
  if minP > rad || maxP < -rad { return false }

  // 2 - test bounding box axes
  minX, maxX := v0x, v0x
  if v1x < minX { minX = v1x }
  if v1x > maxX { maxX = v1x }
  if v2x < minX { minX = v2x }
  if v2x > maxX { maxX = v2x }
  if minX > hx || maxX < -hx { return false }

  minY, maxY := v0y, v0y
  if v1y < minY { minY = v1y }
  if v1y > maxY { maxY = v1y }
  if v2y < minY { minY = v2y }
  if v2y > maxY { maxY = v2y }
  if minY > hy || maxY < -hy { return false }

  minZ, maxZ := v0z, v0z
  if v1z < minZ { minZ = v1z }
  if v1z > maxZ { maxZ = v1z }
  if v2z < minZ { minZ = v2z }
  if v2z > maxZ { maxZ = v2z }
  if minZ > hz || maxZ < -hz { return false }

  // 3 - triangle plane vs box
  nx := e0y*e1z - e0z*e1y
  ny := e0z*e1x - e0x*e1z
  nz := e0x*e1y - e0y*e1x

  d := -(nx*v0x + ny*v0y + nz*v0z)

  r := hx*math.Abs(nx) + hy*math.Abs(ny) + hz*math.Abs(nz)
  s := d

  if s > r || s < -r {
    return false
  }

  return true
}


func FilterFaces(faces []Face, box BoundingBox) []Face {
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

  var childMin, childMax Vertex

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
      boundingBox: bound,
  }

  if len(faces) == 0 {
    node.skipped = true
    return node
  }

  if depth == 0 {
    node.faces = faces
    return node
  }

  var wg sync.WaitGroup
  for i := 0; i < 8; i++ {
    wg.Add(1)
    go func(childIdx int) {
      defer wg.Done()
      childBound := ComputeChildBound(node.boundingBox, childIdx)
      filteredFaces := FilterFaces(faces, childBound)
      node.children[childIdx] = BuildOctree(filteredFaces, childBound, depth-1)
    }(i)
  }
  wg.Wait()
  
	return node
}

func PrintOctree(node *OctreeNode) {
  var printNode func(node *OctreeNode, depth int)
  printNode = func(node *OctreeNode, depth int) {
    if node == nil {
      fmt.Printf("%snil\n", strings.Repeat("│   ", depth))
      return
    }
    prefix := strings.Repeat("│   ", depth)
    fmt.Printf("%sNode(depth=%d, faces=%d)\n", prefix, node.depth, len(node.faces))
    fmt.Printf("%sBoundingBox(min=(%.2f, %.2f, %.2f), max=(%.2f, %.2f, %.2f))\n",
      prefix, node.boundingBox.min.x, node.boundingBox.min.y, node.boundingBox.min.z,
      node.boundingBox.max.x, node.boundingBox.max.y, node.boundingBox.max.z)
    for i, child := range node.children {
      if child != nil {
        fmt.Printf("%s├── Child %d:\n", prefix, i)
        printNode(child, depth+1)
      }
    }
  }
  printNode(node, 0)
}

func ParseObj(filename string) ([]Face, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka file: %w", err)
	}
	defer file.Close()

	var vertices []Vertex
	var faces []Face
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Fields(line)
		switch parts[0] {
		case "v":
			if len(parts) >= 4 {
				x, _ := strconv.ParseFloat(parts[1], 64)
				y, _ := strconv.ParseFloat(parts[2], 64)
				z, _ := strconv.ParseFloat(parts[3], 64)
				vertices = append(vertices, Vertex{x: x, y: y, z: z})
			}
		case "f":
			if len(parts) >= 4 {
				var indices []int
				for i := 1; i < len(parts); i++ {
					vals := strings.Split(parts[i], "/")
					idx, _ := strconv.Atoi(vals[0])
					indices = append(indices, idx-1)
				}
				for i := 1; i+1 < len(indices); i++ {
					if indices[0] < len(vertices) && indices[i] < len(vertices) && indices[i+1] < len(vertices) {
						faces = append(faces, Face{
							a: vertices[indices[0]],
							b: vertices[indices[i]],
							c: vertices[indices[i+1]],
						})
					}
				}
			}
		}
	}

	return faces, scanner.Err()
}

func CollectVoxelBoxes(node *OctreeNode, boxes *[]BoundingBox) {
	if node == nil || node.skipped {
		return
	}
	if node.depth == 0 {
		if len(node.faces) > 0 {
			*boxes = append(*boxes, node.boundingBox)
		}
		return
	}
	for i := 0; i < 8; i++ {
		CollectVoxelBoxes(node.children[i], boxes)
	}
}

func ExportToObj(boxes []BoundingBox, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	vertexOffset := 1

	for _, box := range boxes {
		min, max := box.min, box.max
		vertices := []Vertex{
			{min.x, min.y, min.z}, {max.x, min.y, min.z}, {max.x, max.y, min.z}, {min.x, max.y, min.z},
			{min.x, min.y, max.z}, {max.x, min.y, max.z}, {max.x, max.y, max.z}, {min.x, max.y, max.z},
		}

		for _, v := range vertices {
			fmt.Fprintf(writer, "v %.6f %.6f %.6f\n", v.x, v.y, v.z)
		}

		facesIdx := [][]int{
			{1, 3, 2}, {1, 4, 3}, {5, 6, 7}, {5, 7, 8},
			{1, 2, 6}, {1, 6, 5}, {4, 8, 7}, {4, 7, 3},
			{1, 5, 8}, {1, 8, 4}, {2, 3, 7}, {2, 7, 6},
		}

		for _, f := range facesIdx {
			fmt.Fprintf(writer, "f %d %d %d\n", f[0]-1+vertexOffset, f[1]-1+vertexOffset, f[2]-1+vertexOffset)
		}
		vertexOffset += 8
	}

	return writer.Flush()
}

func CollectStats(node *OctreeNode, maxDepth int, created []int, skipped []int) {
	if node == nil {
		return
	}
	level := maxDepth - node.depth
	created[level]++

	if node.skipped {
		skipped[level]++
		return 
	}
	for i := 0; i < 8; i++ {
		if node.children[i] != nil {
			CollectStats(node.children[i], maxDepth, created, skipped)
		}
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Cara penggunaan: go run main.go <path_file.obj> <depth_maksimum>")
		return
	}

	filename := os.Args[1]
	depth, err := strconv.Atoi(os.Args[2])
	if err != nil || depth < 0 {
		fmt.Println("Error: Parameter kedalaman (depth) harus berupa bilangan bulat positif.")
		return
	}

	fmt.Printf("Memproses file: %s dengan kedalaman maksimal: %d\n", filename, depth)

	faces, err := ParseObj(filename)
	if err != nil {
		fmt.Printf("Error validasi input: %v\n", err)
		return
	}
	if len(faces) == 0 {
		fmt.Println("Error: Tidak ada data bidang (faces) yang valid di file OBJ.")
		return
	}

  startTime := time.Now()
	octree := BuildOctree(faces, ComputeRootBound(faces), depth)
	
	var voxelBoxes []BoundingBox
	CollectVoxelBoxes(octree, &voxelBoxes)

  createdNodes := make([]int, depth+1)
	skippedNodes := make([]int, depth+1)
	CollectStats(octree, depth, createdNodes, skippedNodes)

	elapsedTime := time.Since(startTime)

	baseName := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
	outPath := filepath.Join("test", fmt.Sprintf("%s-voxelized-%d.obj", baseName, depth))
	
	if err := ExportToObj(voxelBoxes, outPath); err != nil {
		fmt.Printf("Error saat menyimpan file output: %v\n", err)
		return
	}

	absPath, _ := filepath.Abs(outPath)
	numVoxels := len(voxelBoxes)

	fmt.Println("\n==================================")
	fmt.Printf("Banyaknya voxel yang terbentuk: %d\n", numVoxels)
	fmt.Printf("Banyaknya vertex yang terbentuk: %d\n", numVoxels*8)
	fmt.Printf("Banyaknya faces yang terbentuk: %d\n", numVoxels*12)
	
	fmt.Println("\nStatistik Octree:")
	fmt.Println("Kedalaman | Terbentuk | Dilewati")
	fmt.Println("----------------------------------")
	for i := 1; i <= depth; i++ {
		fmt.Printf("Level %-3d | %-9d | %-8d\n", i, createdNodes[i], skippedNodes[i])
	}
	
	fmt.Printf("\nKedalaman octree: %d\n", depth)
	fmt.Printf("Lama waktu eksekusi: %v\n", elapsedTime)
	fmt.Printf("Lokasi file disimpan: %s\n", absPath)
	fmt.Println("==================================")
}