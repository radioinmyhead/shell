package shell

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Shell(cmd string) (stdout string, stderr string, err error) {
	c := exec.Command("bash", "-c", cmd)
	var o, e bytes.Buffer
	c.Stdout = &o
	c.Stderr = &e
	err = c.Run()
	return o.String(), e.String(), err
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
