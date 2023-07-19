package template

import (
	"bytes"
	template2 "html/template"
	"prometheus-webhook-wechat/models"
)

func TransferContent(srcContent models.Data, templateFile *string) (alertContent string, err error) {
	buf := new(bytes.Buffer)
	tpl, err := template2.New(*templateFile).Funcs(DefaultFuncs).ParseFiles(*templateFile)
	if err != nil {
		return
	}

	err = tpl.Execute(buf, srcContent)
	if err != nil {
		return
	}
	return buf.String(), nil
}
