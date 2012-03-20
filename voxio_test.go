package glvox_test

import (
	"github.com/shogg/glvox"
	"testing"
	"fmt"
)

func OffTestReadBinvox(t *testing.T) {

	voxels := glvox.NewOctree()
	err := glvox.ReadBinvox("skull.binvox", voxels)
	if err != nil {
		t.Error(err)
	}

	if voxels.WHD != 256 {
		t.Error("dimension 256 expected, was", voxels.WHD)
	}

	indexCount := len(voxels.Index) / 8
	if indexCount != 139680 {
		t.Error("index size 139680 expected, was", indexCount)
	}

	longJump := int32(0)
	for i := int32(0); i < int32(indexCount); i++ {
		for j := int32(0); j < 8; j++ {
			jump := voxels.Index[i*8 + j] - i
			if jump > longJump { longJump = jump }
		}
	}
	fmt.Println("longest jump", longJump)

	avgJump := int32(0)
	for i := int32(0); i < int32(indexCount); i++ {
		for j := int32(0); j < 8; j++ {
			jump := voxels.Index[i*8 + j] - i
			avgJump += jump
		}
	}
	avgJump /= int32(indexCount * 8)
	fmt.Println("average jump", avgJump)

	countBySize := make(map [int32] int32)
	for x := int32(0); x < voxels.WHD; x++ {
		for y := int32(0); y < voxels.WHD; y++ {
			for z := int32(0); z < voxels.WHD; z++ {
				val, size := voxels.Get(x, y, z)
				if val == 0 {
					countBySize[size] += 1
				}
			}
		}
	}

	fmt.Println("voxels by size:")
	var sizes = []int32 { 1, 2, 4, 8, 16, 32, 64, 128, 256 }
	for _, size := range(sizes) {
		count := countBySize[size]
		unit := ""
		if count > 1000 { count /= 1000; unit = "k" }
		fmt.Printf("%d\t%8d%s\n", size, count, unit)
	}
}
