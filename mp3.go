package musicmanager

/*
This code was adopted from https://github.com/xhenner/mp3-go with minor
changes by Lari Rasku.  It is released under the following license:

(The MIT License)

Copyright (c) 2014 Xavier Henner <xhenner@nerim.net>

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import (
	"io"
	"math"
)

// Some documentation :
// http://id3.org/mp3Frame
// http://www.codeproject.com/Articles/8295/MPEG-Audio-Frame-Header

// the number of frames to scan in fast mode

type mp3Infos struct {
	Version  string
	Layer    string
	Type     string
	Mode     string
	Bitrate  int
	Sampling int
	Size     int64
	Length   float64
	vbr      int
}

// constants for decoding frame header
var (
	mp3Version = [4]string{"2.5", "x", "2", "1"}
	mp3Layer   = [4]string{"r", "III", "II", "I"}
	mp3Bitrate = map[string][16]int{
		"1I":     {0, 32, 64, 96, 128, 160, 192, 224, 256, 288, 320, 352, 384, 416, 448},
		"1II":    {0, 32, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 384},
		"1III":   {0, 32, 40, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320},
		"2I":     {0, 32, 48, 56, 64, 80, 96, 112, 128, 144, 160, 176, 192, 224, 256},
		"2II":    {0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160},
		"2III":   {0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160},
		"2.5I":   {0, 32, 48, 56, 64, 80, 96, 112, 128, 144, 160, 176, 192, 224, 256},
		"2.5II":  {0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160},
		"2.5III": {0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160},
	}
	mp3Sampling = map[string][4]int{
		"1":   {44100, 48000, 32000, 0},
		"2":   {22050, 24000, 16000, 0},
		"2.5": {11025, 12000, 8000, 0},
	}
	mp3Channel      = [4]string{"Stereo", "Join Stereo", "Dual", "Mono"}
	frameLengthMult = map[string]int{
		"1I":     48,
		"1II":    144,
		"1III":   144,
		"2I":     24,
		"2II":    144,
		"2III":   72,
		"2.5I":   24,
		"2.5II":  72,
		"2.5III": 144,
	}
)

func mp3Fexamine(f io.ReadSeeker, slow bool) (*mp3Infos, error) {
	header := new(mp3Infos)
	var buf [10]byte
	var err error
	var bitrateSum int64
	var bitrateCount int64

	_, err = f.Read(buf[:10])

	if err != nil {
		return header, err
	}

	header.Size, err = f.Seek(0, 2)

	if err != nil {
		return header, err
	}

	pos, err := f.Seek(getId3Size(&buf), 0)

	if err != nil {
		return header, err
	}
	vbr := 0
	for {
		i, _ := f.Read(buf[:10])
		if i < 10 {
			break
		}
		pos += int64(i)
		// looking for the synchronization bits
		switch {
		case (buf[0] == 255) && (buf[1] >= 224):
			pos, _ = f.Seek(header.analyse(&buf, &vbr)-10, 1)
			bitrateSum += int64(header.Bitrate)
			bitrateCount++
			break
		case string(buf[:3]) == "TAG":
			pos, _ = f.Seek(128-10, 1) // id3v1 tag, bypass it
			break
		default:
			f.Seek(-9, 1) // looking for the next header
		}
	}

	if bitrateCount > 1 && header.Type == "VBR" {
		s := float64(bitrateSum / bitrateCount)
		diff := s
		for _, v := range mp3Bitrate[header.Version+header.Layer] {
			if math.Abs(float64(v)-s) < diff {
				header.Bitrate = v
				diff = math.Abs(float64(v) - s)
			}
		}
	}

	return header, nil
}

func (m *mp3Infos) analyse(buf *[10]byte, vbrCount *int) int64 {
	v := buf[1] & 24 >> 3
	l := buf[1] & 6 >> 1

	b := buf[2] & 240 >> 4
	s := buf[2] & 12 >> 2
	c := buf[3] & 192 >> 6

	// if the values are off, try 1 byte after
	if l == 0 || b == 15 || v == 1 || b == 0 || s == 3 {
		return 11
	}

	pad := int64(buf[2] & 2 >> 1)
	bitrate := mp3Bitrate[mp3Version[v]+mp3Layer[l]][b]

	//fmt.Println(m.Bitrate, bitrate)
	switch {
	case m.Type == "":
		m.Type = "CBR"

	case *vbrCount > 15:
		m.Type = "VBR"

	case bitrate != m.Bitrate:
		*vbrCount++
	}

	m.Bitrate = bitrate

	if m.Version == "" {
		m.Version = mp3Version[v]
		m.Layer = mp3Layer[l]
		m.Sampling = mp3Sampling[mp3Version[v]][s]
		m.Mode = mp3Channel[c]
	}

	samples := 1152
	switch {
	case m.Layer == "I":
		samples = 384
		break
	case m.Layer == "III" && m.Version != "1":
		samples = 576
		break
	}
	m.Length += float64(samples) / float64(m.Sampling)

	if l == 3 { // layer I
		return (int64(12*bitrate*1000/m.Sampling) + pad) * 4
	}

	return int64(144*bitrate*1000/m.Sampling) + pad
}

func getId3Size(buf *[10]byte) int64 {
	if string(buf[:3]) != "ID3" {
		return 0
	}
	var s int64

	// check if there is a footer. add 10 to the size
	if buf[5]&0x10 != 0 {
		s = 10
	}

	// cancel the last bit of each byte, and read the total
	for k := 6; k <= 9; k++ {
		s += int64(127&buf[k]) * (1 << (uint32(9-k) * 7))
	}

	// add the 10 octets of the header
	return s + 10
}
