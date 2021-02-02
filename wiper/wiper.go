package wiper

import (
	"crypto/rand"
	"io"
	"math"
	"os"
)

const blockSize = 1024 * 1024

//Wipe file with basic wiping pattern (zerofill-randfill-randfill)
func Wipe(filename string) error {
	randbuf := make([]byte, blockSize)
	f, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	stat, err := f.Stat()
	if err != nil {
		return err
	}
	size := stat.Size()
	blockCount := int(math.Ceil(float64(size) / blockSize))
	counter := 0
	for pass := 0; pass < 3; pass++ {
		counter = 0
		for i := 0; i < blockCount; i++ {
			/*
				if pass == 0 {
					f.Read(randbuf[:blockSize])
					for i := range randbuf {
						randbuf[i] = 255 - randbuf[i]
					}
				}
			*/
			f.WriteAt(randbuf, int64(counter*blockSize))
			counter++
			if pass != 0 {
				io.ReadFull(rand.Reader, randbuf)
			}
		}
		io.ReadFull(rand.Reader, randbuf)
	}
	f.Sync()
	f.Close()
	return os.Remove(filename)
	//return nil
}
