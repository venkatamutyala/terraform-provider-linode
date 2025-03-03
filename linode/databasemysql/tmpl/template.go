package tmpl

import (
	"github.com/linode/terraform-provider-linode/linode/acceptance"

	"testing"
)

type TemplateData struct {
	Engine          string
	Label           string
	AllowedIP       string
	ReplicationType string
	ClusterSize     int
	Encrypted       bool
	SSLConnection   bool
}

func Basic(t *testing.T, label, engine string) string {
	return acceptance.ExecuteTemplate(t,
		"database_mysql_basic", TemplateData{
			Engine: engine,
			Label:  label,
		})
}

func Complex(t *testing.T, data TemplateData) string {
	return acceptance.ExecuteTemplate(t,
		"database_mysql_complex", data)
}

func DataBasic(t *testing.T, data TemplateData) string {
	return acceptance.ExecuteTemplate(t,
		"database_mysql_data_basic", data)
}
