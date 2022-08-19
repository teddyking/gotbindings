package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const ServiceBindingRootEnvVar = "SERVICE_BINDING_ROOT"

type Status struct {
	GotBindings bool      `json:"gotBindings"`
	Bindings    []Binding `json:"bindings"`
}

type Binding struct {
	Name string `json:"name"`
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

		bindings := getBindings(bindingsRootDir)
		if len(bindings) > 0 {
			status.GotBindings = true
		}
		status.Bindings = bindings

		c.JSON(http.StatusOK, status)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func getBindings(bindingsRootDir string) []Binding {
	bindingDirs, err := os.ReadDir(bindingsRootDir)
	if err != nil {
		return []Binding{}
	}

	if len(bindingDirs) < 1 {
		return []Binding{}
	}

	bindings := []Binding{}

	for _, bindingDir := range bindingDirs {
		bindings = append(bindings, Binding{Name: bindingDir.Name()})
	}

	return bindings
}
