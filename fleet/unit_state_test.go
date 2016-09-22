package fleet

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"errors"
)

func TestClient_UnitStates(t *testing.T) {
	unitStatesJson := `{
	  "nextPageToken": "ZAACAA==",
	  "states": [
	    {
	      "hash": "cf96f618332feae08d2bd5f2544f182f463d49fe",
	      "machineID": "9f08c99f7d9a304499004fd01891b396",
	      "name": "service1.service",
	      "systemdActiveState": "active",
	      "systemdLoadState": "loaded",
	      "systemdSubState": "running"
	    },
	    {
	      "hash": "8fb47a221a8933cbcb36c843914649be4c874c05",
	      "machineID": "287822fbe7a3134a87ebdd94975e9248",
	      "name": "service2.service",
	      "systemdActiveState": "active",
	      "systemdLoadState": "loaded",
	      "systemdSubState": "running"
	    }
	  ]
	}`

	baseUrl, _ := url.ParseRequestURI("http://fleet.example.com:4001/")
	senderMock := &requestSenderMock{
		httpResponse: &http.Response{
			Status:     "OK",
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(unitStatesJson)),
		},
	}

	client := &Client{
		baseUrl:       baseUrl,
		requestSender: senderMock,
	}

	unitStates, err := client.UnitStates()

	if nil != err {
		t.Error("Error supposed to be nil")
	}

	if http.MethodGet != senderMock.httpRequest.Method {
		t.Error("Request header must be GET")
	}

	if "http://fleet.example.com:4001/fleet/v1/state" != senderMock.httpRequest.URL.String() {
		t.Error("Request URL must be http://fleet.example.com:4001/fleet/v1/state")
	}

	if "application/json" != senderMock.httpRequest.Header.Get("Content-Type") {
		t.Error("Content-Type header must be application/json")
	}

	expectedUnitStates := []UnitState{
		{
			Hash:               "cf96f618332feae08d2bd5f2544f182f463d49fe",
			MachineID:          "9f08c99f7d9a304499004fd01891b396",
			Name:               "service1.service",
			SystemdActiveState: "active",
			SystemdLoadState:   "loaded",
			SystemdSubState:    "running",
		}, {
			Hash:               "8fb47a221a8933cbcb36c843914649be4c874c05",
			MachineID:          "287822fbe7a3134a87ebdd94975e9248",
			Name:               "service2.service",
			SystemdActiveState: "active",
			SystemdLoadState:   "loaded",
			SystemdSubState:    "running",
		},
	}

	if len(unitStates) != 2 {
		t.Error("Return slice must contain exactly 2 items")
	}

	for key, expectedUnitState := range expectedUnitStates {
		if unitStates[key].Hash != expectedUnitState.Hash {
			t.Errorf("Wrong unit stat %d Hash", key)
		}
		if unitStates[key].MachineID != expectedUnitState.MachineID {
			t.Errorf("Wrong unit stat %d MachineID", key)
		}
		if unitStates[key].Name != expectedUnitState.Name {
			t.Errorf("Wrong unit stat %d Name", key)
		}
		if unitStates[key].SystemdActiveState != expectedUnitState.SystemdActiveState {
			t.Errorf("Wrong unit stat %d SystemdActiveState", key)
		}
		if unitStates[key].SystemdLoadState != expectedUnitState.SystemdLoadState {
			t.Errorf("Wrong unit stat %d SystemdLoadState", key)
		}
		if unitStates[key].SystemdSubState != expectedUnitState.SystemdSubState {
			t.Errorf("Wrong unit stat %d SystemdSubState", key)
		}
	}

}

func TestClient_UnitStateFiltered(t *testing.T) {
	unitStatesJson := `{
	  "states": []
	}`

	baseUrl, _ := url.ParseRequestURI("http://fleet.example.com:4001/")
	senderMock := &requestSenderMock{
		httpResponse: &http.Response{
			Status:     "OK",
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(unitStatesJson)),
		},
	}

	client := &Client{
		baseUrl:       baseUrl,
		requestSender: senderMock,
	}

	client.UnitStateFiltered(&UnitStateFilter{
		UnitName: "foo.service",
		MachineID: "287822fbe7a3134a87ebdd94975e9248",
	});

	if "287822fbe7a3134a87ebdd94975e9248" != senderMock.httpRequest.URL.Query().Get("machineID")  {
		t.Error("MachineId parameter is missing")
	}

	if "foo.service" != senderMock.httpRequest.URL.Query().Get("unitName")  {
		t.Error("UnitName parameter is missing")
	}

}

func TestClient_UnitStatesRequestError(t *testing.T) {
	baseUrl, _ := url.ParseRequestURI("http://fleet.example.com:4001/")
	senderMock := &requestSenderMock{
		err: errors.New("Request failed"),
	}

	client := &Client{
		baseUrl:       baseUrl,
		requestSender: senderMock,
	}

	machines, err := client.UnitStates()

	if err == nil {
		t.Error("Error supposed not to be nil")
	}

	if len(machines) > 0 {
		t.Error("Machines supposed to be empty")
	}

}

func TestClient_UnitStatesRequestWrongStatusError(t *testing.T) {
	baseUrl, _ := url.ParseRequestURI("http://fleet.example.com:4001/")
	senderMock := &requestSenderMock{
		httpResponse: &http.Response{
			Status:     "OK",
			StatusCode: 500,
			Body:       ioutil.NopCloser(bytes.NewBufferString("")),
		},
	}

	client := &Client{
		baseUrl:       baseUrl,
		requestSender: senderMock,
	}

	machines, err := client.UnitStates()

	if err == nil {
		t.Error("Error supposed not to be nil")
	}

	if len(machines) > 0 {
		t.Error("Machines supposed to be empty")
	}
}
