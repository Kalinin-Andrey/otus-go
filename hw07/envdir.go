package main

import (
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"os/exec"
)

const maxNumberOfFiles = 100

func main() {
	args := os.Args
	path := args[1]
	prog := args[2]
	envs, err := getEnvsFromDir(path)
	if err != nil {
		log.Fatalf("An error has occurred: %v", err.Error())
	}

	err = cmdExec(prog, args[2:], envMap2Slice(envs))
	if err != nil {
		log.Fatalf("An error has occurred: %v", err.Error())
	}
}

func getEnvsFromDir(path string) (envs map[string]string, err error) {
	dir, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return envs, errors.Errorf("Can not open a directory %q", path)
	}
	defer dir.Close()
	filesInfo, err := dir.Readdir(maxNumberOfFiles)
	envs = make(map[string]string, len(filesInfo))

	for _, fileInfo := range filesInfo {
		filePath := path + string(os.PathSeparator) + fileInfo.Name()

		file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
		if err != nil {
			return envs, errors.Errorf("Can not open a file %q\n", filePath)
		}

		fileContent := make([]byte, fileInfo.Size())
		n, err := file.Read(fileContent)
		if err != nil && err != io.EOF {
			return envs, errors.Errorf("Can not read a file %q, an error has occurred: %v\n", filePath, err.Error())
		}
		if int64(n) != fileInfo.Size() {
			return envs, errors.Errorf("Can not read a full file %q, read bites: %v; expected: %v\n", filePath, n, fileInfo.Size())
		}
		file.Close()

		envs[fileInfo.Name()] = string(fileContent)
	}
	return envs, nil
}


func setEnvs(envs map[string]string) {

	for name, val := range envs {

		if val != "" {
			os.Setenv(name, val)
		} else {
			os.Unsetenv(name)
		}
	}
}


func envMap2Slice(envsMap map[string]string) (envsSlice []string) {
	envsSlice = make([]string, 0, len(envsMap))

	for k, v := range envsMap {
		envsSlice = append(envsSlice, k + "=" + v)
	}
	return envsSlice
}


func cmdExec(prog string, args []string, envs []string) error {

	cmd := exec.Command(prog, args...)
	cmd.Env = append(os.Environ(), envs...)

	dir, err := os.Getwd()
	if err != nil {
		return errors.Wrapf(err, "Can not get the current dir")
	}
	cmd.Dir = dir

	cmd.Stdout  = os.Stdout
	cmd.Stderr  = os.Stderr
	cmd.Stdin   = os.Stdin

	err = cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "Command run error, command: %v, envs: %v", prog, envs)
	}

	return nil
}
