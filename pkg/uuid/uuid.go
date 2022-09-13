package uuid

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type UUID struct {
	High int64
	Low  int64
}

func NewUUID(uuid string) (*UUID, error) {
	if uuid[8] != '-' || uuid[13] != '-' || uuid[18] != '-' || uuid[23] != '-' {
		return nil, fmt.Errorf("invalid UUID: %s", uuid)
	}

	data := make([]byte, 16)
	for i := 0; i < 4; i++ {
		data[15-i] = hexPairToChar(uuid[2*i], uuid[2*i+1])
	}
	data[11] = hexPairToChar(uuid[9], uuid[10])
	data[10] = hexPairToChar(uuid[11], uuid[12])
	data[9] = hexPairToChar(uuid[14], uuid[15])
	data[8] = hexPairToChar(uuid[16], uuid[17])
	data[7] = hexPairToChar(uuid[19], uuid[20])
	data[6] = hexPairToChar(uuid[21], uuid[22])

	for i := 10; i < 16; i++ {
		data[15-i] = hexPairToChar(uuid[2*i+4], uuid[2*i+5])
	}

	return &UUID{
		High: int64(binary.BigEndian.Uint64(data[:8])),
		Low:  int64(binary.BigEndian.Uint64(data[8:])),
	}, nil
}

func hexPairToChar(a uint8, b uint8) uint8 {
	convert := func(a uint8) uint8 {
		if a >= 97 {
			return a - 87
		} else {
			if a >= 65 {
				return a - 55
			} else {
				return a - 48
			}
		}
	}
	return convert(a)<<4 + convert(b)
}

var uuidMap = map[int64]struct{}{}

func (uuid *UUID) Squash() int64 {
	low := uuid.Low
	for {
		if _, ok := uuidMap[low]; ok {
			low++
		} else {
			break
		}
	}
	return low
}

func (uuid *UUID) HexString() string {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf, uint64(uuid.High))
	binary.BigEndian.PutUint64(buf[8:], uint64(uuid.Low))
	return hex.EncodeToString(buf)
}
