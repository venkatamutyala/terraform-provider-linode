package databaseaccesscontrols_test

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/linode/linodego"
	"github.com/linode/terraform-provider-linode/linode/acceptance"
	"github.com/linode/terraform-provider-linode/linode/databaseaccesscontrols/tmpl"
	"github.com/linode/terraform-provider-linode/linode/helper"
)

var engineVersion string

func init() {
	client, err := acceptance.GetClientForSweepers()
	if err != nil {
		log.Fatalf("failed to get client: %s", err)
	}

	v, err := helper.ResolveValidDBEngine(context.Background(), *client, "mysql")
	if err != nil {
		log.Fatalf("failde to get db engine version: %s", err)
	}

	engineVersion = v.ID
}

func TestAccResourceDatabaseMySQLAccessControls_basic(t *testing.T) {
	t.Parallel()

	resName := "linode_database_access_controls.foobar"
	dbName := acctest.RandomWithPrefix("tf_test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: tmpl.Basic(t, dbName, engineVersion, "0.0.0.0/0"),
				Check: resource.ComposeTestCheckFunc(
					checkMySQLDatabaseExists,
					resource.TestCheckResourceAttr(resName, "allow_list.#", "1"),
					resource.TestCheckResourceAttr(resName, "allow_list.0", "0.0.0.0/0"),
				),
			},
			{
				Config: tmpl.Basic(t, dbName, engineVersion, "192.168.0.25/32"),
				Check: resource.ComposeTestCheckFunc(
					checkMySQLDatabaseExists,
					resource.TestCheckResourceAttr(resName, "allow_list.#", "1"),
					resource.TestCheckResourceAttr(resName, "allow_list.0", "192.168.0.25/32"),
				),
			},
			{
				ResourceName:      resName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func checkMySQLDatabaseExists(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_database_mysql" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}

		_, err = client.GetMySQLDatabase(context.Background(), id)
		if err != nil {
			return fmt.Errorf("error retrieving state of mysql database %s: %s", rs.Primary.Attributes["label"], err)
		}
	}

	return nil
}

func checkDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*helper.ProviderMeta).Client
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "linode_database_mysql" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing %v to int", rs.Primary.ID)
		}
		if id == 0 {
			return fmt.Errorf("Would have considered %v as %d", rs.Primary.ID, id)

		}

		_, err = client.GetMySQLDatabase(context.Background(), id)

		if err == nil {
			return fmt.Errorf("mysql database with id %d still exists", id)
		}

		if apiErr, ok := err.(*linodego.Error); ok && apiErr.Code != 404 {
			return fmt.Errorf("error requesting mysql database with id %d", id)
		}
	}

	return nil
}
