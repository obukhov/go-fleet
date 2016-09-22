package fleet

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestClient_Destroy(t *testing.T) {
	baseUrl, _ := url.ParseRequestURI("http://fleet.example.com:4001/")
	senderMock := &requestSenderMock{
		httpResponse: &http.Response{
			Status:     "OK",
			StatusCode: 204,
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
		},
	}

	client := &Client{
		baseUrl:       baseUrl,
		requestSender: senderMock,
	}

	err := client.Destroy("foo.service")

	if nil != err {
		t.Error("Error supposed to be nil")
	}

	if http.MethodDelete != senderMock.httpRequest.Method {
		t.Error("Request header must be DELETE")
	}

	if "http://fleet.example.com:4001/fleet/v1/units/foo.service" != senderMock.httpRequest.URL.String() {
		t.Error("Request URL must be http://fleet.example.com:4001/fleet/v1/units/foo.service")
	}

	if "application/json" != senderMock.httpRequest.Header.Get("Content-Type") {
		t.Error("Content-Type header must be application/json")
	}
}

func TestClient_Destroy_NotFound(t *testing.T) {
	baseUrl, _ := url.ParseRequestURI("http://fleet.example.com:4001/")
	senderMock := &requestSenderMock{
		httpResponse: &http.Response{
			Status:     "Not found",
			StatusCode: 404,
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
		},
	}

	client := &Client{
		baseUrl:       baseUrl,
		requestSender: senderMock,
	}

	err := client.Destroy("foo.service")

	if nil == err {
		t.Error("Error supposed not to be nil")
	}

	if nil != err && ERROR_UNIT_NOT_FOUND != err.Error() {
		t.Errorf("Error supposed to be %s", ERROR_UNIT_NOT_FOUND)
	}
}

func TestClient_Destroy_InternalError(t *testing.T) {
	baseUrl, _ := url.ParseRequestURI("http://fleet.example.com:4001/")
	senderMock := &requestSenderMock{
		httpResponse: &http.Response{
			Status:     "Internal server error",
			StatusCode: 500,
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
		},
	}

	client := &Client{
		baseUrl:       baseUrl,
		requestSender: senderMock,
	}

	err := client.Destroy("foo.service")

	if nil == err {
		t.Error("Error supposed not to be nil")
	}
}

func TestClient_Destroy_RequestError(t *testing.T) {
	baseUrl, _ := url.ParseRequestURI("http://fleet.example.com:4001/")
	senderMock := &requestSenderMock{
		err: errors.New("Request error"),
	}

	client := &Client{
		baseUrl:       baseUrl,
		requestSender: senderMock,
	}

	err := client.Destroy("foo.service")

	if nil == err {
		t.Error("Error supposed not to be nil")
	}
}
