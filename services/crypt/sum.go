package crypt

import (
	"encoding/binary"
	"strconv"
	"strings"
)

// Luhn validates a number using the Luhn algorithm.
func (s *Service) Luhn(payload string) bool {
	payload = strings.ReplaceAll(payload, " ", "")
	sum := 0
	isSecond := false
	for i := len(payload) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(payload[i]))
		if err != nil {
			return false // Contains non-digit
		}

		if isSecond {
			digit = digit * 2
			if digit > 9 {
				digit = digit - 9
			}
		}

		sum += digit
		isSecond = !isSecond
	}
	return sum%10 == 0
}

// Fletcher16 computes the Fletcher-16 checksum.
func (s *Service) Fletcher16(payload string) uint16 {
	data := []byte(payload)
	var sum1, sum2 uint16
	for _, b := range data {
		sum1 = (sum1 + uint16(b)) % 255
		sum2 = (sum2 + sum1) % 255
	}
	return (sum2 << 8) | sum1
}

// Fletcher32 computes the Fletcher-32 checksum.
func (s *Service) Fletcher32(payload string) uint32 {
	data := []byte(payload)
	// Pad with 0 to make it even length for uint16 conversion
	if len(data)%2 != 0 {
		data = append(data, 0)
	}

	var sum1, sum2 uint32
	for i := 0; i < len(data); i += 2 {
		val := binary.LittleEndian.Uint16(data[i : i+2])
		sum1 = (sum1 + uint32(val)) % 65535
		sum2 = (sum2 + sum1) % 65535
	}
	return (sum2 << 16) | sum1
}

// Fletcher64 computes the Fletcher-64 checksum.
func (s *Service) Fletcher64(payload string) uint64 {
	data := []byte(payload)
	// Pad to multiple of 4
	if len(data)%4 != 0 {
		padding := 4 - (len(data) % 4)
		data = append(data, make([]byte, padding)...)
	}

	var sum1, sum2 uint64
	for i := 0; i < len(data); i += 4 {
		val := binary.LittleEndian.Uint32(data[i : i+4])
		sum1 = (sum1 + uint64(val)) % 4294967295
		sum2 = (sum2 + sum1) % 4294967295
	}
	return (sum2 << 32) | sum1
}
