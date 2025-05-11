//go:build mage

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var supportedLangs = []string{"go"}

func Gen_proto() error {
	if exec.Command("proman").Run() != nil {
		if err := exec.Command("go", "install", "github.com/ondbyte/proman@latest").Run(); err != nil {
			return err
		}
	}

	for _, v := range supportedLangs {
		opDir := filepath.Join("internal", "protos", v)
		if _, err := os.Stat(opDir); err != nil {
			if err := os.MkdirAll(opDir, os.ModePerm); err != nil {
				return err
			}
		}
		cmd := exec.Command(
			"proman",
			"gen",
			"--lang",
			v,
			"--in",
			filepath.Join("internal", "protos"),
			"--out",
			opDir,
			"--grpc",
		)
		stdErr := bytes.NewBuffer(nil)
		cmd.Stderr = stdErr
		cmd.Stdout = stdErr
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("err while proman gen: %v\n%v", err, string(stdErr.Bytes()))
		}
	}
	fmt.Println("generated protos for langs: ", strings.Join(supportedLangs, ", "))
	return nil
}
