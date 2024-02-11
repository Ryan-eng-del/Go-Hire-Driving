package biz

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-kratos/kratos/v2/log"
)

// GreeterUsecase is a Greeter usecase.
type MapServiceBiz struct {
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewMapServiceBiz(logger log.Logger) *MapServiceBiz {
	return &MapServiceBiz{log: log.NewHelper(logger)}
}

type DrivePaths struct {
	Distance string `json:"distance`
	Duration string `json:"duration`
}

type DriveServiceRes struct {
	Status string `json:"status"`
	Info string `json:"info"`
	Route struct {
		Origin string `json:"origin"`
		Destination string `json:"destination"`
		Paths  []DrivePaths `json:"paths"`
	}
}

func (s *MapServiceBiz) GetDriveInfo(origin string, destination string) ( string, string, error) {
	key := "96bb5f00921daa6fdaaf6c85bc1c9dea"
	api := "https://restapi.amap.com/v3/direction/driving"
	parameters := fmt.Sprintf("?origin=%s&destination=%s&extensions=base&key=%s", origin, destination, key)
	resp, err := http.Get(api + parameters)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	driveInfoResp := DriveServiceRes{}
	err = json.Unmarshal(body, &driveInfoResp)
	if err != nil {
		return "", "", err
	}
	firstPath := driveInfoResp.Route.Paths[0]
	return firstPath.Distance, firstPath.Duration, nil
}