package glvox

import (
//	"fmt"
	"bufio"
	"io"
	"os"
	"strings"
	"strconv"
	"errors"
)

func ReadBinvox(filename string, voxels DimSetter) (err error) {

	f, err := os.Open(filename)
	if err != nil { return }

	buf := bufio.NewReader(f)

	var line string
	line, err = buf.ReadString('\n')
	if err != nil { return }

	if !strings.HasPrefix(line, "#binvox") {
		err = errors.New("not a binvox file")
		return
	}

	w, h, d := 0, 0, 0
	data := false
	for {
		line, err = buf.ReadString('\n')
		if strings.HasPrefix(line, "dim") {
			dims := strings.Fields(line)
			d, err = strconv.Atoi(dims[1])
			if err != nil { return }
			h, err = strconv.Atoi(dims[2])
			if err != nil { return }
			w, err = strconv.Atoi(dims[3])
			if err != nil { return }
		} else if strings.HasPrefix(line, "data") {
			data = true
			break;
		}
	}

	if w == 0 && h == 0 && d == 0 {
		err = errors.New("no dim header")
		return
	}

	if !data {
		err = errors.New("no data")
		return
	}

	voxels.Dim(int32(w), int32(h), int32(d))

	x, y, z := 0, 0, 0
	for {
		var val, cnt byte

		val, err = buf.ReadByte()
		if err != nil {
			if err == io.EOF { err = nil }
			return
		}

		cnt, err = buf.ReadByte()
		if err != nil {
			return
		}

		for i := 0; i < int(cnt); i++ {

			voxels.Set(int32(x), int32(y), int32(z), int32(val))

			x++
			if x >= w {
				x = 0; y++
				if y >= h {
					y = 0; z++
					if z >= d { /* error */ }
				}
			}
		}
	}

	return
}
