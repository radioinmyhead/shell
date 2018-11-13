package shell

import (
	"bytes"
	"fmt"
	"os/exec"
)

func ShellContext(ctx, cmd string) (stdout string, stderr string, err error) {
	c := exec.CommandContext(ctx, "bash", "-c", cmd)
	var o, e bytes.Buffer
	c.Stdout = &o
	c.Stderr = &e
	err = c.Run()
	return o.String(), e.String(), err
	return
}

func Shell(cmd string) (stdout string, stderr string, err error) {
	c := exec.Command("bash", "-c", cmd)
	var o, e bytes.Buffer
	c.Stdout = &o
	c.Stderr = &e
	err = c.Run()
	return o.String(), e.String(), err
}

func OutContext(ctx, cmd string) (str string, err error) {
	str, stderr, err := ShellContext(cmd)
	if err != nil {
		err = fmt.Errorf("%s\n%s\n%v", str, stderr, err)
		return
	}
	if stderr != "" {
		err = fmt.Errorf("%s\n%s", str, stderr)
		return
	}
	return
}

func Out(cmd string) (str string, err error) {
	str, stderr, err := Shell(cmd)
	if err != nil {
		err = fmt.Errorf("%s\n%s\n%v", str, stderr, err)
		return
	}
	if stderr != "" {
		err = fmt.Errorf("%s\n%s", str, stderr)
		return
	}
	return
}
