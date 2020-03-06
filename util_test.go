package shell

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShell(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	assert := assert.New(t)

	stdout, _, _, err := Bash(ctx, `grep -io ABC <(echo -ne '123abc!@#\n$' | grep -P '\d+')`)
	assert.Nil(err, "exit code")
	assert.Equal("abc\n", stdout)

	stdout, _, _, err = Bash(ctx, `echo 123|grep abc`)
	assert.NotNil(err, "exit code")
}

func TestOut(t *testing.T) {
	assert := assert.New(t)

	stdout, err := Output(`grep -io ABC <(echo -ne '123abc!@#\n$' | grep -P '\d+')`)

	assert.Nil(err, "exit code")
	assert.Equal("abc\n", stdout)
}

func TestRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	assert := assert.New(t)
	ch, code, err := Run(ctx, []string{"seq", "1", "3"})
	assert.Nil(err, "start command")
	ret := ""
	for line := range ch {
		ret += line
	}
	assert.Equal("123", ret)
	n := <-code
	assert.Equal(0, n)
}
