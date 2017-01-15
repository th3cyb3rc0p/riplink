package requests_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/mschwager/riplink/src/requests"
)

type MockClient struct {
	Body string
	Code int
	Err  error
}

func (client MockClient) Do(request *http.Request) (response *http.Response, err error) {
	response = &http.Response{
		StatusCode: client.Code,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(client.Body))),
	}

	return response, client.Err
}

func TestDoBasic(t *testing.T) {
	body := ""
	code := 200

	client := MockClient{
		Body: body,
		Code: code,
	}

	result_body, result_code, result_err := requests.Request(client, "UNUSED", "UNUSED", nil)

	if result_body != body || result_code != code || result_err != nil {
		t.Error("Failed to parse request: ", client)
	}
}

func TestDoError(t *testing.T) {
	body := ""
	code := 0
	err := errors.New("")

	client := MockClient{
		Body: body,
		Code: code,
		Err:  err,
	}

	result_body, result_code, result_err := requests.Request(client, "UNUSED", "UNUSED", nil)

	if result_body != body || result_code != code || result_err != err {
		t.Error("Failed to parse request: ", client)
	}
}