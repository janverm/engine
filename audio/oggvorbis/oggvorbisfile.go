// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oggvorbis

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"unsafe"
)

type VorbisInfo struct {
	Version        int
	Channels       int
	Rate           int
	BitrateUpper   int
	BitrateNominal int
	BitrateLower   int
}

type OVFile struct {
	OsHandle *os.File
	FsHandle fs.File
	OvReader *Reader
}

// Open opens an ogg vorbis file for decoding
// Returns an opaque pointer to the internal decode structure and an error
func Open(path string) (*OVFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Error:%s from Open", err)
	}
	var ovf OVFile
	ovf.OsHandle = f
	ovf.OvReader, err = NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("Error:%s from NewReader", err)
	}

	return &ovf, nil
}

func OpenEmbedded(path string, efs *embed.FS) (*OVFile, error) {
	f, err := efs.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Error:%s from OpenEmbedded", err)
	}
	var ovf OVFile
	ovf.FsHandle = f
	ovf.OvReader, err = NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("Error:%s from NewReader", err)
	}

	return &ovf, nil
}

// Clear clears the decoded buffers and closes the file
func Clear(f *OVFile) error {

	if f.OsHandle != nil {
		f.OsHandle.Close()
	}
	if f.FsHandle != nil {
		f.FsHandle.Close()
	}
	return nil
}

// Read decodes next data from the file updating the specified buffer contents and
// returns the number of bytes read, the number of current logical bitstream and an error
func Read(f *OVFile, buffer unsafe.Pointer, nrBytes int) (int, int, error) {
	nrSamples := nrBytes / 4
	temp := make([]float32, nrSamples)
	n, err := f.OvReader.Read(temp)
	if err != nil {
		log.Printf("ov.Read(): n = %d, err = %v", n, err)
	}
	bufferSlice := unsafe.Slice((*int16)(buffer), nrBytes/2)
	for i := 0; i < n; i = i + 1 {
		bufferSlice[i] = int16(temp[i] * float32(32768.))
	}
	return n * 2, 0, err
}

// Info updates the specified VorbisInfo structure with contains basic
// information about the audio in a vorbis stream
func Info(f *OVFile, link int, info *VorbisInfo) error {

	info.Version = 1
	info.Channels = f.OvReader.Channels()
	info.Rate = f.OvReader.SampleRate()
	info.BitrateUpper = f.OvReader.Bitrate().Maximum
	info.BitrateNominal = f.OvReader.Bitrate().Nominal
	info.BitrateLower = f.OvReader.Bitrate().Minimum
	return nil
}

// Seekable returns indication whether or not the bitstream is seekable
func Seekable(f *OVFile) bool {

	return true
}

// Seek seeks to the offset specified (in number pcm samples) within the physical bitstream.
// This function only works for seekable streams.
// Updates everything needed within the decoder, so you can immediately call Read()
// and get data from the newly seeked to position.
func PcmSeek(f *OVFile, pos int64) error {
	return f.OvReader.SetPosition(pos)
}

// PcmTotal returns the total number of pcm samples of the physical bitstream or a specified logical bit stream.
// To retrieve the total pcm samples for the entire physical bitstream, the 'link' parameter should be set to -1
func PcmTotal(f *OVFile, i int) (int64, error) {
	l := f.OvReader.Length()
	if l == 0 {
		return 0, fmt.Errorf("Error:Ogg Vorbis file has 0 samples.")
	}
	return l, nil
}

// TimeTotal returns the total time in seconds of the physical bitstream or a specified logical bitstream
// To retrieve the time total for the entire physical bitstream, 'i' should be set to -1.
func TimeTotal(f *OVFile, i int) (float64, error) {

	return -1, nil
}

// TimeTell returns the current decoding offset in seconds.
func TimeTell(f *OVFile) (float64, error) {

	return -1, nil
}
