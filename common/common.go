package common

import (
	"image"
	"image/draw"
)

// MaxEncodeSize given an image will find how many bytes can be stored in that image using least significant bit encoding
// ((width * height * 3) / 8 ) - 4
// The result must be at least 4,
func MaxEncodeSize(img image.Image) uint32 {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	eval := ((width * height * 3) / 8) - 4
	if eval < 4 {
		eval = 0
	}
	return uint32(eval)
}

// setBitInByte sets a specific bit in a byte to a given value and returns the new byte
func SetBitInByte(b byte, indexInByte uint32, bit byte) byte {
	var mask byte = 0x80
	mask = mask >> uint(indexInByte)

	if bit == 0 {
		mask = ^mask
		b = b & mask
	} else if bit == 1 {
		b = b | mask
	}
	return b
}

// getNextBitFromString each call will return the next subsequent bit in the string
func GetNextBitFromString(byteArray []byte, ch chan byte) {

	var offsetInBytes int
	var offsetInBitsIntoByte int
	var choiceByte byte

	lenOfString := len(byteArray)

	for {
		if offsetInBytes >= lenOfString {
			close(ch)
			return
		}

		choiceByte = byteArray[offsetInBytes]
		ch <- GetBitFromByte(choiceByte, offsetInBitsIntoByte)

		offsetInBitsIntoByte++

		if offsetInBitsIntoByte >= 8 {
			offsetInBitsIntoByte = 0
			offsetInBytes++
		}
	}
}

// getLSB given a byte, will return the least significant bit of that byte
func GetLSB(b byte) byte {
	if b%2 == 0 {
		return 0
	}
	return 1
}

// setLSB given a byte will set that byte's least significant bit to a given value (where true is 1 and false is 0)
func SetLSB(b *byte, bit byte) {
	if bit == 1 {
		*b = *b | 1
	} else if bit == 0 {
		var mask byte = 0xFE
		*b = *b & mask
	}
}

// getBitFromByte given a bit will return a bit from that byte
func GetBitFromByte(b byte, indexInByte int) byte {
	b = b << uint(indexInByte)
	var mask byte = 0x80

	var bit = mask & b

	if bit == 128 {
		return 1
	}
	return 0
}

// combineToInt given four bytes, will return the 32 bit unsigned integer which is the composition of those four bytes (one is MSB)
func CombineToInt(one, two, three, four byte) (ret uint32) {
	ret = uint32(one)
	ret = ret << 8
	ret = ret | uint32(two)
	ret = ret << 8
	ret = ret | uint32(three)
	ret = ret << 8
	ret = ret | uint32(four)
	return
}

// splitToBytes given an unsigned integer, will split this integer into its four bytes
func SplitToBytes(x uint32) (one, two, three, four byte) {
	one = byte(x >> 24)
	var mask uint32 = 255

	two = byte((x >> 16) & mask)
	three = byte((x >> 8) & mask)
	four = byte(x & mask)
	return
}

// imageToNRGBA converts image.Image to image.NRGBA
func ImageToNRGBA(src image.Image) *image.NRGBA {
	bounds := src.Bounds()

	var m *image.NRGBA
	var width, height int

	width = bounds.Dx()
	height = bounds.Dy()

	m = image.NewNRGBA(image.Rect(0, 0, width, height))

	draw.Draw(m, m.Bounds(), src, bounds.Min, draw.Src)
	return m
}
