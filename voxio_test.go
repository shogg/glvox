package glvox_test

import (
	"github.com/shogg/glvox"
	"testing"
	"fmt"
)

func TestReadBinvox(t *testing.T) {

	voxels := glvox.NewOctree()
	voxels.Dim(1256, 1256, 1256)
	err := glvox.ReadBinvox("skull.binvox", voxels, 0, 0, 0)
	if err != nil {
		t.Error(err)
	}

	if voxels.WHD <= 0 {
		t.Error("dimension > 0 expected, was", voxels.WHD)
	}

	indexCount := len(voxels.Index) / 8
	if indexCount != 642216 {
		t.Error("index size 642216 expected, was", indexCount)
	}

	longJump := int32(0)
	for i := int32(0); i < int32(indexCount); i++ {
		for j := int32(0); j < 8; j++ {
			idx := voxels.Index[i<<3 + j]
			if idx <= 0 { continue; }

			jump := idx - i
			if jump > longJump { longJump = jump }
		}
	}
	fmt.Println("longest jump", longJump)

	avgJump := int32(0)
	jumpCount := int32(0)
	for i := int32(0); i < int32(indexCount); i++ {
		for j := int32(0); j < 8; j++ {
			idx := voxels.Index[i<<3 + j]
			if idx <= 0 { continue; }

			avgJump += idx - i
			jumpCount++
		}
	}
	avgJump /= jumpCount
	fmt.Println("average jump", avgJump)
/*
	countBySize := make(map [int32] int32)
	for x := int32(0); x < voxels.WHD; x++ {
		for y := int32(0); y < voxels.WHD; y++ {
			for z := int32(0); z < voxels.WHD; z++ {
				_, size := voxels.Get(x, y, z)
				countBySize[size] += 1
			}
		}
	}

	fmt.Println("voxels by size:")
	var sizes = []int32 { 1, 2, 4, 8, 16, 32, 64, 128, 256 }
	for _, size := range(sizes) {
		count := countBySize[size] / (size*size*size)
		unit := ""
//		if count > 10000 { count /= 1000; unit = "k" }
		fmt.Printf("%d\t%8d%s\n", size, count, unit)
	}*/
}
