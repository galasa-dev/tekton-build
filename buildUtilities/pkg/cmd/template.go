//
// Licensed Materials - Property of IBM
//
// (c) Copyright IBM Corp. 2021.
//

package cmd

import (
	"fmt"
	"io/ioutil"
	"bytes"
	"text/template"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"galasa.dev/buildUtilities/pkg/galasayaml"
)

var (
	templateFile       string
	releaseMetadata    string
	outputFile         string
	requireObr         bool
	requireBom         bool
	requireMvp         bool
	requireIsolated    bool
	requireJavadoc     bool
	requireManagerdoc  bool

	templateCmd = &cobra.Command{
		Use:   "template",
		Short: "generates files from a template",
		Long:  "Generates files from a template using the Galasa release metadata file",
		Run:   execute,
	}

	release galasayaml.Release
)

type templateData struct {
	Release string
	Artifacts []artifact
}
type artifact struct {
	GroupId       string
	ArtifactId    string
	Version       string
	Type          string
}


func init() {
	templateCmd.PersistentFlags().StringVarP(&templateFile, "template", "t", "", "template file")
	templateCmd.PersistentFlags().StringVarP(&releaseMetadata, "releaseMetadata", "r", "", "release metadata file")
	templateCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "output file")

	templateCmd.PersistentFlags().BoolVarP(&requireObr, "obr", "", false, "require maven artifacts for OBR")
	templateCmd.PersistentFlags().BoolVarP(&requireBom, "bom", "", false, "require maven artifacts for BOM")
	templateCmd.PersistentFlags().BoolVarP(&requireMvp, "mvp", "", false, "require maven artifacts for mvp zip")
	templateCmd.PersistentFlags().BoolVarP(&requireIsolated, "isolated", "", false, "require maven artifacts for isolated zip")
	templateCmd.PersistentFlags().BoolVarP(&requireJavadoc, "javadoc", "", false, "require maven artifacts for javadoc")
	templateCmd.PersistentFlags().BoolVarP(&requireManagerdoc, "managerdoc", "", false, "require maven artifacts for manager docs")

	rootCmd.AddCommand(templateCmd)
}

func execute(cmd *cobra.Command, args []string) {
	fmt.Println("Galasa Build - Template")

	if releaseMetadata == "" {
		panic("Release metadata file has not been provided")
	}
	
	if templateFile == "" {
		panic("Template file has not been provided")
	}
	
	if outputFile == "" {
		fmt.Println("Output file has not been provided")
	}
	
	b, err := ioutil.ReadFile(releaseMetadata)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(b, &release)
	if err != nil {
		panic(err)
	}

	var requires = 0;
	if requireObr {
		requires++
		fmt.Println("OBR artifact type requested")
	}
	if requireBom {
		requires++
		fmt.Println("BOM artifact type requested")
	}
	if requireMvp {
		requires++
		fmt.Println("MVP artifact type requested")
	}
	if requireIsolated {
		requires++
		fmt.Println("Isolated artifact type requested")
	}
	if requireJavadoc {
		requires++
		fmt.Println("Javadoc artifact type requested")
	}
	if requireManagerdoc {
		requires++
		fmt.Println("Manager Docs artifact type requested")
	}

	if requires == 0 {
		panic("Artifact type has not been provided")
	}

	if requires > 1 {
		panic("Too many artifact types have been requested")
	}

	t := templateData {}

	t.Release = release.Release.Version

	for _, bundle := range release.Framework.Bundles {
		if bundle.Group == "" {
			bundle.Group = "dev.galasa"
		}

		selected := false
		
		if requireObr {
			selected = true
		} else if requireBom {
			selected = bundle.Bom
		} else if requireMvp {
			selected = bundle.Mvp
		} else if requireIsolated {
			selected = true
		} else if requireJavadoc {
			selected = bundle.Javadoc
		} else if requireManagerdoc {
			selected = bundle.Managerdoc
		}
		
		if selected {
			artifact := artifact {
				GroupId: bundle.Group,
				ArtifactId: bundle.Artifact,
				Version: bundle.Version,
				Type: bundle.Type,
			}

			t.Artifacts = append(t.Artifacts, artifact)

			fmt.Printf("    Added framework artifact %v:%v:%v\n", artifact.GroupId, artifact.ArtifactId, artifact.Version)
		}	
	}
	
	for _, bundle := range release.Api.Bundles {
		if bundle.Group == "" {
			bundle.Group = "dev.galasa"
		}

		selected := false
		
		if requireObr {
			selected = true
		} else if requireBom {
			selected = bundle.Bom
		} else if requireMvp {
			selected = bundle.Mvp 
		} else if requireIsolated {
			selected = true
		} else if requireJavadoc {
			selected = bundle.Javadoc
		} else if requireManagerdoc {
			selected = bundle.Managerdoc
		}
		
		if selected {
			artifact := artifact {
				GroupId: bundle.Group,
				ArtifactId: bundle.Artifact,
				Version: bundle.Version,
				Type: bundle.Type,
			}

			t.Artifacts = append(t.Artifacts, artifact)

			fmt.Printf("    Added framework artifact %v:%v:%v\n", artifact.GroupId, artifact.ArtifactId, artifact.Version)
		}	
	}
	
	for _, bundle := range release.Managers.Bundles {
		if bundle.Group == "" {
			bundle.Group = "dev.galasa"
		}

		selected := false
		
		if requireObr {
			selected = true
		} else if requireBom {
			selected = true
		} else if requireMvp {
			selected = bundle.Mvp 
		} else if requireIsolated {
			selected = true
		} else if requireJavadoc {
			selected = true
		} else if requireManagerdoc {
			selected = true
		}
		
		if selected {
			artifact := artifact {
				GroupId: bundle.Group,
				ArtifactId: bundle.Artifact,
				Version: bundle.Version,
				Type: bundle.Type,
			}

			t.Artifacts = append(t.Artifacts, artifact)

			fmt.Printf("    Added framework artifact %v:%v:%v\n", artifact.GroupId, artifact.ArtifactId, artifact.Version)
		}	
	}
	
	for _, bundle := range release.External.Bundles {
		if bundle.Group == "" {
			bundle.Group = "dev.galasa"
		}

		selected := false
		
		if requireObr {
			selected = bundle.Obr
		} else if requireBom {
			selected = bundle.Bom
		} else if requireMvp {
			selected = bundle.Mvp 
		} else if requireIsolated {
			selected = bundle.Isolated
		} else if requireJavadoc {
			selected = false
		} else if requireManagerdoc {
			selected = false
		}
		
		if selected {
			artifact := artifact {
				GroupId: bundle.Group,
				ArtifactId: bundle.Artifact,
				Version: bundle.Version,
				Type: bundle.Type,
			}

			t.Artifacts = append(t.Artifacts, artifact)

			fmt.Printf("    Added framework artifact %v:%v:%v\n", artifact.GroupId, artifact.ArtifactId, artifact.Version)
		}	
	}
	
	b, err = ioutil.ReadFile(templateFile)
	if err != nil {
		panic(err)
	}

	templString := string(b)

	tmpl, err := template.New("convert").Parse(templString)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, t)
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(outputFile, buf.Bytes(), 0644)

}