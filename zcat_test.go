package shell

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnzip(t *testing.T) {
	assert := assert.New(t)

	stdout, err := Output(`echo 12345| gzip`)
	assert.Nil(err, "exit code")

	ret, err := Unzip(bytes.NewBufferString(stdout))
	assert.Nil(err, "unzip")
	assert.Equal("12345\n", string(ret))
}
