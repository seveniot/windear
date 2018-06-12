package codec

import (
	"bytes"
	"github.com/SevenIOT/windear/ex"
)

/**
 *
 * @author: schbook
 * @email: seekerxu@163.com
 * @date: 2018/6/5
 *
 */

func DecodeUtf8String(buff []byte) (res string, content []byte, err error) {
	len := len(buff)

	if len < 2 {
		err = ex.IllegalPayloadData
		return
	}

	l := uint16(buff[0])<<8 + uint16(buff[1])

	if len-2 < int(l) {
		err = ex.IllegalPayloadData

		return
	}

	reLen := uint32(2 + l)
	res = string(bytes.TrimSpace(buff[2:reLen])) //not include 0x00
	content = buff[reLen:]

	return
}

func EncodeUtf8String(content string) []byte {
	buffer := []byte(content)

	l := len(buffer)

	return append([]byte{uint8(l >> 8), uint8(l & 0xFF)}, buffer...)
}
