package main

import (
	"fmt"
	"go/build"
	"io"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/dave/jennifer/jen"
)

func hints(w io.Writer) error {

	// notest

	file := NewFile("jen")

	file.Comment("\tThis file is generated by genjen - do not edit!\n")
	file.Line()

	packages, err := getStandardLibraryPackages()
	if err != nil {
		return err
	}
	/*
		// Hints is a map containing hints for the names of all standard library packages
		var Hints = map[string]string{
			...
		}
	*/
	file.Comment("Hints is a map containing hints for the names of all standard library packages")
	file.Var().Id("Hints").Op("=").Map(String()).String().Values(DictFunc(func(d Dict) {
		for path, name := range packages {
			if name == "main" {
				continue
			}
			d[Lit(path)] = Lit(name)
		}
	}))

	return file.Render(w)
}

func getStandardLibraryPackages() (map[string]string, error) {

	// notest

	cmd := exec.Command("go", "list", "-f", "{{ .ImportPath }} {{ .Name }}", "./...")
	cmd.Env = []string{
		fmt.Sprintf("GOPATH=%s", build.Default.GOPATH),
		fmt.Sprintf("GOROOT=%s", build.Default.GOROOT),
	}
	cmd.Dir = filepath.Join(build.Default.GOROOT, "src")
	b, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	all := strings.Split(strings.TrimSpace(string(b)), "\n")

	packages := map[string]string{}
	for _, j := range all {
		parts := strings.Split(j, " ")
		path := parts[0]
		name := parts[1]
		packages[path] = name
	}
	return packages, nil
}