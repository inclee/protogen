package gen_dart

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

func formatCode(filename string) error {
	cmd := exec.Command("dartfmt", "-w", filename) // 格式化当前目录下的所有Dart文件
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
func Gen(protoDir, codeDir string) error {
	// 获取特定后缀名的文件列表
	files, err := getFilesWithExtension(protoDir, ".proto")
	if err != nil {
		return err
	}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	// 读取并打印每个文件的内容
	for _, file := range files {
		logger.Printf("read proto file %s", file)
		fname := utool.FileNameWithoutExt(file) + ".gen.dart"
		entityfile := filepath.Join(codeDir, "entity/"+fname)
		gener := map[string]func(string, string) (string, error){
			entityfile: parseEntity,
		}
		for fpath, gen := range gener {
			content, err := gen(protoDir, file)
			if err != nil {
				return err
			}
			if err = writeFile(fpath, content); err != nil {
				return err
			}
			if err = formatCode(fpath); err != nil {
				log.Printf("WARNNING: formate code failed %s", err.Error())
			}
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
		fname := utool.FileNameWithoutExt(file) + ".gen.dart"
		fpaths := []string{filepath.Join(codeDir, "service/"+fname),
			filepath.Join(codeDir, "entity/"+fname),
		}
		for _, fpath := range fpaths {
			os.Remove(fpath)
			logger.Printf("remove file %s", fpath)
		}
	}
	logger.Println("all file remove success. :-)")
	return nil
}
