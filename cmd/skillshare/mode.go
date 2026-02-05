package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type runMode int

const (
	modeAuto runMode = iota
	modeProject
	modeGlobal
)

func parseModeArgs(args []string) (runMode, []string, error) {
	mode := modeAuto
	rest := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--project", "-p":
			if mode == modeGlobal {
				return modeAuto, nil, fmt.Errorf("--project and --global are mutually exclusive")
			}
			mode = modeProject
		case "--global", "-g":
			if mode == modeProject {
				return modeAuto, nil, fmt.Errorf("--project and --global are mutually exclusive")
			}
			mode = modeGlobal
		default:
			rest = append(rest, args[i])
		}
	}

	return mode, rest, nil
}

func projectConfigExists(root string) bool {
	_, err := os.Stat(filepath.Join(root, ".skillshare", "config.yaml"))
	return err == nil
}
