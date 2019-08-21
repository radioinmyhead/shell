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

	stdout, err := Bash(ctx, `grep -io ABC <(echo -ne '123abc!@#\n$' | grep -P '\d+')`)
	assert.Nil(err, "exit code")
	assert.Equal("abc\n", stdout)

	stdout, err = Bash(ctx, `echo 123|grep abc`)
	assert.NotNil(err, "exit code")
}

func TestOut(t *testing.T) {
	assert := assert.New(t)

	stdout, err := Output(`grep -io ABC <(echo -ne '123abc!@#\n$' | grep -P '\d+')`)

	assert.Nil(err, "exit code")
	assert.Equal("abc\n", stdout)
}
