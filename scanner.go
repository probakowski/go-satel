package satel

import "bytes"

func scan(data []byte, _ bool) (advance int, token []byte, err error) {
	i := 0
	for ; i < len(data) && data[i] == 0xFE; i++ {
	}
	if i > 0 {
		data = data[i:]
	}
	startIndex := bytes.Index(data, []byte{0xFE, 0xFE})
	index := bytes.Index(data, []byte{0xFE, 0x0D})
	if startIndex > 0 && (index < 0 || startIndex < index) {
		return i + startIndex + 2, nil, nil
	}
	if index > 0 {
		return i + index + 2, data[:index], nil
	}
	return 0, nil, nil
}
