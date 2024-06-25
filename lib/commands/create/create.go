package create

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lspaccatrosi16/go-cli-tools/input"
	"github.com/lspaccatrosi16/go-libs/structures/set"
)

func Run() error {
	name, err := inputS("Module name:")
	if err != nil {
		return err
	}
	oPath, err := inputS("Output path:")
	if err != nil {
		return err
	}

	selMods := set.NewSet[string]()
	opts := []input.SelectOption{
		{
			Name:  "build",
			Value: "build",
		},
		{
			Name:  "output",
			Value: "output",
		},
		{
			Name:  "static",
			Value: "static",
		},
		{
			Name:  "version",
			Value: "version",
		},
		{
			Name:  "done",
			Value: "done",
		},
	}

selMod:
	for {
		sel, err := input.GetSelection("Select modules to add", opts)
		if err != nil {
			return err
		}

		switch sel {
		case "done":
			break selMod
		default:
			selMods.Add(sel)
		}
	}
	sel := []string{}

	next := selMods.GetIterator()
	for v, ok := next(); ok; v, ok = next() {
		sel = append(sel, v)
	}

	yB := bytes.NewBuffer(nil)
	fmt.Fprintf(yB, "name: %s\n", name)

	if len(sel) > 0 {
		fmt.Fprintf(yB, "modules:\n")
		for _, s := range sel {
			mB := bytes.NewBuffer(nil)
			fmt.Fprintf(mB, "- name: %s\n", s)
			fmt.Fprintf(mB, "  config:\n")
			fmt.Fprintln(mB, addIndent(templates[s], 4))
			fmt.Fprintln(yB, addIndent(mB.String(), 2))
		}
	}

	f, err := os.Create(oPath)
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(f, yB)

	return nil
}

func addIndent(l string, n int) string {
	lines := strings.Split(l, "\n")
	for i, line := range lines {
		lines[i] = fmt.Sprintf("%s%s", strings.Repeat(" ", n), line)
	}
	return strings.Join(lines, "\n")
}

func inputS(msg string) (string, error) {
	var s string
	fmt.Print(msg)
	_, err := fmt.Scanln(&s)
	if err != nil {
		return "", err
	}

	return s, nil
}
