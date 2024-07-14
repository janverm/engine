// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oggvorbis

import (
	"fmt"
	"io/fs"
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
	BitrateWindow  int
}

type OVFile struct {
	OsHandle *os.File
	FsHandle *fs.File
}

// Fopen opens an ogg vorbis file for decoding
// Returns an opaque pointer to the internal decode structure and an error
func Open(path string) (*OVFile, error) {
	return nil, fmt.Errorf("Error:%s from Fopen", 1)
}

// Clear clears the decoded buffers and closes the file
func Clear(f *OVFile) error {

	return nil
}

// Read decodes next data from the file updating the specified buffer contents and
// returns the number of bytes read, the number of current logical bitstream and an error
func Read(f *OVFile, buffer unsafe.Pointer, length int, bigendianp bool, word int, sgned bool) (int, int, error) {

	return 10, 4, nil
}

// Info updates the specified VorbisInfo structure with contains basic
// information about the audio in a vorbis stream
func Info(f *OVFile, link int, info *VorbisInfo) error {

	/*
		info.Version = int(vi.version)
		info.Channels = int(vi.channels)
		info.Rate = int(vi.rate)
		info.BitrateUpper = int(vi.bitrate_upper)
		info.BitrateNominal = int(vi.bitrate_nominal)
		info.BitrateLower = int(vi.bitrate_lower)
		info.BitrateWindow = int(vi.bitrate_window)
	*/
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

	return nil
}

// PcmTotal returns the total number of pcm samples of the physical bitstream or a specified logical bit stream.
// To retrieve the total pcm samples for the entire physical bitstream, the 'link' parameter should be set to -1
func PcmTotal(f *OVFile, i int) (int64, error) {

	return 10, nil
}

// TimeTotal returns the total time in seconds of the physical bitstream or a specified logical bitstream
// To retrieve the time total for the entire physical bitstream, 'i' should be set to -1.
func TimeTotal(f *OVFile, i int) (float64, error) {

	return 10, nil
}

// TimeTell returns the current decoding offset in seconds.
func TimeTell(f *OVFile) (float64, error) {

	return 10, nil
}
