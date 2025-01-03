package editor

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type StringOrUint32 interface {
	~string | ~uint32
}

// ModifyInfo 定义需要修改的新旧值
type ModifyInfo[T StringOrUint32] struct {
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

// Modify 修改二进制 manifest 文件中的内容
// 返回修改后的数据和可能的错误
func Modify(data []byte, modifyInfos ...any) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty input data")
	}

	result := data
	for _, info := range modifyInfos {
		switch v := info.(type) {
		case ModifyInfo[string]:
			oldBytes := changeString([]byte(v.Old))
			newBytes := changeString([]byte(v.New))
			newBytes = adjustStringLength(oldBytes, newBytes)
			result = bytes.Replace(result, oldBytes, newBytes, -1)

		case ModifyInfo[uint32]:
			new := changeUInt32(v.New)
			old := changeUInt32(v.Old)
			result = bytes.Replace(result, old, new, -1)

		default:
			return nil, fmt.Errorf("unsupported type: %T", info)
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
