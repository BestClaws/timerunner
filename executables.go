package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"timerunner/pkg/jsondisk"
)

func AddExecutable(execReader io.Reader) (*RunInfo, error) {

	b, err := io.ReadAll(execReader)
	if err != nil {
		return nil, err
	}
	// execReader has bytes of zip file.
	// unzip in lib. iterate and save files.
	zipReader, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return nil, fmt.Errorf("unable to create zip reader: ", err)
	}
	if len(zipReader.File) < 2 {
		return nil, fmt.Errorf("package should contain atleast a description.json and an executable file")
	}

	// extract info from describe.json
	var runInfo RunInfo
	for _, eachFile := range zipReader.File {
		if eachFile.Name != "describe.json" {
			continue
		}

		file, err := eachFile.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		b, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read all bytes from description file, err: %s", err)
		}
		err = json.Unmarshal(b, &runInfo)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal, err: %s", err)
		}
		break
	}

	if runInfo.GetPeriod() == 0 || runInfo.GetOffset() == 0 || runInfo.ExecPath == "" || runInfo.PackageName == "" {
		return nil, fmt.Errorf("not enough info")
	}

	log.Printf("execpath: %s, start: %s, frequency: %v", runInfo.ExecPath, runInfo.GetOffset(), runInfo.GetPeriod())

	execDir := path.Join(executablesPath, runInfo.PackageName)
	err = os.Mkdir(execDir, 0744)
	if err != nil {
		if strings.HasSuffix(err.Error(), "file exists") {
			// folder already exists delete and create a fresh one to replace.
			log.Println("directory already exists. replacing, ", execDir)
			os.RemoveAll(execDir)
			os.Mkdir(execDir, 0744)
		}
	}

	for _, eachFile := range zipReader.File {
		file, err := eachFile.Open()
		if err != nil {
			return nil, fmt.Errorf("error trying to open file: %s", eachFile.Name)
		}

		data, err := io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read file into bytes. err: %s", err)
		}

		err = os.WriteFile(path.Join(executablesPath, runInfo.PackageName, eachFile.Name), data, 0744)
		if err != nil {
			return nil, fmt.Errorf("failed to persist files from zip package to disk. err: %s", err)
		}

	}

	// update descriptions.json
	descriptions, err := jsondisk.Load[[]RunInfo]("data/descriptions.json")
	if err != nil {
		return nil, err
	}
	// remove already present entry by swapping last element with the element to remove, and then removing last element
	for i, description := range *descriptions {
		if description.PackageName == runInfo.PackageName {
			(*descriptions)[i] = (*descriptions)[len(*descriptions)-1]
			*descriptions = (*descriptions)[:len(*descriptions)-1]
			break
		}
	}
	// add the new entry.
	*descriptions = append(*descriptions, runInfo)

	err = jsondisk.Save(*descriptions, "data/descriptions.json")
	if err != nil {
		return nil, err
	}

	return &runInfo, nil
}
