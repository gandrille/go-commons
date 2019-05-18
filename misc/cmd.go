package misc

import (
	"bytes"
	"io"
	"os"
	"os/exec"

	"github.com/gandrille/go-commons/result"
)

// RunCmd executes a command with stdin/out/err piped from/to os defaults.
func RunCmd(cmd *exec.Cmd, displayName string) result.Result {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return result.Failure(displayName + ": " + err.Error())
	}

	return result.Success(displayName)
}

// RunCmdStdIn executes a command sending an input string to stdin.
// out/err are piped to os defaults.
func RunCmdStdIn(commandName, input string, cmd *exec.Cmd) result.Result {
	stdin, err := cmd.StdinPipe()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err != nil {
		return result.Failure(commandName + ": can't pipe stdin (" + err.Error() + ")")
	}
	err = cmd.Start()
	if err != nil {
		return result.Failure(commandName + ": starting error (" + err.Error() + ")")
	}

	io.Copy(stdin, bytes.NewBufferString(input))
	stdin.Close()

	err = cmd.Wait()
	if err != nil {
		return result.Failure(commandName + ": failed (" + err.Error() + ")")
	}
	return result.Success(commandName)
}
