package shell

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"reflect"
)

type Process struct{}

// run command as bash
func Bash(ctx context.Context, cmd string) (ret string, err error) {
	var o, e bytes.Buffer
	var se string

	parms := []string{"-o0", "-e0", "bash", "-c", cmd}

	c := exec.CommandContext(ctx, "stdbuf", parms...)
	c.Stdout = &o
	c.Stderr = &e

	err = c.Run()
	if err != nil {
		err = fmt.Errorf("run shell: %s: %s: %s", cmd, err, e.String())
		return
	}

	se = e.String()
	if se != "" {
		err = fmt.Errorf("run shell: %s: %s", cmd, se)
		return
	}

	ret = o.String()
	return
}

func Run(ctx context.Context, oricmd []string) (ch chan string, err error) {
	// refer to: https://unix.stackexchange.com/questions/25372/turn-off-buffering-in-pipe
	// script run command in a pseudo terminal (pty). man script for help
	parms := append([]string{"-o0", "-e0"}, oricmd...)
	cmd := exec.CommandContext(ctx, "stdbuf", parms...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		err = fmt.Errorf("shell run: %v", err)
		return
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Start()
	if err != nil {
		err = fmt.Errorf("shell run: %v: %v", err, stderr.String())
		return
	}

	ch = make(chan string)
	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(stdout)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			m := scanner.Text()
			ch <- m
		}
		cmd.Wait()
	}()
	return
}

func Output(cmd string) (ret string, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ret, err = Bash(ctx, cmd)
	if err != nil {
		err = fmt.Errorf("shell output: %v", err)
		return
	}
	return
}

func (p *Process) SetField(obj interface{}, name string, value interface{}) error {
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
