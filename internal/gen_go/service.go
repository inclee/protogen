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

func genService(workdir, filepath string) (string, error) {
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
	tmpl := template.Must(template.New("code").Funcs(template.FuncMap{"lower": utool.Lower, "lowerFirst": utool.LowerFirst, "camelToSnakeCase": utool.CamelToSnakeCase}).Parse(codeTemplate))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	// 将生成的代码输出为字符串
	return buf.String(), nil
}
