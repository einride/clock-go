package main

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/magefile/mage/mg"
	"go.einride.tech/mage-tools/mglogr"
	"go.einride.tech/mage-tools/mgmake"
	"go.einride.tech/mage-tools/mgpath"
	"go.einride.tech/mage-tools/mgtool"
	"go.einride.tech/mage-tools/tools/mgconvco"
	"go.einride.tech/mage-tools/tools/mggit"
	"go.einride.tech/mage-tools/tools/mggo"
	"go.einride.tech/mage-tools/tools/mggolangcilint"
	"go.einride.tech/mage-tools/tools/mggoreview"
	"go.einride.tech/mage-tools/tools/mgmarkdownfmt"
	"go.einride.tech/mage-tools/tools/mgyamlfmt"
)

func init() {
	mgmake.GenerateMakefiles(
		mgmake.Makefile{
			Path:          mgpath.FromGitRoot("Makefile"),
			DefaultTarget: All,
		},
	)
}

func All() {
	mg.Deps(
		ConvcoCheck,
		GolangciLint,
		Goreview,
		GoTest,
		FormatMarkdown,
		FormatYaml,
	)
	mg.SerialDeps(
		GoModTidy,
		GitVerifyNoDiff,
	)
}

func GitVerifyNoDiff(ctx context.Context) error {
	ctx = logr.NewContext(ctx, mglogr.New("git-verify-no-diff"))
	logr.FromContextOrDiscard(ctx).Info("verifying that git has no diff...")
	return mggit.VerifyNoDiff(ctx)
}

func GoModTidy(ctx context.Context) error {
	ctx = logr.NewContext(ctx, mglogr.New("go-mod-tidy"))
	logr.FromContextOrDiscard(ctx).Info("tidying Go module files...")
	return mgtool.Command(ctx, "go", "mod", "tidy", "-v").Run()
}

func Goreview(ctx context.Context) error {
	ctx = logr.NewContext(ctx, mglogr.New("goreview"))
	logr.FromContextOrDiscard(ctx).Info("running...")
	return mggoreview.Command(ctx, "-c", "1", "./...").Run()
}

func GoTest(ctx context.Context) error {
	ctx = logr.NewContext(ctx, mglogr.New("go-test"))
	logr.FromContextOrDiscard(ctx).Info("running Go tests...")
	return mggo.TestCommand(ctx).Run()
}

func GolangciLint(ctx context.Context) error {
	ctx = logr.NewContext(ctx, mglogr.New("golangci-lint"))
	logr.FromContextOrDiscard(ctx).Info("linting Go files...")
	return mggolangcilint.RunCommand(ctx).Run()
}

func FormatMarkdown(ctx context.Context) error {
	ctx = logr.NewContext(ctx, mglogr.New("format-markdown"))
	logr.FromContextOrDiscard(ctx).Info("formatting Markdown files...")
	return mgmarkdownfmt.Command(ctx, "-w", ".").Run()
}

func FormatYaml(ctx context.Context) error {
	ctx = logr.NewContext(ctx, mglogr.New("format-yaml"))
	logr.FromContextOrDiscard(ctx).Info("formatting YAML files...")
	return mgyamlfmt.FormatYAML(ctx)
}

func ConvcoCheck(ctx context.Context) error {
	ctx = logr.NewContext(ctx, mglogr.New("convco-check"))
	logr.FromContextOrDiscard(ctx).Info("checking...")
	return mgconvco.Command(ctx, "check", "origin/master..HEAD").Run()
}
