package gen_go

import (
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/inclee/protogen/internal/errors"
	"github.com/inclee/protogen/internal/utool"
)

func getFilesWithExtension(dir string, ext string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ext) {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func writeFile(filePath, content string) error {
	dir := filepath.Dir(filePath)

	// 创建所有不存在的目录
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func formatGoCode(filename string) error {
	// 读取文件内容
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// 格式化代码
	formattedContent, err := format.Source(content)
	if err != nil {
		return err
	}

	// 保存格式化后的内容到原文件
	err = ioutil.WriteFile(filename, formattedContent, 0644)
	if err != nil {
		return err
	}

	return nil
}
func Gen(protoDir, codeDir string) error {
	// 获取特定后缀名的文件列表
	files, err := getFilesWithExtension(protoDir, ".proto")
	if err != nil {
		return err
	}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	if err := writeValide(codeDir); err != nil {
		return err
	}
	// 读取并打印每个文件的内容
	for _, file := range files {
		logger.Printf("read proto file %s", file)
		fname := utool.FileNameWithoutExt(file) + ".gen.go"
		servicefile := filepath.Join(codeDir, "service/"+fname)
		entityfile := filepath.Join(codeDir, "entity/"+fname)
		handerfile := filepath.Join(codeDir, "handler/"+fname)
		gener := map[string]func(string, string) (string, error){
			servicefile: genService,
			entityfile:  genEntity,
			handerfile:  genHandler,
		}
		for fpath, gen := range gener {
			content, err := gen(protoDir, file)
			if err != nil {
				if err == errors.ErrNotNeedGen {
					continue
				}
				return err
			}
			if err = writeFile(fpath, content); err != nil {
				return err
			}
			formatGoCode(fpath)
			logger.Printf("generate file %s success", fpath)
		}
	}
	logger.Println("all file generate success. :-)")
	return nil
}

func Clean(protoDir, codeDir string) error {
	// 获取特定后缀名的文件列表
	files, err := getFilesWithExtension(protoDir, ".proto")
	if err != nil {
		return err
	}
	// 读取并打印每个文件的内容
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	for _, file := range files {
		fname := utool.FileNameWithoutExt(file) + ".gen.go"
		fpaths := []string{filepath.Join(codeDir, "service/"+fname),
			filepath.Join(codeDir, "entity/"+fname),
			filepath.Join(codeDir, "handler/"+fname),
		}
		for _, fpath := range fpaths {
			os.Remove(fpath)
			logger.Printf("remove file %s", fpath)
		}
	}
	logger.Println("all file remove success. :-)")
	return nil
}
