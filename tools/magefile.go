//+build mage

package main

import (
	"context"
	"log"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"

	tools "github.com/mrcrilly/terraform-provider-awx/tools"

	// mage:import
	"github.com/nolte/plumbing/cmd/kind"
	_ "github.com/nolte/plumbing/cmd/kind"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Delete the cluster.
func InstallAWX() error {
	log.Printf("InstallAWX to Kind Cluster")

	return sh.Run(
		"sh", "./installAwx.sh")
}
func GenDocumentuation() error {
	return tools.GenerateProviderCoumentation()
}

// ReCreate a kind Cluster with Awx support.
func ReCreate(ctx context.Context) {
	log.Printf("Create Kind Cluster with AWX")
	mg.CtxDeps(ctx, kind.Kind.Recreate)
	mg.CtxDeps(ctx, InstallAWX)
}
