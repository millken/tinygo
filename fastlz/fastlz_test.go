package fastlz

import (
	"bytes"
	"crypto/rand"
	"os"
	"testing"
)

func TestFastLZ(t *testing.T) {
	// Test FlzCompress
	ib, err := os.ReadFile("fastlz.go")
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}
	ob := FlzCompress(ib)
	//compressing the file should not increase the size
	if len(ob) > len(ib) {
		t.Errorf("Compressed size is greater than original size")
	}
	//calculate the decompressed size
	dob := FlzDecompress(ob)
	if len(dob) != len(ib) {
		t.Errorf("Decompressed size is not equal to original size")
	}
	//calculate the compress ratio
	ratio := float64(len(ob)) / float64(len(ib))
	t.Logf("Compress ratio: %f", ratio)
	//compare the decompressed data with the original data
	for i := 0; i < len(ib); i++ {
		if ib[i] != dob[i] {
			t.Errorf("Decompressed data is not equal to original data")
		}
	}

}
func genData(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

func TestAll(t *testing.T) {
	for _, i := range []int{10, 128, 1000, 1024 * 10, 1024 * 100, 1024 * 1024, 1024 * 1024 * 7} {
		data, err := genData(i)
		if err != nil {
			t.Error(err)
			continue
		}
		comp := FlzCompress(data)
		decomp := FlzDecompress(comp)
		if !bytes.Equal(decomp, data) {
			t.Fatalf("deflate->inflate does not match original for %d", i)
		}

	}
}
