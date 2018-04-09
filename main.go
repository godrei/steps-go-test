package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
)

// Config ...
type Config struct {
	Packages string `env:"packages,required"`
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

	packages := strings.Split(cfg.Packages, ",")

	log.Infof("\nRunning go test...")

	for _, p := range packages {
		cmd := command.NewWithStandardOuts("go", "test", "-v", p)

		log.Printf("$ %s", cmd.PrintableCommandArgs())

		if err := cmd.Run(); err != nil {
			failf("go test failed: %s", err)
		}
	}
}
