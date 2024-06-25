package create

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/lspaccatrosi16/go-cli-tools/input"
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

	opts := []input.SelectOption{
		{
			Name:  "Build",
			Value: "build",
		},
		{
			Name:  "Output",
			Value: "output",
		},
		{
			Name:  "Static",
			Value: "static",
		},
		{
			Name:  "Version",
			Value: "version",
		},
		{
			Name:  "Done",
			Value: "done",
		},
	}

	selPackages := []string{}

selMod:
	for {
		sel, idx, err := input.GetSelectionIdx("Select modules to add", opts)
		if err != nil {
			return err
		}

		switch sel {
		case "done":
			break selMod
		default:
			selPackages = append(selPackages, sel)
			fmt.Printf("%s added\n", opts[idx].Name)
			opts = append(opts[:idx], opts[idx+1:]...)
		}
	}
	yB := bytes.NewBuffer(nil)
	fmt.Fprintf(yB, "name: %s\n", name)

	if len(selPackages) > 0 {
		fmt.Fprintf(yB, "modules:\n")
		for _, s := range selPackages {
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
