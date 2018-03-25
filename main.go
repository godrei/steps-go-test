package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
	glob "github.com/ryanuber/go-glob"
)

type config struct {
	Include string `env:"include"`
	Exclude string `env:"exclude"`
}

func listFiles(include, exclude string) ([]string, error) {
	cmd := command.New("go", "list", "./...")
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return nil, err
	}

	split := strings.Split(out, "\n")
	var packages []string

	for _, p := range split {
		p = strings.TrimSpace(p)

		if include != "" && !glob.Glob(include, p) {
			return nil, nil
		}

		if exclude != "" && glob.Glob(exclude, p) {
			return nil, nil
		}

		packages = append(packages, p)
	}

	return packages, nil
}

func main() {
	var cfg config
	if err := stepconf.Parse(&cfg); err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}
	stepconf.Print(cfg)

	files, err := listFiles(cfg.Include, cfg.Exclude)
	if err != nil {
		log.Errorf("Failed to list files: %s", err)
		os.Exit(1)
	}

	cmd := command.NewWithStandardOuts("go", "test", strings.Join(files, "\n"))

	fmt.Println()
	log.Infof("$ %s", cmd.PrintableCommandArgs())

	if err := cmd.Run(); err != nil {
		log.Errorf("go test failed: %s", err)
		os.Exit(1)
	}
}
