package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	generator "github.com/SAP/jenkins-library/pkg/documentation/generator"
	"github.com/SAP/jenkins-library/pkg/generator/helper"
)

type sliceFlags struct {
	list []string
}

func (f *sliceFlags) String() string {
	return ""
}

func (f *sliceFlags) Set(value string) error {
	f.list = append(f.list, value)
	return nil
}

func main() {
	var metadataPath string
	var docTemplatePath string
	var customDefaultFiles sliceFlags

	flag.StringVar(&metadataPath, "metadataDir", "./resources/metadata", "The directory containing the step metadata. Default points to \\'resources/metadata\\'.")
	flag.StringVar(&docTemplatePath, "docuDir", "./documentation/docs/steps/", "The directory containing the docu stubs. Default points to \\'documentation/docs/steps/\\'.")
	// flag.StringVar(&customDefaultPaths, "customDefaultPath", "./resources/default_pipeline_environment.yml", "Comma-separated paths to custom default configuration files.")
	flag.Var(&customDefaultFiles, "customDefaultFile", "Path to a custom default configuration file.")
	flag.Parse()

	fmt.Println("using Metadata Directory:", metadataPath)
	fmt.Println("using Documentation Directory:", docTemplatePath)
	fmt.Println("using Custom Default Files:", strings.Join(customDefaultFiles.list, ", "))

	metadataFiles, err := helper.MetadataFiles(metadataPath)
	checkError(err)
	err = generator.GenerateStepDocumentation(metadataFiles, customDefaultFiles.list, generator.DocuHelperData{
		DocTemplatePath:     docTemplatePath,
		OpenDocTemplateFile: openDocTemplateFile,
		DocFileWriter:       writeFile,
		OpenFile:            openFile,
	})
	checkError(err)
}

func openDocTemplateFile(docTemplateFilePath string) (io.ReadCloser, error) {
	//check if template exists otherwise print No Template found
	if _, err := os.Stat(docTemplateFilePath); os.IsNotExist(err) {
		err := fmt.Errorf("no template found: %v", docTemplateFilePath)
		return nil, err
	}

	return os.Open(docTemplateFilePath)
}

func writeFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}

func openFile(name string) (io.ReadCloser, error) {
	return os.Open(name)
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("Error occured: %v\n", err)
		os.Exit(1)
	}
}