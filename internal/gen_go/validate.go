package gen_go

import "path/filepath"

var validatorTemplate = `
package handler
/*GENERATE BY protogen(https://github.com/inclee/protogen); PLEASE DON'T EDIT IT */
import(
	"github.com/go-playground/validator/v10"
)
var validate *validator.Validate = validator.New()
`

func writeValide(codeDir string) error {
	handerfile := filepath.Join(codeDir, "handler/private_hander.gen.go")
	return writeFile(handerfile, validatorTemplate)
}
