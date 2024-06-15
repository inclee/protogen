package gen_go

import (
	"bytes"
	"text/template"

	"github.com/inclee/gokit/container/maps"
	"github.com/inclee/protogen/internal/entity"
	"github.com/inclee/protogen/internal/errors"
	"github.com/inclee/protogen/internal/parse"
	"github.com/inclee/protogen/internal/utool"
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
	"{{ .Module }}/pkg/contexthelper"
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
	ctx := contexthelper.FromGin(c)
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

func genHandler(workdir, filepath string) (string, error) {
	compnent, err := parse.Parsefile(workdir, filepath)
	if err != nil {
		return "", err
	}
	if len(compnent.Services) == 0 {
		return "", errors.ErrNotNeedGen
	}
	//
	// 准备模板数据
	data := struct {
		Package  string
		Module   string
		Messages []*entity.Message
		Services []*entity.Service
	}{
		Package:  "",
		Module:   compnent.Module,
		Messages: maps.ValuesWith(compnent.Messages, func(msg *entity.Message) string { return msg.Name }, true),
		Services: maps.ValuesWith(compnent.Services, func(msg *entity.Service) string { return msg.Name }, true),
	}

	// 创建模板并解析
	tmpl := template.Must(template.New("code").Funcs(template.FuncMap{"toLower": utool.Lower, "lowerFirst": utool.LowerFirst, "camelToSnakeCase": utool.CamelToSnakeCase, "camelToSplitCase": utool.CamelToSplitCase}).Parse(handlerTemplate))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	// 将生成的代码输出为字符串
	return buf.String(), nil
}
