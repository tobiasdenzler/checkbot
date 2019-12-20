package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"

	log "github.com/sirupsen/logrus"
)

// Sandbox can test check scripts
type Sandbox struct {
	Script string
	Result string
	Error  string
}

// Run the check and return the result.
func (app *application) runSandbox(script string) *Sandbox {

	log.Debug("Execute sandbox script")

	sandbox := new(Sandbox)

	// Remove Windows or other CF/LF characters
	re := regexp.MustCompile(`\r?\n`)
	sandbox.Script = re.ReplaceAllString(script, "\n")

	// Write sandbox script to file
	data := []byte(sandbox.Script)
	err := ioutil.WriteFile(os.TempDir()+"/sandbox.sh", data, 0644)

	defer func() {
		os.Remove(os.TempDir() + "/sandbox.sh")
	}()

	if err != nil {
		log.Warnf("Error creating file: %v", err)
		sandbox.Error = err.Error()
	}

	// Execute sandbox script
	cmd := exec.Command(determineBash(), os.TempDir()+"/sandbox.sh")
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	errRun := cmd.Run()

	if errRun != nil {
		log.Warnf("Error running sandbox script: %v", errRun)
		sandbox.Error = errRun.Error()
	} else {
		sandbox.Result = out.String()
		sandbox.Error = stderr.String()
	}

	return sandbox
}
