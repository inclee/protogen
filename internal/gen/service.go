package gen

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"text/template"
	"unicode"
)

// 定义服务处理函数结构体
type HandlerFunc struct {
	Method   string
	Name     string
	Request  string
	Response string
}

// 定义服务结构体
type Service struct {
	Name     string
	Handlers []HandlerFunc
}

// 生成代码的模板
const codeTemplate = `
package service

/*GENERATE BY protogen(https://github.com/inclee/protogen); PLEASE DON'T EDIT IT */

import (
	"context"
	"gorm.io/gorm"
	"{{ .Module }}/internal/entity"
	"{{ .Module }}/internal/repository/query"
)

{{range .Services}}
// {{.Name}} 服务接口
type I{{.Name}}Service interface {
	{{- range .Handlers }}
	{{ .Name }}(context.Context, *entity.{{ .Request }}) (*entity.{{ .Response }}, error)
	{{- end }}
}

// 实现 {{.Name}} 服务
type {{ .Name }}Service struct{
	db *gorm.DB
	q  *query.Query
}

func New{{ .Name }}Service(db *gorm.DB) *{{ .Name }}Service{
	return &{{ .Name }}Service{
		db: db,
		q:  query.Use(db),
	}
}

type default{{ .Name }}Service struct{
}

{{- $serviceName := .Name }}
{{- range .Handlers }}
func (s *default{{ $serviceName }}Service) {{ .Name }}(ctx context.Context, req *entity.{{ .Request }}) (rsp *entity.{{ .Response }}, err error) {
	return &entity.{{ .Response }}{}, nil
}
{{- end }}
{{end}}
`

func parseService(filepath string) (string, error) {
	// 打开协议文件
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 创建消息和服务列表
	var services []Service

	// 当前解析的消息和服务
	var currentService *Service
	var currentPackage string
	var currentModule string

	// 创建一个扫描器
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 忽略空行和注释行
		if len(line) == 0 || strings.HasPrefix(line, "//") {
			continue
		}
		if strings.HasPrefix(line, "package") {
			currentPackage = strings.Fields(line)[1]
			continue
		}
		if strings.HasPrefix(line, "module") {
			currentModule = strings.Fields(line)[1]
			continue
		}

		// 解析服务定义
		if strings.HasPrefix(line, "service") {
			name := strings.TrimSuffix(strings.Fields(line)[1], " {")
			currentService = &Service{Name: name}
		} else if line == "}" && currentService != nil {
			services = append(services, *currentService)
			currentService = nil
		} else if currentService != nil && strings.Contains(line, "(") && strings.Contains(line, ")") {
			parts := strings.Fields(line)
			handlerName := parts[1]
			request := strings.Trim(parts[2], "()")
			response := strings.Trim(parts[3], "()")
			currentService.Handlers = append(currentService.Handlers, HandlerFunc{Name: handlerName, Request: request, Response: response})
		}
	}

	// 检查是否有扫描错误
	if err := scanner.Err(); err != nil {
		return "", err
	}

	// 准备模板数据
	data := struct {
		Package  string
		Services []Service
		Module   string
	}{
		Package:  currentPackage,
		Module:   currentModule,
		Services: services,
	}

	// 创建模板并解析
	tmpl := template.Must(template.New("code").Funcs(template.FuncMap{"lower": lower, "lowerFirst": lowerFirst, "camelToSnakeCase": camelToSnakeCase}).Parse(codeTemplate))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	// 将生成的代码输出为字符串
	return buf.String(), nil
}

// 工具函数：首字母小写
func lowerFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func lower(s string) string {
	return strings.ToLower(s)
}

// camelToSnakeCase 将驼峰命名转换为下划线命名
func camelToSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

func camelToSplitCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '/')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}
