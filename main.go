package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/manifoldco/promptui"
	homedir "github.com/mitchellh/go-homedir"
	"gopkg.in/ini.v1"
)

func getProfiles() ([]string, error) {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	cfg, err := ini.Load(filepath.Join(home, ".aws", "credentials"))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	sections := cfg.SectionStrings()
	// Remove DEFAULT section added by ini parser
	fixedSections := append(sections[:0], sections[1:]...)
	return fixedSections, nil
}

func main() {
	sections, err := getProfiles()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	prompt := promptui.Select{
		Label: "AWS Profile",
		Items: sections,
		Size:  15,
		Searcher: func(input string, idx int) bool {
			if strings.Contains(sections[idx], input) {
				return true
			}
			return false
		},
	}
	_, result, err := prompt.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Copied to clipboard: export AWS_PROFILE=%s\n", result)
	export := fmt.Sprintf("export AWS_PROFILE=%s", result)
	clipboard.WriteAll(export)
}
