/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
)

var (
	host = flag.String("host", "", "Remote host name or IP address")
	port = flag.Int("port", 22, "SSH port")
	user = flag.String("login", "", "Login User")
	pwd  = flag.String("pwd", "", "Login Password")
	cmd  = flag.String("cmd", "", "Command to run on remote server")
)

func main() {
	flag.Parse()

	config := &ssh.ClientConfig{
		User: *user,
		Auth: []ssh.AuthMethod{ssh.Password(*pwd)},
	}

	fmt.Println("-> Executing Command: ", cmd)
	report := executeCmd(*cmd, *host+":"+string(*port), config)
	fmt.Println("#### STDOUT ####")
	fmt.Println(report.stdout)
	fmt.Println("################")
	fmt.Println("-> ExitCode: ", report.exitCode)
	os.Exit(report.exitCode)
}

func executeCmd(cmd, hostname string, config *ssh.ClientConfig) ExecReport {
	conn, _ := ssh.Dial("tcp", hostname, config)
	session, _ := conn.NewSession()
	defer session.Close()

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	err := session.Run(cmd)

	var exitCode = 0
	if err != nil && err == err.(*ssh.ExitError) {
		exitCode = err.(*ssh.ExitError).ExitStatus()
	}

	return ExecReport{stdout: stdoutBuf.String(), stderr: stderrBuf.String(), exitCode: exitCode}
}

type ExecReport struct {
	stdout   string
	stderr   string
	exitCode int
}
