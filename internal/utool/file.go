package utool

import "path/filepath"

func FileNameWithoutExt(path string) string {
	filename := filepath.Base(path)
	return filename[:len(filename)-len(filepath.Ext(filename))]
}
