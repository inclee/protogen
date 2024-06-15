package parse

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/inclee/gokit/container/maps"
	"github.com/inclee/protogen/internal/entity"
	"github.com/inclee/protogen/internal/utool"
	"github.com/inclee/protogen/internal/utool/containers/slice"
)

func Parsefile(workdir string, fpath string) (compnent entity.CodeCompnent, err error) {
	compnent = entity.CodeCompnent{
		Messages: make(map[string]*entity.Message),
		Services: make(map[string]*entity.Service),
	}
	// imports []string, messages map[string]*entity.Message, services map[string]*entity.Service,
	file, err := os.Open(fpath)
	if err != nil {
		return
	}
	defer file.Close()

	// 创建消息和服务列表

	// 当前解析的消息和服务
	var currentMessage *entity.Message
	var currentService *entity.Service
	baseMessages := []string{}
	// 创建一个扫描器
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 忽略空行和注释行
		if len(line) == 0 || strings.HasPrefix(line, "//") {
			continue
		}
		if len(line) == 0 || strings.HasPrefix(line, "module") {
			compnent.Module = strings.Fields(line)[1]
			continue
		}
		if strings.HasPrefix(line, "import") {
			importPath := strings.Fields(line)[1]
			compnent.Imports = append(compnent.Imports, utool.FileNameWithoutExt(importPath)+".gen.dart")
			parent, err := Parsefile(workdir, filepath.Join(workdir, strings.Fields(line)[1]))
			if err != nil {
				return compnent, err
			}
			for msg := range parent.Messages {
				baseMessages = append(baseMessages, msg)
			}
			maps.Merge(compnent.Messages, parent.Messages, false)
			maps.Merge(compnent.Services, parent.Services, false)
		}
		// 解析消息定义
		if strings.HasPrefix(line, "message") && (strings.HasSuffix(line, "{")) {
			name := strings.Replace(strings.TrimSuffix(line, "{"), "message", "", -1)
			currentMessage = &entity.Message{Name: strings.TrimSpace(name)}
		} else if line == "}" && currentMessage != nil {
			compnent.Messages[currentMessage.Name] = currentMessage
			currentMessage = nil
		} else if currentMessage != nil {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				currentMessage.Fields = append(currentMessage.Fields, &entity.Field{Name: strings.TrimSpace(parts[0]), Type: parts[1], Validate: slice.Get(parts, 2, "")})
			} else {
				parent := compnent.Messages[strings.TrimSpace(parts[0])]
				if parent == nil {
					panic(fmt.Sprintf("not found message:%s ", parts[0]))
				}
				currentMessage.Parent = parent
			}
		}
		// 解析服务定义
		if strings.HasPrefix(line, "service") {
			name := strings.TrimSuffix(strings.Fields(line)[1], " {")
			currentService = &entity.Service{Name: name}
		} else if line == "}" && currentService != nil {
			compnent.Services[currentService.Name] = currentService
			currentService = nil
		} else if currentService != nil && strings.Contains(line, "(") && strings.Contains(line, ")") {
			parts := strings.Fields(line)
			method := parts[0]
			handlerName := parts[1]
			request := strings.Trim(parts[2], "()")
			response := strings.Trim(parts[3], "()")
			compnent.Messages[request].AttachTag(tagMap(strings.ToLower(method)))
			currentService.Handlers = append(currentService.Handlers, entity.HandlerFunc{Method: strings.ToUpper(method), Name: handlerName, Request: request, Response: response})
		}
	}
	for _, base := range baseMessages {
		delete(compnent.Messages, base)
	}
	err = scanner.Err()
	return
}

func tagMap(method string) string {
	m := map[string]string{"get": "form", "post": "json"}
	return m[method]
}
