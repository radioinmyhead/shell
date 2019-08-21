package shell

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func Curl(url string) (ret []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	ret, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}

func CurlZcat(url string) (ret []byte, err error) {
	b, err := Curl(url)
	if err != nil {
		return
	}

	ret, err := Unzip(bytes.NewBufferString(b))
	if err != nil {
		return
	}
	return
}
