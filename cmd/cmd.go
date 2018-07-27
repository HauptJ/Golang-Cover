/*
DESC: Runs a bash command
Author: Joshua Haupt
Last Modified: 07-13-2018
*/

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

/*
DESC: runs a bash command
IN: command name: cmdName, command arguments cmdArgs
OUT: nill on success
SOURCE: https://nathanleclaire.com/blog/2014/12/29/shelled-out-commands-in-golang/
*/
func Run_cmd(cmdName string, cmdArgs []string) error {

	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for CMD", err)
		panic(err)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("command output | %s\n", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting CMD", err)
		panic(err)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for CMD", err)
		panic(err)
	}

	return nil
}
