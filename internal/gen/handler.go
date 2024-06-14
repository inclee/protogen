package gen

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"text/template"
)

// 生成代码的模板
const handlerTemplate = `
package handler

/*GENERATE BY protogen(https://github.com/inclee/protogen); PLEASE DON'T EDIT IT */

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"{{ .Module }}/internal/service"
	"{{ .Module }}/internal/entity"
	"{{ .Module }}/pkg/context"
	"{{ .Module }}/pkg/log"
)
{{range .Services}}
// 实现 {{.Name}} 服务
type {{ .Name }}Handler struct{
	srv service.I{{ .Name }}Service
}

{{- $serviceName := .Name }}
func New{{ .Name }}Handler(srv service.I{{ .Name }}Service) *{{ .Name }}Handler{
	return &{{ .Name }}Handler{
		srv: srv,
	}
}

func (h *{{ .Name }}Handler) Register(group *gin.RouterGroup) {
	{{- range .Handlers }}
		group.{{ .Method }}("{{ .Name | camelToSplitCase }}", h.{{ .Name }})	
	{{- end }}		
}
{{- range .Handlers }}
func (h *{{ $serviceName }}Handler) {{ .Name }}(c *gin.Context) () {
	ctx := context.FromGin(c)
	log := log.New(ctx)
	req := &entity.{{ .Request }}{}
	{{- if eq (toLower .Method) "get" }}
		if err := c.ShouldBindQuery(req);err != nil{
	{{- else }}
		if err := c.Bind(req);err != nil{
	{{- end }}
		c.Status(http.StatusBadRequest)
		log.Errorf("%v.%v params error:%v", "{{ $serviceName }}Handler", "{{ .Name }}", err)
		return 	
	}
	if err := validate.Struct(req); err != nil {
		c.Status(http.StatusBadRequest)
		log.Errorf("%v.%v params error:%v", "{{ $serviceName }}Handler", "{{ .Name }}", err)
		return
	}
	rsp ,err := h.srv.{{ .Name }}(ctx,req)
	if err != nil{
		c.Status(http.StatusInternalServerError)
		log.Errorf("%v.%v process error:%v", "{{ $serviceName }}Handler", "{{ .Name }}", err)
		return
	}
	c.JSON(http.StatusOK,rsp)
}
{{- end }}
{{end}}
`

func parseHandler(filepath string) (string, error) {
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
	var currentModule string
	var currentPackage string

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
			method := parts[0]
			handlerName := parts[1]
			request := strings.Trim(parts[2], "()")
			response := strings.Trim(parts[3], "()")
			currentService.Handlers = append(currentService.Handlers, HandlerFunc{Method: method, Name: handlerName, Request: request, Response: response})
		}
	}

	// 检查是否有扫描错误
	if err := scanner.Err(); err != nil {
		return "", err
	}
	// 准备模板数据
	data := struct {
		Package  string
		Module   string
		Services []Service
	}{
		Package:  currentPackage,
		Module:   currentModule,
		Services: services,
	}

	// 创建模板并解析
	tmpl := template.Must(template.New("code").Funcs(template.FuncMap{"toLower": lower, "lowerFirst": lowerFirst, "camelToSnakeCase": camelToSnakeCase, "camelToSplitCase": camelToSplitCase}).Parse(handlerTemplate))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	// 将生成的代码输出为字符串
	return buf.String(), nil
}
