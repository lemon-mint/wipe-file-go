package wiper

import (
	"crypto/rand"
	"encoding/base32"
	"io"
	"math"
	"os"
	"path/filepath"
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
	blockCount := int(math.Ceil(float64(size)/blockSize)) + 1
	for pass := 0; pass < 3; pass++ {
		for i := 0; i < blockCount; i++ {
			/*
				if pass == 0 {
					f.Read(randbuf[:blockSize])
					for i := range randbuf {
						randbuf[i] = 255 - randbuf[i]
					}
				}
			*/
			f.WriteAt(randbuf, int64(i*blockSize))
			if pass != 0 {
				io.ReadFull(rand.Reader, randbuf)
			}
		}
		io.ReadFull(rand.Reader, randbuf)
		f.Sync()
	}
	f.Close()
	dir, _ := filepath.Split(filename)
	newname := filepath.Join(dir, randb32(max(len(filename), 20)))
	//fmt.Println(filename)
	//fmt.Println(dir)
	for i := 0; i < 10; i++ {
		//fmt.Println(newname)
		os.Rename(filename, newname)
		filename = newname
		newname = filepath.Join(dir, randb32(max(len(filename), 20)))
		//time.Sleep(time.Second * 2)
	}
	return os.Remove(filename)
	//return nil
}

//RandFill : fill file with random
func RandFill(filename string, count int) error {
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
	blockCount := int(math.Ceil(float64(size)/blockSize)) + 1
	for pass := 0; pass < count; pass++ {
		for i := 0; i < blockCount; i++ {
			io.ReadFull(rand.Reader, randbuf)
			f.WriteAt(randbuf, int64(i*blockSize))
		}
		f.Sync()
	}
	f.Close()
	return nil
}

//ZeroFill : fill file with zero
func ZeroFill(filename string, count int) error {
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
	blockCount := int(math.Ceil(float64(size)/blockSize)) + 1
	for pass := 0; pass < count; pass++ {
		for i := 0; i < blockCount; i++ {
			f.WriteAt(randbuf, int64(i*blockSize))
		}
		f.Sync()
	}
	f.Close()
	return nil
}

//MixFileName : change file name to random
func MixFileName(filename string, count int) (newname string, err error) {
	dir, _ := filepath.Split(filename)
	newname = filepath.Join(dir, randb32(20))
	//fmt.Println(filename)
	//fmt.Println(dir)
	for i := 0; i < count; i++ {
		//fmt.Println(newname)
		err = os.Rename(filename, newname)
		if err != nil {
			return
		}
		filename = newname
		newname = filepath.Join(dir, randb32(max(len(filename), 20)))
		//time.Sleep(time.Second * 2)
	}
	newname = filename
	err = nil
	return
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func randb32(size int) string {
	buf := make([]byte, size)
	io.ReadFull(rand.Reader, buf)
	return base32.StdEncoding.EncodeToString(buf)
}
