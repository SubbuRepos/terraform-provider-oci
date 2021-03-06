// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const (
	ExportSetRequiredOnlyResource = ExportSetResourceDependencies + `
resource "oci_file_storage_export_set" "test_export_set" {
	#Required
	mount_target_id = "${oci_file_storage_mount_target.test_mount_target.id}"
}
`

	ExportSetResourceConfig = ExportSetResourceDependencies + `
resource "oci_file_storage_export_set" "test_export_set" {
	#Required
	mount_target_id = "${oci_file_storage_mount_target.test_mount_target.id}"

	# Optional
	display_name = "${var.export_set_display_name}"
	max_fs_stat_bytes = "${var.max_bytes}"
	max_fs_stat_files = "${var.max_files}"
}
`
	ExportSetPropertyVariables = `
variable "export_set_availability_domain" { default = "kIdk:PHX-AD-1" }
variable "export_set_display_name" { default = "export set display name" }
variable "max_bytes" { default = 23843202333 }
variable "max_files" { default = 223442 }
variable "export_set_state" { default = "ACTIVE" }
`
	ExportSetResourceDependencies = MountTargetPropertyVariables + MountTargetResourceConfig
)

func TestFileStorageExportSetResource_basic(t *testing.T) {
	provider := testAccProvider
	config := testProviderConfig()

	compartmentId := getRequiredEnvSetting("compartment_id_for_create")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	resourceName := "oci_file_storage_export_set.test_export_set"
	datasourceName := "data.oci_file_storage_export_sets.test_export_sets"

	var resId, resId2 string

	resource.Test(t, resource.TestCase{
		Providers: map[string]terraform.ResourceProvider{
			"oci": provider,
		},
		Steps: []resource.TestStep{
			// verify create - note that we don't really create an export set, see provider for details.
			{
				ImportState:       true,
				ImportStateVerify: true,
				Config:            config + ExportSetPropertyVariables + compartmentIdVariableStr + ExportSetRequiredOnlyResource,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "availability_domain", "kIdk:PHX-AD-1"),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					//resource.TestCheckResourceAttrSet(resourceName, "display_name"),
					resource.TestCheckResourceAttrSet(resourceName, "max_fs_stat_bytes"),
					resource.TestCheckResourceAttrSet(resourceName, "max_fs_stat_files"),
					resource.TestCheckResourceAttrSet(resourceName, "mount_target_id"),
					resource.TestCheckResourceAttr(resourceName, "state", "ACTIVE"),
					resource.TestCheckResourceAttrSet(resourceName, "time_created"),

					func(s *terraform.State) (err error) {
						resId, err = fromInstanceState(s, resourceName, "id")
						return err
					},
				),
			},
			// This step serves the purpose of both "create with optionals" and "update non-forcenew fields". See provider for details.
			{
				Config: config + ExportSetPropertyVariables + compartmentIdVariableStr + ExportSetResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "availability_domain", "kIdk:PHX-AD-1"),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "display_name", "export set display name"),
					resource.TestCheckResourceAttr(resourceName, "max_fs_stat_bytes", "23843202333"),
					resource.TestCheckResourceAttr(resourceName, "max_fs_stat_files", "223442"),
					resource.TestCheckResourceAttrSet(resourceName, "mount_target_id"),
					resource.TestCheckResourceAttr(resourceName, "state", "ACTIVE"),
					resource.TestCheckResourceAttrSet(resourceName, "time_created"),

					func(s *terraform.State) (err error) {
						resId2, err = fromInstanceState(s, resourceName, "id")
						if resId != resId2 {
							return fmt.Errorf("Resource recreated when it was supposed to be updated.")
						}
						return err
					},
				),
			},
			// verify updates to Force New parameters.
			{
				Config: config + `
variable "export_set_availability_domain" { default = "kIdk:PHX-AD-1" }
variable "export_set_display_name" { default = "export set on mount target 2" }
variable "max_bytes" { default = 23843202333 }
variable "max_files" { default = 223442 }
variable "export_set_state" { default = "ACTIVE" }

# creating a second mount target
resource "oci_file_storage_mount_target" "test_mount_target_2" {
	#Required
	availability_domain = "${var.mount_target_availability_domain}"
	compartment_id = "${var.compartment_id}"
	subnet_id = "${oci_core_subnet.test_subnet.id}"
}

# Using the same test_export_set variable, but specifying different mount target 
resource "oci_file_storage_export_set" "test_export_set" {
	#Required
	mount_target_id = "${oci_file_storage_mount_target.test_mount_target_2.id}"

	# Optional
	display_name = "${var.export_set_display_name}"
	max_fs_stat_bytes = "${var.max_bytes}"
	max_fs_stat_files = "${var.max_files}"
}
                ` + compartmentIdVariableStr + MountTargetPropertyVariables + MountTargetResourceDependencies,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "availability_domain", "kIdk:PHX-AD-1"),
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "display_name", "export set on mount target 2"),
					resource.TestCheckResourceAttr(resourceName, "max_fs_stat_bytes", "23843202333"),
					resource.TestCheckResourceAttr(resourceName, "max_fs_stat_files", "223442"),
					resource.TestCheckResourceAttrSet(resourceName, "mount_target_id"),
					resource.TestCheckResourceAttr(resourceName, "state", "ACTIVE"),
					resource.TestCheckResourceAttrSet(resourceName, "time_created"),

					func(s *terraform.State) (err error) {
						resId2, err = fromInstanceState(s, resourceName, "id")
						if resId == resId2 {
							return fmt.Errorf("Resource was expected to be recreated but it wasn't.")
						}
						return err
					},
				),
			},
			// verify datasource
			{
				Config: config + `
variable "export_set_availability_domain" { default = "kIdk:PHX-AD-1" }
variable "export_set_display_name" { default = "export set display name" }
variable "max_bytes" { default = 23843202333 }
variable "max_files" { default = 223442 }
variable "export_set_state" { default = "ACTIVE" }

data "oci_file_storage_export_sets" "test_export_sets" {
	#Required
	availability_domain = "${var.export_set_availability_domain}"
	compartment_id = "${var.compartment_id}"

	#Optional
	display_name = "${var.export_set_display_name}"
	id = "${oci_file_storage_mount_target.test_mount_target.export_set_id}"
	state = "${var.export_set_state}"

    filter {
    	name = "id"
    	values = ["${oci_file_storage_mount_target.test_mount_target.export_set_id}"]
    }
}
                ` + compartmentIdVariableStr + ExportSetResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "availability_domain", "kIdk:PHX-AD-1"),
					resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),

					resource.TestCheckResourceAttr(datasourceName, "export_sets.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceName, "export_sets.0.compartment_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "export_sets.0.display_name"),
					resource.TestCheckResourceAttrSet(datasourceName, "export_sets.0.id"),
					resource.TestCheckResourceAttr(datasourceName, "export_sets.0.state", "ACTIVE"),
					resource.TestCheckResourceAttrSet(datasourceName, "export_sets.0.time_created"),
					// resource.TestCheckResourceAttrSet(datasourceName, "export_sets.0.vcn_id"),
				),
			},
		},
	})
}
