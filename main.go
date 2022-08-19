package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const ServiceBindingRootEnvVar = "SERVICE_BINDING_ROOT"

type Status struct {
	GotBindings bool      `json:"gotBindings"`
	Bindings    []Binding `json:"bindings"`
}

type Binding struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Entries []string `json:"entries"`
}

func main() {
	bindingsRootDir := os.Getenv(ServiceBindingRootEnvVar)
	if bindingsRootDir == "" {
		log.Fatalf("%s env var not set", ServiceBindingRootEnvVar)
	}

	_, err := os.Stat(bindingsRootDir)
	if os.IsNotExist(err) {
		log.Fatalf("%s dir does not exist", bindingsRootDir)
	}

	log.Printf("bindings root dir set to %s", bindingsRootDir)

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		var status Status

		bindings, err := getBindings(bindingsRootDir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		if len(bindings) > 0 {
			status.GotBindings = true
		}
		status.Bindings = bindings

		c.JSON(http.StatusOK, status)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func getBindings(bindingsRootDir string) ([]Binding, error) {
	bindingDirs, err := os.ReadDir(bindingsRootDir)
	if err != nil {
		return []Binding{}, err
	}

	if len(bindingDirs) < 1 {
		return []Binding{}, nil
	}

	bindings := []Binding{}
	for _, bindingDir := range bindingDirs {
		binding := Binding{
			Name: bindingDir.Name(),
		}

		bindingEntryFiles, err := os.ReadDir(filepath.Join(bindingsRootDir, bindingDir.Name()))
		if err != nil {
			return []Binding{}, err
		}

		entries := []string{}
		for _, bindingEntryFile := range bindingEntryFiles {
			if bindingEntryFile.Name() == "type" {
				bindingType, err := os.ReadFile(filepath.Join(bindingsRootDir, bindingDir.Name(), bindingEntryFile.Name()))
				if err != nil {
					return []Binding{}, err
				}

				binding.Type = strings.TrimSpace(string(bindingType))
			} else {
				entries = append(entries, bindingEntryFile.Name())
			}
		}

		binding.Entries = entries
		bindings = append(bindings, binding)
	}

	return bindings, nil
}
