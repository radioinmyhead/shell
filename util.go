package shell

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

func ShellContext(ctx context.Context, cmd string) (stdout string, stderr string, err error) {
	c := exec.CommandContext(ctx, "bash", "-c", cmd)
	c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	go func() {
		<-ctx.Done()
		syscall.Kill(-c.Process.Pid, syscall.SIGKILL)
	}()
	var o, e bytes.Buffer
	c.Stdout = &o
	c.Stderr = &e
	err = c.Run()
	return o.String(), e.String(), err
}

func Shell(cmd string) (stdout string, stderr string, err error) {
	c := exec.Command("bash", "-c", cmd)
	var o, e bytes.Buffer
	c.Stdout = &o
	c.Stderr = &e
	err = c.Run()
	return o.String(), e.String(), err
}

func OutContext(ctx context.Context, cmd string) (str string, err error) {
	str, stderr, err := ShellContext(ctx, cmd)
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

////////////////////////////////////////////////////////////////
func runWithTimeout(cmd string, num int) *exec.Cmd {
	c := exec.Command("bash", "-c", cmd)
	c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	time.AfterFunc(time.Duration(num)*time.Millisecond, func() {
		syscall.Kill(-c.Process.Pid, syscall.SIGKILL)
	})
	return c
}

func OutWithTimeout(cmd string, num int) (string, error) {
	c := runWithTimeout(cmd, num)

	out, err := c.Output()
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}
func CombinedOutWithTimeout(cmd string, num int) (string, error) {
	c := runWithTimeout(cmd, num)

	out, err := c.CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return string(out), err
}

////////////////////////////////////////////////////////
// ch,err:=shell.Run(ctx,cmd)
//
func Run(ctx context.Context, strcmd string) (ch chan string, err error) {

	cmd := exec.Command("stdbuf", "-o0", "-e0", "bash", "-c", strcmd)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	go func() {
		<-ctx.Done()
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}()

	out, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	err = cmd.Start()
	if err != nil {
		return
	}

	ch = make(chan string)
	go func() {
		defer close(ch)
		defer cmd.Wait()
		scanner := bufio.NewScanner(out)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
	}()
	return
}
