package glvox_test

import (
	"github.com/shogg/glvox"
	"testing"
	"fmt"
)

func TestReadBinvox(t *testing.T) {

	voxels := glvox.NewOctree(1256)
	err := glvox.ReadBinvox("skull.binvox", voxels, 0, 0, 0)
	if err != nil {
		t.Error(err)
	}

	if voxels.Size <= 0 {
		t.Error("dimension > 0 expected, was", voxels.Size)
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
}
