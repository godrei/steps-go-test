package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
	"github.com/godrei/steps-golint/gotool"
)

// Config ...
type Config struct {
	Exclude string `env:"exclude"`
}

func installedInPath(name string) bool {
	cmd := exec.Command("which", name)
	outBytes, err := cmd.Output()
	return err == nil && strings.TrimSpace(string(outBytes)) != ""
}

func failf(format string, args ...interface{}) {
	log.Errorf(format, args...)
	os.Exit(1)
}

func main() {
	var cfg Config
	if err := stepconf.Parse(&cfg); err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}
	stepconf.Print(cfg)

	dir, err := os.Getwd()
	if err != nil {
		failf("Failed to get working directory: %s", err)
	}

	excludes := strings.Split(cfg.Exclude, ",")

	packages, err := gotool.ListPackages(dir, excludes...)
	if err != nil {
		failf("Failed to list packages: %s", err)
	}

	log.Infof("\nRunning go test...")

	for _, p := range packages {
		cmd := command.NewWithStandardOuts("go", "test", p)

		log.Printf("$ %s", cmd.PrintableCommandArgs())

		if err := cmd.Run(); err != nil {
			failf("go test failed: %s", err)
		}
	}
}
