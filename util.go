package aliyundrive_open

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

// randomString 生成随机字符串
func randomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return string(b)
}

// 合并自定义的字符串类型
func joinCustomString[T fmt.Stringer](items []T, separator string) string {
	switch len(items) {
	case 0:
		return ""
	case 1:
		return fmt.Sprint(items[0])
	default:
		var b strings.Builder
		b.WriteString(fmt.Sprint(items[0]))
		for _, s := range items[1:] {
			b.WriteString(separator)
			b.WriteString(fmt.Sprint(s))
		}
		return b.String()
	}
}

// SplitFile 处理文件分片信息
func SplitFile(file *os.File) (partInfoList []FileUpdatePartInfo, err error) {

	stat, err := file.Stat()
	if err != nil {
		return partInfoList, err
	}

	var partInfo = FileUpdatePartInfo{}
	if stat.Size() <= DefaultPartSize {
		partInfo.PartNumber = 1
		partInfo.ParallelSha1Ctx.PartOffset = 0
		partInfo.ParallelSha1Ctx.PartSize = stat.Size()
		partInfoList = append(partInfoList, partInfo)
		return partInfoList, nil
	}

	var n = stat.Size() / DefaultPartSize
	var otherSize = stat.Size() % DefaultPartSize

	for i := int64(0); i < n; i++ {
		partInfo.PartNumber = i + 1
		partInfo.ParallelSha1Ctx.PartOffset = i * DefaultPartSize
		partInfo.ParallelSha1Ctx.PartSize = DefaultPartSize
		if i == n-1 {
			partInfo.ParallelSha1Ctx.PartSize = DefaultPartSize + otherSize
		}
		partInfoList = append(partInfoList, partInfo)
	}

	return partInfoList, nil
}

// SplitFile 处理文件分片信息
func SplitFileC(file *os.File) (partInfoList []FileUpdatePartInfo, err error) {

	stat, err := file.Stat()
	if err != nil {
		return partInfoList, err
	}

	var partInfo = FileUpdatePartInfo{}

	if stat.Size() <= DefaultPartSize {
		partInfo.PartNumber = 1
		partInfo.ParallelSha1Ctx.PartOffset = 0
		partInfo.ParallelSha1Ctx.PartSize = stat.Size()
		partInfoList = append(partInfoList, partInfo)
		return partInfoList, nil
	}

	otherSize := stat.Size() % DefaultPartSize
	var n = stat.Size() / DefaultPartSize
	var h []uint32

	for i := int64(0); i < n; i++ {
		partInfo.PartNumber = i + 1
		partInfo.ParallelSha1Ctx.PartOffset = i * DefaultPartSize
		partInfo.ParallelSha1Ctx.PartSize = DefaultPartSize
		if i == n-1 {
			partInfo.ParallelSha1Ctx.PartSize = DefaultPartSize + otherSize
		}
		if i > 0 {
			partInfo.ParallelSha1Ctx.H = h
		}

		data, err := io.ReadAll(io.LimitReader(file, partInfo.ParallelSha1Ctx.PartSize))
		if err != nil {
			return partInfoList, err
		}

		hasher := sha1.New()
		hash := hasher.Sum(data)

		h = make([]uint32, 5)
		for i := 0; i < 5; i++ {
			h[i] = binary.BigEndian.Uint32(hash[i*4 : (i+1)*4])
		}

		partInfoList = append(partInfoList, partInfo)
	}

	return partInfoList, nil
}
