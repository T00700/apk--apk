package main

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	defaultAlignment = 4
	pageAlignment    = 4096
	zipHeaderSize    = 30 // ZIP 本地文件头的固定大小
)

// getAlignment 根据文件类型决定对齐方式
func getAlignment(pageAlignSharedLibs bool, defaultAlignment int, fileName string) int {
	if !pageAlignSharedLibs {
		return defaultAlignment
	}

	if strings.HasSuffix(fileName, ".so") {
		return pageAlignment
	}

	return defaultAlignment
}

// fileNeedsAlignment 判断文件是否需要对齐
func fileNeedsAlignment(file *zip.File) bool {
	// 压缩文件不需要对齐
	if file.Method == zip.Deflate {
		return false
	}
	return true
}

// createAlignedExtra 创建对齐的扩展字段
func createAlignedExtra(oldExtra []byte, padding int) []byte {
	// 如果填充小于等于0，直接返回原扩展字段
	if padding <= 0 {
		return oldExtra
	}

	// 创建新的扩展字段
	newExtra := make([]byte, len(oldExtra)+padding)
	copy(newExtra, oldExtra)

	// 添加填充字节
	for i := len(oldExtra); i < len(newExtra); i++ {
		newExtra[i] = 0
	}

	return newExtra
}

// 添加一个结构体来存储文件信息和偏移量
type fileInfo struct {
	file    *zip.File
	content []byte
	needPad bool
	alignTo int
	offset  uint32 // 添加offset字段来存储文件的偏移量
}

// copyAndAlign 将输入zip文件的内容复制到输出zip文件，并进行对齐
func copyAndAlign(inFile, outFile string, alignment int, pageAlignSharedLibs bool) error {
	reader, err := zip.OpenReader(inFile)
	if err != nil {
		return fmt.Errorf("无法打开输入文件: %v", err)
	}
	defer reader.Close()

	// 创建临时缓冲区
	buf := new(bytes.Buffer)

	// 收集所有文件信息
	var files []fileInfo

	// 读取所有文件内容
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}

		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return err
		}

		needPad := fileNeedsAlignment(file)
		alignTo := getAlignment(pageAlignSharedLibs, alignment, file.Name)

		files = append(files, fileInfo{
			file:    file,
			content: content,
			needPad: needPad,
			alignTo: alignTo,
		})
	}

	// 按目录路径排序
	sort.Slice(files, func(i, j int) bool {
		dirI := filepath.Dir(files[i].file.Name)
		dirJ := filepath.Dir(files[j].file.Name)
		if dirI != dirJ {
			return dirI < dirJ
		}
		return files[i].file.Name < files[j].file.Name
	})

	// 写入文件数据
	currentOffset := int64(0)
	for i := range files {
		// 保存文件的偏移量
		files[i].offset = uint32(currentOffset)

		// 写入本地文件头
		header := files[i].file.FileHeader

		// 计算文件头大小
		headerSize := int64(30 + len(header.Name) + len(header.Extra)) // 30 是固定头部大小

		// 如果需要对齐，计算并添加填充
		if files[i].needPad {
			dataOffset := currentOffset + headerSize
			padding := (int64(files[i].alignTo) - (dataOffset % int64(files[i].alignTo))) % int64(files[i].alignTo)
			if padding > 0 {
				header.Extra = append(header.Extra, make([]byte, padding)...)
				headerSize += padding
			}
		}

		// 写入本地文件头
		binary.Write(buf, binary.LittleEndian, uint32(0x04034b50)) // 本地文件头签名
		binary.Write(buf, binary.LittleEndian, uint16(20))         // 版本
		binary.Write(buf, binary.LittleEndian, uint16(0))          // 通用标志位
		binary.Write(buf, binary.LittleEndian, uint16(files[i].file.Method))
		binary.Write(buf, binary.LittleEndian, uint16(header.Modified.Hour()<<11|header.Modified.Minute()<<5|header.Modified.Second()/2))
		binary.Write(buf, binary.LittleEndian, uint16((header.Modified.Year()-1980)<<9|int(header.Modified.Month())<<5|header.Modified.Day()))
		binary.Write(buf, binary.LittleEndian, header.CRC32)
		binary.Write(buf, binary.LittleEndian, header.CompressedSize64)
		binary.Write(buf, binary.LittleEndian, header.UncompressedSize64)
		binary.Write(buf, binary.LittleEndian, uint16(len(header.Name)))
		binary.Write(buf, binary.LittleEndian, uint16(len(header.Extra)))

		// 写入文件名和扩展字段
		buf.WriteString(header.Name)
		buf.Write(header.Extra)

		// 写入文件内容
		buf.Write(files[i].content)

		// 更新偏移量
		currentOffset += headerSize + int64(len(files[i].content))
	}

	// 写入中央目录
	centralDirOffset := currentOffset
	for _, f := range files {
		header := f.file.FileHeader
		binary.Write(buf, binary.LittleEndian, uint32(0x02014b50)) // 中央目录头签名
		binary.Write(buf, binary.LittleEndian, uint16(20))         // 版本
		binary.Write(buf, binary.LittleEndian, uint16(20))         // 最小版本
		binary.Write(buf, binary.LittleEndian, uint16(0))          // 通用标志位
		binary.Write(buf, binary.LittleEndian, uint16(f.file.Method))
		binary.Write(buf, binary.LittleEndian, uint16(header.Modified.Hour()<<11|header.Modified.Minute()<<5|header.Modified.Second()/2))
		binary.Write(buf, binary.LittleEndian, uint16((header.Modified.Year()-1980)<<9|int(header.Modified.Month())<<5|header.Modified.Day()))
		binary.Write(buf, binary.LittleEndian, header.CRC32)
		binary.Write(buf, binary.LittleEndian, header.CompressedSize64)
		binary.Write(buf, binary.LittleEndian, header.UncompressedSize64)
		binary.Write(buf, binary.LittleEndian, uint16(len(header.Name)))
		binary.Write(buf, binary.LittleEndian, uint16(len(header.Extra)))
		binary.Write(buf, binary.LittleEndian, uint16(0)) // 文件注释长度
		binary.Write(buf, binary.LittleEndian, uint16(0)) // 磁盘开始号
		binary.Write(buf, binary.LittleEndian, uint16(0)) // 内部文件属性
		binary.Write(buf, binary.LittleEndian, uint32(0)) // 外部文件属性
		binary.Write(buf, binary.LittleEndian, f.offset)  // 使用保存的偏移量
		buf.WriteString(header.Name)
		buf.Write(header.Extra)
		currentOffset += int64(46 + len(header.Name) + len(header.Extra))
	}

	// 写入中央目录结束标记
	binary.Write(buf, binary.LittleEndian, uint32(0x06054b50))                     // 结束标记签名
	binary.Write(buf, binary.LittleEndian, uint16(0))                              // 当前磁盘编号
	binary.Write(buf, binary.LittleEndian, uint16(0))                              // 中央目录开始磁盘编号
	binary.Write(buf, binary.LittleEndian, uint16(len(files)))                     // 本磁盘上的条目数
	binary.Write(buf, binary.LittleEndian, uint16(len(files)))                     // 中央目录条目总数
	binary.Write(buf, binary.LittleEndian, uint32(currentOffset-centralDirOffset)) // 中央目录大小
	binary.Write(buf, binary.LittleEndian, uint32(centralDirOffset))               // 中央目录偏移
	binary.Write(buf, binary.LittleEndian, uint16(0))                              // 注释长度

	// 写入最终文件
	return os.WriteFile(outFile, buf.Bytes(), 0644)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("用法: zipalign <input.zip> <output.zip>")
		os.Exit(1)
	}

	inFile := os.Args[1]
	outFile := os.Args[2]

	// 检查输入文件是否存在
	if _, err := os.Stat(inFile); os.IsNotExist(err) {
		fmt.Printf("错误: 输入文件 '%s' 不存在\n", inFile)
		os.Exit(1)
	}

	// 检查输入输出文件是否相同
	if inFile == outFile {
		fmt.Println("错误: 输入和输出文件不能相同")
		os.Exit(1)
	}

	// 使用默认4字节对齐，启用.so文件页面对齐
	err := copyAndAlign(inFile, outFile, defaultAlignment, true)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("成功: '%s' 已对齐并保存到 '%s'\n", inFile, outFile)
}
