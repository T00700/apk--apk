package editor

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Modifier interface {
	~string | ~uint32
}

// ModifyInfo 定义需要修改的新旧值
type ModifyInfo[T Modifier] struct {
	Old T
	New T
}

// adjustStringLength 调整字符串长度，使new的长度与old一致
// 如果new更长则截断，如果更短则用空格补齐
func adjustStringLength(old, new []byte) []byte {
	if len(new) > len(old) {
		// 如果new更长，截取到old的长度
		return new[:len(old)]
	} else if len(new) < len(old) {
		// 如果new更短，用空格补齐
		padding := make([]byte, len(old)-len(new)) // 需要补充的字节数
		return append(new, padding...)
	}
	return new
}

// ModifyAll 支持同时处理不同类型的修改
func ModifyAll(data []byte, modifications ...any) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty input data")
	}

	result := data
	for _, mod := range modifications {
		switch m := mod.(type) {
		case ModifyInfo[string]:
			oldBytes := changeString([]byte(m.Old))
			newBytes := changeString([]byte(m.New))
			newBytes = adjustStringLength(oldBytes, newBytes)
			result = bytes.Replace(result, oldBytes, newBytes, -1)
		case ModifyInfo[uint32]:
			old := changeUInt32(m.Old)
			new := changeUInt32(m.New)
			result = bytes.Replace(result, old, new, -1)
		default:
			return nil, fmt.Errorf("unsupported modification type: %T", mod)
		}
	}

	return result, nil
}

// changeString 将字符串转换为 UTF-16LE 格式
func changeString(data []byte) []byte {
	var ret = make([]byte, len(data)*2)
	for i, b := range data {
		ret[i*2] = b
		ret[i*2+1] = 0
	}
	return ret
}

// changeUInt32 将 uint32 转换为大端字节序
func changeUInt32(i uint32) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, i)
	return buf
}
