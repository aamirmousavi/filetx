package filetx_test

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"testing"

	"github.com/aamirmousavi/filetx"
)

func TestTx(t *testing.T) {
	tx, err := filetx.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	op, err := tx.Create("/Users/amir/dev/go/lab/filetx/my_my_test.txt")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := op.Write([]byte("Hello, World!\n2121")); err != nil {
		t.Fatal(err)
	}

	img, err := tx.Create("/Users/amir/dev/go/lab/filetx/my_my_test.png")
	if err != nil {
		t.Fatal(err)
	}

	realImg, err := os.ReadFile("/Users/amir/dev/go/lab/filetx/img.png")
	if err != nil {
		t.Fatal(err)
	}

	decodedImg, _, err := image.Decode(bytes.NewReader(realImg))
	if err != nil {
		t.Fatal(err)
	}

	if err := png.Encode(img, decodedImg); err != nil {
		t.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
}
