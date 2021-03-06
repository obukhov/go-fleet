package fleet

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestClient_Units(t *testing.T) {
	unitsJson := `{
	  "units": [
	    {
	      "currentState": "loaded",
	      "desiredState": "launched",
	      "machineID": "9f08c99f7d9a304499004fd01891b396",
	      "name": "foo.service",
	      "options": [
		{
		  "name": "After",
		  "section": "Unit",
		  "value": "docker.service"
		},
		{
		  "name": "Restart",
		  "section": "Service",
		  "value": "on-failure"
		},
		{
		  "name": "RestartSec",
		  "section": "Service",
		  "value": "25s"
		},
		{
		  "name": "ExecStart",
		  "section": "Service",
		  "value": "/bin/bash -c \"while true; do echo 'hello' && sleep 1; done\""
		}
	      ]
	    }
	  ]
	}`

	baseUrl, _ := url.ParseRequestURI("http://fleet.example.com:4001/")
	senderMock := &requestSenderMock{
		httpResponse: &http.Response{
			Status:     "OK",
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(unitsJson)),
		},
	}

	client := &Client{
		baseUrl:       baseUrl,
		requestSender: senderMock,
	}

	units, err := client.Units()

	if nil != err {
		t.Error("Error supposed to be nil")
	}

	if http.MethodGet != senderMock.httpRequest.Method {
		t.Error("Request header must be GET")
	}

	if "http://fleet.example.com:4001/fleet/v1/units" != senderMock.httpRequest.URL.String() {
		t.Error("Request URL must be http://fleet.example.com:4001/fleet/v1/units")
	}

	if "application/json" != senderMock.httpRequest.Header.Get("Content-Type") {
		t.Error("Content-Type header must be application/json")
	}

	if len(units) != 1 {
		t.Error("Return slice must contain exactly 1 items")
	}

	validateUnit(&units[0], t)
}

func TestClient_Unit(t *testing.T) {
	unitsJson := `{
	  "currentState": "loaded",
  	  "desiredState": "launched",
	  "machineID": "9f08c99f7d9a304499004fd01891b396",
	  "name": "foo.service",
	  "options": [
	    {
	      "name": "After",
	      "section": "Unit",
	      "value": "docker.service"
	    },
	    {
	      "name": "Restart",
	      "section": "Service",
	      "value": "on-failure"
	    },
	    {
	      "name": "RestartSec",
	      "section": "Service",
	      "value": "25s"
	    },
	    {
	      "name": "ExecStart",
	      "section": "Service",
	      "value": "/bin/bash -c \"while true; do echo 'hello' && sleep 1; done\""
	    }
	  ]
	}`

	baseUrl, _ := url.ParseRequestURI("http://fleet.example.com:4001/")
	senderMock := &requestSenderMock{
		httpResponse: &http.Response{
			Status:     "OK",
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(unitsJson)),
		},
	}

	client := &Client{
		baseUrl:       baseUrl,
		requestSender: senderMock,
	}

	unit, err := client.Unit("foo.service")

	if nil != err {
		t.Error("Error supposed to be nil")
	}

	if http.MethodGet != senderMock.httpRequest.Method {
		t.Error("Request header must be GET")
	}

	if "http://fleet.example.com:4001/fleet/v1/units/foo.service" != senderMock.httpRequest.URL.String() {
		t.Error("Request URL must be http://fleet.example.com:4001/fleet/v1/units/foo.service")
	}

	if "application/json" != senderMock.httpRequest.Header.Get("Content-Type") {
		t.Error("Content-Type header must be application/json")
	}

	validateUnit(unit, t)
}

func validateUnit(unit *Unit, t *testing.T) {
	if unit.CurrentState != "loaded" {
		t.Error("Wrong unit current state")
	}

	if unit.DesiredState != "launched" {
		t.Error("Wrong unit desired state")
	}

	if unit.MachineID != "9f08c99f7d9a304499004fd01891b396" {
		t.Error("Wrong unit machinedID")
	}

	if unit.Name != "foo.service" {
		t.Error("Wrong unit name")
	}

	if len(unit.Options) != 4 {
		t.Error("Options slice must contain exactly 4 items")
	}
	expectedOptions := []UnitOption{
		{
			Name:    "After",
			Section: "Unit",
			Value:   "docker.service",
		},
		{
			Name:    "Restart",
			Section: "Service",
			Value:   "on-failure",
		},
		{
			Name:    "RestartSec",
			Section: "Service",
			Value:   "25s",
		},
		{
			Name:    "ExecStart",
			Section: "Service",
			Value:   "/bin/bash -c \"while true; do echo 'hello' && sleep 1; done\"",
		},
	}

	for key, expectedOption := range expectedOptions {
		if unit.Options[key].Name != expectedOption.Name {
			t.Errorf("Wrong option %d name", key)
		}

		if unit.Options[key].Section != expectedOption.Section {
			t.Errorf("Wrong option %d section", key)
		}

		if unit.Options[key].Value != expectedOption.Value {
			t.Errorf("Wrong option %d value", key)
		}
	}
}

