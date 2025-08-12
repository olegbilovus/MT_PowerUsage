package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func checkClientNil(client *http.Client) error {
	if client == nil {
		return fmt.Errorf(NilClientError)
	}
	return nil
}

type Plug interface {
	Name() string
	Load() (float64, error)
	TurnOn() error
	TurnOff() error
}

const NilClientError = "the client is nil"

type ShellyPlugS struct {
	ip     string
	client *http.Client
}

func (s ShellyPlugS) Name() string {
	return "ShellyPlugS"
}

func (s ShellyPlugS) Load() (float64, error) {
	if err := checkClientNil(s.client); err != nil {
		return 0, err
	}

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%v/meter/0", s.ip), nil)
	res, err := s.client.Do(req)
	if err != nil {
		return 0, err
	}

	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}

	return data["power"].(float64), nil
}

func (s ShellyPlugS) power(on bool) error {
	url := fmt.Sprintf("http://%v/relay/0?turn=", s.ip)
	if on {
		url += "on"
	} else {
		url += "off"
	}

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	res, err := s.client.Do(req)
	if err != nil {
		return err
	}

	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	if data["ison"].(bool) != on {
		return fmt.Errorf("could not set plug to %v", on)
	}
	return nil
}

func (s ShellyPlugS) TurnOn() error {
	if err := checkClientNil(s.client); err != nil {
		return err
	}
	return s.power(true)
}

func (s ShellyPlugS) TurnOff() error {
	if err := checkClientNil(s.client); err != nil {
		return err
	}
	return s.power(false)
}

type ShellyPlugSv2 struct {
	ip     string
	client *http.Client
}

func (s ShellyPlugSv2) Name() string {
	return "ShellyPlugSv2"
}

func (s ShellyPlugSv2) Load() (float64, error) {
	if err := checkClientNil(s.client); err != nil {
		return 0, err
	}

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%v/rpc/Switch.GetStatus?id=0", s.ip), nil)
	res, err := s.client.Do(req)
	if err != nil {
		return 0, err
	}

	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}

	return data["apower"].(float64), nil
}

func (s ShellyPlugSv2) power(on bool) error {
	url := fmt.Sprintf("http://%v/rpc/Switch.Set?id=0&on=", s.ip)
	if on {
		url += "true"
	} else {
		url += "false"
	}

	req, _ := http.NewRequest(http.MethodGet, url, nil)
	res, err := s.client.Do(req)
	if err != nil {
		return err
	}

	body, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	var data map[string]any
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}

	return nil
}

func (s ShellyPlugSv2) TurnOn() error {
	if err := checkClientNil(s.client); err != nil {
		return err
	}
	return s.power(true)
}

func (s ShellyPlugSv2) TurnOff() error {
	if err := checkClientNil(s.client); err != nil {
		return err
	}
	return s.power(false)
}
