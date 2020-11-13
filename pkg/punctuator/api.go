package punctuator

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func Punctuate(text string) (string, error) {
	var punctuated string
	reply, err := http.PostForm("http://bark.phon.ioc.ee/punctuator", url.Values{"text": {text}})
	defer reply.Body.Close()

	if err != nil {
		return punctuated, fmt.Errorf("error doing POST to punctuator: %w", err)
	}

	body, err := ioutil.ReadAll(reply.Body)
	if err != nil {
		return punctuated, fmt.Errorf("error reading reply from punctuator: %w", err)
	}

	punctuated = strings.TrimSpace(string(body))

	return punctuated, nil
}
