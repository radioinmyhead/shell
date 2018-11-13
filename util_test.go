package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShell(t *testing.T) {
	assert := assert.New(t)

	stdout, stderr, err := Shell(`grep -io ABC <(echo -ne '123abc!@#\n$' | grep -P '\d+')`)
	assert.Nil(err, "exit code")
	assert.Equal("", stderr)
	assert.Equal("abc\n", stdout)

	stdout, stderr, err = Shell(`echo 123|grep abc`)
	assert.NotNil(err, "exit code")
}

func TestOut(t *testing.T) {
	assert := assert.New(t)

	stdout, err := Out(`grep -io ABC <(echo -ne '123abc!@#\n$' | grep -P '\d+')`)

	assert.Nil(err, "exit code")
	assert.Equal("abc\n", stdout)
}
