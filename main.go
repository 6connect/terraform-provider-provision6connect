package main

import (
	"context"
	"terraform-provider-provision6connect/provision6connect"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name provision6connect
func main() {
	providerserver.Serve(context.Background(), provision6connect.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/6connect/provision6connect",
	})
}
