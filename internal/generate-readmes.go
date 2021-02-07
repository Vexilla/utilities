package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/vexilla/utilities/internal/color"

	"github.com/BurntSushi/toml"
)

type ReadmeContext struct {
	LanguageName string
	LanguageDisplayName string
	InstallInstructions string
	CustomInstanceHash string
	FetchFlags string
	SetupSnippet string
	Should string
	UsageSnippet string
	Example string
}

func GenerateReadmes() {

	// TODO: use os.WalkDir when migrating to Go 1.16
	files, err := ioutil.ReadDir("..")
	if err != nil {
		fmt.Println("Error reading parent")
	}
	fmt.Printf("found %v 'files' \n", len(files))

	for _, file := range files {

		// Log the IsDir()
		lineColor := color.Red
		if file.IsDir() {
			lineColor = color.Green
		}
		// fmt.Println(lineColor, file.IsDir(),color.Reset)

		if file.IsDir() &&
		strings.HasPrefix(file.Name(), "client-") {

			fmt.Println(file.Name())
			fmt.Println(lineColor, file.IsDir(), color.Reset)

			readmeFile, readmeError := os.OpenFile( filepath.Join("../", file.Name(), "/README.md"), os.O_WRONLY|os.O_CREATE, 0755)

			readmeTomlBytes, readmeTomlError := ioutil.ReadFile( filepath.Join("../", file.Name(), "/README.toml"))

			if readmeError != nil || readmeTomlError != nil {
				fmt.Println(color.Red, "Could not create or open", file.Name(), color.Reset)
			} else {
				defer readmeFile.Close()

				var readmeData ReadmeContext
				// json5Error := json5.Unmarshal(readmeJson5Bytes, &readmeData)
				_, tomlError := toml.Decode(string(readmeTomlBytes), &readmeData)

				if tomlError != nil {
					fmt.Println(color.Red, "Error Marshalling TOML", file.Name(), color.Reset, tomlError)
				}

				readmeTemplate, err := template.ParseFiles("./templates/README.tmpl")

				if err != nil {
					fmt.Println(color.Red, "Error rendering template", file.Name(), color.Reset)
				} else {
					fmt.Println(color.Green, readmeTemplate, file.Name(), color.Reset)
					err := readmeTemplate.Execute(readmeFile, readmeData)

					if err != nil {
						fmt.Println(color.Red, "Error writing file", file.Name(), color.Reset, err)
					}
				}

			}



		}

	}
}
