package tmpl

import (
	"testing"

	"github.com/linode/terraform-provider-linode/linode/acceptance"
	databasemysqltmpl "github.com/linode/terraform-provider-linode/linode/databasemysql/tmpl"
)

type TemplateData struct {
	DB     databasemysqltmpl.TemplateData
	Engine string
	Label  string
}

func ByLabel(t *testing.T, engineVersion, instLabel, dsLabel string) string {
	return acceptance.ExecuteTemplate(t,
		"databases_data_by_label", TemplateData{
			DB:    databasemysqltmpl.TemplateData{Engine: engineVersion, Label: instLabel},
			Label: dsLabel,
		})
}

func ByEngine(t *testing.T, engineVersion, label, engine string) string {
	return acceptance.ExecuteTemplate(t,
		"databases_data_by_engine", TemplateData{
			DB:     databasemysqltmpl.TemplateData{Engine: engineVersion, Label: label},
			Engine: engine,
		})
}
