package hashcash

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

const (
	maxIterations    int    = 1 << 20        // Max iterations to find a solution
	bytesToRead      int    = 8              // Bytes to read for random token
	bitsPerHexChar   int    = 4              // Each hex character takes 4 bits
	zero             rune   = 48             // ASCII code for number zero
	hashcashV1Length int    = 7              // Number of items in a V1 hashcash header
	timeFormat       string = "060102150405" // YYMMDDhhmmss

	defaultTokenSize = 16
)

type Config struct {
	Bits    int
	Expired time.Time
	Future  time.Time
}

var DefaultConfig = &Config{
	Bits:    20,
	Future:  time.Now().AddDate(0, 0, 2),
	Expired: time.Now().AddDate(0, 0, -30),
}

type Hashcash struct {
	version   int
	bits      int
	created   time.Time
	resource  []byte
	extension string
	rand      string
	counter   int
	expired   time.Time
	future    time.Time
}

func (h *Hashcash) Compute() (string, error) {
	var (
		wantZeros = h.bits / bitsPerHexChar
		header    = h.createHeader()
		hash      = sha1Hash(header)
	)
	for !acceptableHeader(hash, zero, wantZeros) {
		h.counter++
		header = h.createHeader()
		hash = sha1Hash(header)
		if h.counter >= maxIterations {
			return "", ErrSolutionFail
		}
	}
	return header, nil
}

func (h *Hashcash) Verify(header string) (bool, error) {
	vals := strings.Split(header, ":")
	if len(vals) != hashcashV1Length {
		return false, ErrInvalidHeader
	}
	var (
		hash      = sha1Hash(header)
		wantZeros = h.bits / bitsPerHexChar
	)
	if !acceptableHeader(hash, zero, wantZeros) {
		return false, ErrNoCollision
	}
	created, err := parseHashcashTime(vals[2])
	if err != nil {
		return false, err
	}
	if created.After(h.future) || created.Before(h.expired) {
		return false, ErrTimestamp
	}

	return true, nil
}

func New(res []byte) (*Hashcash, error) {
	if res == nil {
		return nil, ErrResourceEmpty
	}

	bytes, err := randomBytes(bytesToRead)
	if err != nil {
		return nil, err
	}

	return &Hashcash{
		version:   1,
		bits:      DefaultConfig.Bits,
		created:   time.Now(),
		resource:  res,
		extension: "",
		rand:      base64EncodeBytes(bytes),
		counter:   1,
		expired:   DefaultConfig.Expired,
		future:    DefaultConfig.Future,
	}, nil
}

func acceptableHeader(hash string, char rune, n int) bool {
	for _, val := range hash[:n] {
		if val != char {
			return false
		}
	}
	return true
}

func (h *Hashcash) createHeader() string {
	return fmt.Sprintf("%d:%d:%s:%s:%s:%s:%s", h.version,
		h.bits,
		h.created.Format(timeFormat),
		h.resource,
		h.extension,
		h.rand,
		base64EncodeInt(h.counter))
}

func parseHashcashTime(msgTime string) (date time.Time, err error) {
	// In a hashcash header the date parts year, month and day are mandatory but
	// hours, minutes and seconds are optional. So a valid date can be in format:
	//
	// - YYMMDD
	// - YYMMDDhhmm
	// - YYMMDDhhmmss
	switch len(msgTime) {
	case 6:
		f := timeFormat[:6]
		date, err = time.Parse(f, msgTime)
	case 10:
		f := timeFormat[:10]
		date, err = time.Parse(f, msgTime)
	case 12:
		f := timeFormat[:12]
		date, err = time.Parse(f, msgTime)
	}
	return date, err
}

func NewResourceToken(targetBits uint64) []byte {
	buf := make([]byte, defaultTokenSize)
	target := uint64(1) << (64 - targetBits)

	binary.BigEndian.PutUint64(buf[:bytesToRead], target)
	_, _ = rand.Read(buf[bytesToRead:])

	return buf
}
