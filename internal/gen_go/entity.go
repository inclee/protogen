package gen_go

import (
	"bytes"
	"text/template"

	"github.com/inclee/gokit/container/maps"
	"github.com/inclee/protogen/internal/entity"
	"github.com/inclee/protogen/internal/parse"
	"github.com/inclee/protogen/internal/utool"
)

// 生成代码的模板
const entiryTemplate = `
package entity

/*GENERATE BY protogen(https://github.com/inclee/protogen); PLEASE DON'T EDIT IT */

{{range .Messages}}
// {{.Name}} 请求结构体
{{- $Method := .Method }}
{{- $Parent := .Parent }}
type {{.Name}} struct {
	{{- if $Parent }}
		{{ $Parent.Name }} 
	{{- end }}
	{{- range .Fields }} 
		{{- $Validate := "" }}
		{{- if .Validate }}
			{{- $Validate = printf "validate:\"%s\"" .Validate }}
		{{- end }}
		{{- if .Tag }}
			{{ .Name }} {{ .Type }} ` + "`{{.Tag}}:\"{{ .Name | camelToSnakeCase }}\" {{ $Validate }} `" + `
		{{- else }}
			{{ .Name }} {{ .Type }} 
		{{- end }}
	{{- end }}
}
{{end}}

`

func genEntity(workdir, filepath string) (string, error) {
	compnent, err := parse.Parsefile(workdir, filepath)
	if err != nil {
		return "", err
	}
	//
	// 准备模板数据
	data := struct {
		Package  string
		Module   string
		Messages []*entity.Message
	}{
		Package:  "",
		Module:   compnent.Module,
		Messages: maps.ValuesWith(compnent.Messages, func(msg *entity.Message) string { return msg.Name }, true),
	}

	// 创建模板并解析
	tmpl := template.Must(template.New("code").Funcs(template.FuncMap{"lowerFirst": utool.LowerFirst, "camelToSnakeCase": utool.CamelToSnakeCase}).Parse(entiryTemplate))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	// 将生成的代码输出为字符串
	return buf.String(), nil
}
