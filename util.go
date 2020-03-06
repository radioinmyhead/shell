package shell

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"syscall"
)

// run command as bash
func Bash(ctx context.Context, cmd string) (stdout, stderr string, code int, err error) {
	var o, e bytes.Buffer

	parms := []string{"-o0", "-e0", "bash", "-c", cmd}

	c := exec.CommandContext(ctx, "stdbuf", parms...)
	c.Stdout = &o
	c.Stderr = &e

	err = c.Run()
	if err != nil {
		err = fmt.Errorf("run shell: %s: %s: %s", cmd, err, e.String())
		return
	}

	stderr = e.String()
	stdout = o.String()
	code = c.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()

	return
}

func Run(ctx context.Context, oricmd []string) (ch chan string, code chan int, err error) {
	// refer to: https://unix.stackexchange.com/questions/25372/turn-off-buffering-in-pipe
	// script run command in a pseudo terminal (pty). man script for help
	parms := append([]string{"-o0", "-e0"}, oricmd...)
	cmd := exec.CommandContext(ctx, "stdbuf", parms...)

	pr, pw, err := os.Pipe()
	if err != nil {
		return
	}
	defer pw.Close()
	cmd.Stdout = pw
	cmd.Stderr = pw

	err = cmd.Start()
	if err != nil {
		return
	}

	ch = make(chan string)
	go func() {
		defer close(ch)
		defer close(code)
		defer pr.Close()
		scanner := bufio.NewScanner(pr)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			m := scanner.Text()
			ch <- m
		}
		cmd.Wait()
		code <- cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	}()
	return
}

func Output(cmd string) (ret string, err error) {
	var stderr string
	var code int
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ret, stderr, code, err = Bash(ctx, cmd)
	if err != nil {
		err = fmt.Errorf("shell output: %v", err)
		return
	}
	if code != 0 {
		ret += stderr
		err = fmt.Errorf("exit code=%v", code)
	}
	return
}

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}
