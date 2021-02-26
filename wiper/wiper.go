// MIT License
//
// Copyright (c) 2021 lemon-mint
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package wiper

import (
	"crypto/rand"
	"encoding/base32"
	"io"
	"math"
	mrand "math/rand"
	"os"
	"path/filepath"
	"time"
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

//ComplementFill : fill file with Complement
func ComplementFill(filename string) error {
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
	defer func() {
		f.Sync()
		f.Close()
	}()
	blockCount := int(math.Ceil(float64(size) / blockSize))
	for i := 0; i < blockCount; i++ {
		_, err = f.Read(randbuf[:blockSize])
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		for i := range randbuf {
			randbuf[i] = 255 - randbuf[i]
		}
		f.WriteAt(randbuf, int64(i*blockSize))
	}
	return nil
}

//MixFileName : change file name to random
func MixFileName(filename string, count int) (newname string, err error) {
	originalname := filename
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
		newname = filepath.Join(dir, randb32(max(len(originalname), 20)))
		//time.Sleep(time.Second * 2)
	}
	newname = filename
	err = nil
	return
}

//MixTime : change file date
func MixTime(filename string, count int) error {
	for i := 0; i < count; i++ {
		err := os.Chtimes(filename, time.Unix(mrand.Int63n(time.Now().Unix()), 0), time.Unix(mrand.Int63n(time.Now().Unix()), 0))
		if err != nil {
			return err
		}
	}
	return nil
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

func init() {
	mrand.Seed(time.Now().UTC().UnixNano())
	mrand.Seed(time.Now().UTC().UnixNano())
}

//Wipe7pass : (randfill-zerofill-randfill-complementfill-randfill-randfill-randfill)
func Wipe7pass(filename string) error {
	err := RandFill(filename, 1)
	if err != nil {
		return err
	}
	err = ZeroFill(filename, 1)
	if err != nil {
		return err
	}
	err = RandFill(filename, 3)
	if err != nil {
		return err
	}
	err = ComplementFill(filename)
	if err != nil {
		return err
	}
	filename, err = MixFileName(filename, 10)
	if err != nil {
		return err
	}
	err = MixTime(filename, 10)
	if err != nil {
		return err
	}
	filename, err = MixFileName(filename, 5)
	if err != nil {
		return err
	}
	err = MixTime(filename, 5)
	if err != nil {
		return err
	}
	return os.Remove(filename)
}

//Wipe8pass : (randfill-zerofill-randfill-complementfill-randfill-randfill-randfill-zerofill)
func Wipe8pass(filename string) error {
	err := RandFill(filename, 1)
	if err != nil {
		return err
	}
	err = ZeroFill(filename, 1)
	if err != nil {
		return err
	}
	err = RandFill(filename, 3)
	if err != nil {
		return err
	}
	err = ComplementFill(filename)
	if err != nil {
		return err
	}
	filename, err = MixFileName(filename, 10)
	if err != nil {
		return err
	}
	err = MixTime(filename, 10)
	if err != nil {
		return err
	}
	filename, err = MixFileName(filename, 5)
	if err != nil {
		return err
	}
	err = MixTime(filename, 5)
	if err != nil {
		return err
	}
	err = ZeroFill(filename, 1)
	if err != nil {
		return err
	}
	return os.Remove(filename)
}
