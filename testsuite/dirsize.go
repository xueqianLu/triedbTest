package testsuite

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetDirSize(path string) (ByteSize, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return ByteSize(size), err
}

type ByteSize float64

const (
	B  ByteSize = 1
	KB          = B * 1024
	MB          = KB * 1024
	GB          = MB * 1024
	TB          = GB * 1024
	PB          = TB * 1024
)

func (b ByteSize) String() string {
	switch {
	case b >= PB:
		return fmt.Sprintf("%.2f PB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%.2f TB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2f GB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2f MB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2f KB", b/KB)
	default:
		return fmt.Sprintf("%.2f B", b)
	}
}
