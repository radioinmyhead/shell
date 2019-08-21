package shell

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func Zcat(ctx context.Context, filename string) ([]byte, error) {
	fi, err := os.Open(filename)
	if err != nil {
		err = fmt.Errorf("in zcat: %v", err)
		return nil, err
	}
	defer fi.Close()

	return Unzip(fi)
}

func Unzip(fi io.Reader) ([]byte, error) {
	fz, err := gzip.NewReader(fi)
	if err != nil {
		err = fmt.Errorf("in unzip: %v", err)
		return nil, err
	}
	defer fz.Close()

	s, err := ioutil.ReadAll(fz)
	if err != nil {
		err = fmt.Errorf("in unzip: %v", err)
		return nil, err
	}
	return s, nil
}
