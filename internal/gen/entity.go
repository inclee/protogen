package gen

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"text/template"
)

// 定义字段结构体
type Field struct {
	Name string
	Type string
}

// 定义消息结构体
type Message struct {
	Name   string
	Fields []Field
	Method string
}

// 生成代码的模板
const entiryTemplate = `
package entity

/*GENERATE BY protogen(https://github.com/inclee/protogen); PLEASE DON'T EDIT IT */

{{range .Messages}}
// {{.Name}} 请求结构体
{{- $Method := .Method }}
type {{.Name}} struct {
	{{- range .Fields }} 
		{{- if .Type }}
			{{- if eq $Method "GET" }}
				{{ .Name }} {{ .Type }} ` + "`form:\"{{ .Name | camelToSnakeCase }}\"`" + `
			{{- else }}
				{{ .Name }} {{ .Type }} ` + "`json:\"{{ .Name | camelToSnakeCase }}\"`" + `
			{{- end }}
		{{- else }}
			{{ .Name }} 
		{{- end }}
	{{- end }}
}
{{end}}

`

func parseEntity(filepath string) (string, error) {
	// 打开协议文件
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 创建消息和服务列表
	var messages []Message

	// 当前解析的消息和服务
	var currentMessage *Message
	var currentPackage string
	var currentService *Service
	var getMessages = map[string]bool{}

	// 创建一个扫描器
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 忽略空行和注释行
		if len(line) == 0 || strings.HasPrefix(line, "//") {
			continue
		}
		// 解析消息定义
		if strings.HasPrefix(line, "message") && (strings.HasSuffix(line, "{")) {
			name := strings.Replace(strings.TrimSuffix(line, "{"), "message", "", -1)
			currentMessage = &Message{Name: name}
		} else if line == "}" && currentMessage != nil {
			messages = append(messages, *currentMessage)
			currentMessage = nil
		} else if currentMessage != nil {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				currentMessage.Fields = append(currentMessage.Fields, Field{Name: parts[0], Type: parts[1]})
			} else {
				currentMessage.Fields = append(currentMessage.Fields, Field{Name: parts[0], Type: ""})
			}
		}

		if strings.HasPrefix(line, "service") {
			name := strings.TrimSuffix(strings.Fields(line)[1], " {")
			currentService = &Service{Name: name}
		} else if line == "}" && currentService != nil {
			currentService = nil
		} else if currentService != nil && strings.Contains(line, "(") && strings.Contains(line, ")") {
			parts := strings.Fields(line)
			method := parts[0]
			request := strings.Trim(parts[2], "()")
			if strings.ToLower(method) == "get" {
				getMessages[request] = true
			}

		}
	}

	// 检查是否有扫描错误
	if err := scanner.Err(); err != nil {
		return "", err
	}
	for idx, message := range messages {
		if getMessages[strings.TrimSpace(message.Name)] {
			message.Method = "GET"
			messages[idx] = message
		}
	}
	// 准备模板数据
	data := struct {
		Package  string
		Messages []Message
	}{
		Package:  currentPackage,
		Messages: messages,
	}

	// 创建模板并解析
	tmpl := template.Must(template.New("code").Funcs(template.FuncMap{"lowerFirst": lowerFirst, "camelToSnakeCase": camelToSnakeCase}).Parse(entiryTemplate))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	// 将生成的代码输出为字符串
	return buf.String(), nil
}
