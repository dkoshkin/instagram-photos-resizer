package instagram

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		name   string
		file   string
		splits int
		width  int
		height int
		err    error
	}{
		{
			name: "1080x960 error",
			file: "1080x960.jpg",
			err:  fmt.Errorf("the provided image is is not a supported aspect ratio"),
		},
		{
			name:   "2160x1080 split",
			file:   "2160x1080.jpg",
			splits: 2,
			width:  1080,
			height: 1080,
		},
		{
			name:   "2160x1350 split",
			file:   "2160x1350.jpg",
			splits: 2,
			width:  1080,
			height: 1350,
		},
		{
			name:   "2160x2160 resize",
			file:   "2160x2160.jpg",
			splits: 1,
			width:  1080,
			height: 1080,
		},
		{
			name:   "2160x2700 resize",
			file:   "2160x2700.jpg",
			splits: 1,
			width:  1080,
			height: 1350,
		},
		{
			name:   "2700x2700 resize",
			file:   "2700x2700.jpg",
			splits: 1,
			width:  1080,
			height: 1080,
		},
		{
			name:   "3240x1350.jpg split",
			file:   "3240x1350.jpg",
			splits: 3,
			width:  1080,
			height: 1350,
		},
		{
			name:   "4000x2500.jpg.jpg split resize",
			file:   "4000x2500.jpg",
			splits: 2,
			width:  1080,
			height: 1350,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			b, err := ioutil.ReadFile(test.file)
			if err != nil {
				t.Fatal(err)
			}
			preparedImgs, err := Prepare(bytes.NewReader(b))
			if (err != nil && test.err == nil) || (err != nil && err.Error() != test.err.Error()) || (err == nil && test.err != nil) {
				t.Fatalf("expected: %v, got: %v", test.err, err)
			}
			if err == nil {
				if len(preparedImgs.Reader) != test.splits {
					t.Fatalf("expected splits: %d, got: %d", test.splits, len(preparedImgs.Reader))
				}
				for n := range preparedImgs.Reader {
					if preparedImgs.Reader[n] == nil {
						t.Fatalf("nil reader for split %d", n)
					}
				}
				if preparedImgs.Width != test.width {
					t.Fatalf("expected width: %d, got: %d", test.width, preparedImgs.Width)
				}
				if preparedImgs.Height != test.height {
					t.Fatalf("expected width: %d, got: %d", test.height, preparedImgs.Height)
				}
			}
		})
	}
}

func TestCanSplitInto(t *testing.T) {
	tests := []struct {
		width  int
		height int
		ratio  float64
		splits int
	}{
		{
			width:  1080,
			height: 1080,
			ratio:  1.0,
			splits: 1,
		},
		{
			width:  1080,
			height: 1350,
			ratio:  1.25,
			splits: 1,
		},
		{
			width:  2160,
			height: 1080,
			ratio:  1.0,
			splits: 2,
		},
		{
			width:  2160,
			height: 1350,
			ratio:  1.25,
			splits: 2,
		},
		{
			width:  4320,
			height: 1080,
			ratio:  1.0,
			splits: 4,
		},
		{
			width:  4320,
			height: 1350,
			ratio:  1.25,
			splits: 4,
		},
		{
			width:  1080,
			height: 960,
			ratio:  1.0,
			splits: 0,
		},
		{
			width:  1080,
			height: 960,
			ratio:  1.25,
			splits: 0,
		},
		{
			width:  960,
			height: 1080,
			ratio:  1.0,
			splits: 0,
		},
		{
			width:  960,
			height: 1350,
			ratio:  1.25,
			splits: 0,
		},
		{
			width:  2260,
			height: 1080,
			ratio:  1.0,
			splits: 0,
		},
		{
			width:  2260,
			height: 1350,
			ratio:  1.25,
			splits: 0,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test), func(t *testing.T) {
			splits := canSplitInto(test.width, test.height, test.ratio)
			if test.splits != splits {
				t.Fatalf("expected: %d, got %d", test.splits, splits)
			}
		})
	}
}
