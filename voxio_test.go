package glvox_test

import (
	"github.com/shogg/glvox"
	"testing"
)

func TestReadBinvox(t *testing.T) {

	voxels := glvox.NewOctree()
	err := glvox.ReadBinvox("skull.binvox", voxels)
	if err != nil {
		t.Error(err)
	}

	if voxels.WHD != 256 {
		t.Error("dimension 256 expected, was", voxels.WHD)
	}

	indexCount := len(voxels.Index) / 9
	if indexCount != 139680 {
		t.Error("index size expected 139680, was", indexCount)
	}
}
