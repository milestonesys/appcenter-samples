package vms

import (
	"encoding/json"
	"errors"
)

type Camera struct {
	ID           string `json:"id"`
	Name         string `json:"displayName"`
	Description  string `json:"description"`
	Enabled      bool   `json:"enabled"`
	LastModified string `json:"lastModified"`
	Channel      int    `json:"channel"`
}

func (c *Camera) ToString() string {
	return c.Name
}

type CamerasList struct {
	Cameras []*Camera
}

func (cl *CamerasList) ToJSON() (string, error) {
	if cl == nil {
		return "", errors.New("CamerasList is nil")
	}
	jsonData, err := json.Marshal(cl.Cameras)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func NewCamerasList() *CamerasList {
	return &CamerasList{
		Cameras: []*Camera{},
	}
}

func (cl *CamerasList) Add(c *Camera) {
	if c != nil {
		cl.Cameras = append(cl.Cameras, c)
	}
}

func (cl *CamerasList) Empty() bool {
	return len(cl.Cameras) == 0
}

func (cl *CamerasList) NotEmpty() bool {
	return !cl.Empty()
}

func (cl *CamerasList) ToString() []string {
	var s []string
	for _, item := range cl.Cameras {
		s = append(s, item.ToString())
	}
	return s
}
