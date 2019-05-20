package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.raad.cloud/cloud/hermes/config"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
	)

type Service struct {
	BitwiseAbilities int64 `json:"bitwise_abilities,omitempty" structs:"bitwise_abilities,omitempty"`
	IsEnabled        bool  `json:"is_enabled,omitempty" structs:"is_enabled,omitempty"`

	ModifiedAt *time.Time `json:"modified_at,omitempty" structs:"modified_at,omitempty"`
}

type Application struct {
	ID           string              `json:"_id,omitempty" structs:"_id,omitempty"`
	Name         string              `json:"name,omitempty" structs:"name,omitempty"`
	ClientID     string              `json:"client_id,omitempty" structs:"client_id,omitempty"`
	ClientSecret string              `json:"client_secret,omitempty" structs:"client_secret,omitempty"`
	Services     map[string]*Service `json:"services,omitempty" structs:"services,omitempty"`
	IsEnabled    bool                `json:"is_enabled,omitempty" structs:"is_enabled,omitempty"`

	CreatedAt  time.Time  `json:"created_at,omitempty" structs:"created_at,omitempty"`
	ModifiedAt *time.Time `json:"modified_at,omitempty" structs:"modified_at,omitempty"`

	LoadedAt time.Time
}

var applications = make(map[string]*Application)

const (
	APPLICATION_CACHE_TIME_SECONDS = 60 * 60
)

// GetApplicationInfo gets application info with given id
func GetApplicationInfo(applicationID string) (*Application, error) {
	if application, ok := applications[applicationID]; ok {
		duration := time.Until(application.CreatedAt)
		if duration < APPLICATION_CACHE_TIME_SECONDS {
			return application, nil
		}
	}
	fmt.Println("Get Applications")
	url := "https://api.paygear.ir/application/v3" + "/applications/" + applicationID + "?services=true"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "bearer "+config.AuthToken)
	req.Header.Set("api-key", "5aa7e856ae7fbc00016ac5a01c65909797d94a16a279f46a4abb5faa")
	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println("Get Applications")
	fmt.Println(url)
	fmt.Println(resp)
	if err != nil {
		logrus.Errorf("error on client.Do to get application info. \n error is : %s\n ", err)
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		application := Application{}
		json.Unmarshal(body, &application)
		application.LoadedAt = time.Now()
		application.ClientSecret = config.ClientSecret
		applications[applicationID] = &application
		return &application, nil
	}
	logrus.Errorf("app-service returned invalid status code %d\n", resp.StatusCode)
	return nil, errors.New("app-service returned invalid status code")
}

// GetApplicationAuthKeys gets client id & client secret of application
func GetApplicationAuthKeys(applicationID string) (string, string, error) {
	application, err := GetApplicationInfo(applicationID)
	if err != nil {
		logrus.Errorf("error on GetApplicationAuthKeys on get application info. \n error is : %s\n ", err)
		return "", "", err
	}
	return application.ClientID, application.ClientSecret, nil
}

// GetApplicationPermissions gets information of permissions of application
func GetApplicationPermissions(applicationID string) (bool, map[string]*Service, error) {
	application, err := GetApplicationInfo(applicationID)
	if err != nil {
		logrus.Errorf("error on GetApplicationPermissions on get application info. \n error is : %s\n ", err)
		return false, nil, err
	}
	return application.IsEnabled, application.Services, nil
}
