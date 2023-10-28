package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

// go ...
//    __________  ________ __                _
//   / ____/ __ \/ ____/ //_/   ____ _____  (_)
//  / /   / / / / __/ / ,<     / __ `/ __ \/ /
// / /___/ /_/ / /___/ /| |   / /_/ / /_/ / /
// \____/_____/_____/_/ |_|   \__,_/ .___/_/
//                                /_/

type CdekApi struct {
	baseurl  string
	token    string
	login    string
	password string
}

var cdekApi = CdekApi{}.Init()

// go for const
func (_cdekApi CdekApi) Init() CdekApi {
	_cdekApi.baseurl = "https://api.cdek.ru/v2"
	_cdekApi.login = env("CDEK_ACCOUNT", "#_PRIVATE_#")
	_cdekApi.password = env("CDEK_PASSWORD", "#_PRIVATE_#")
	return _cdekApi
}

// go get token
func (_cdekApi CdekApi) Auth() CdekApi {
	body := "grant_type=client_credentials&client_id=" + cdekApi.login + "&client_secret=" + cdekApi.password
	headers := make(map[string][]string)
	response := cdekApi.Send("/oauth/token?"+body, headers, "POST", []byte(""))
	result := cdekAuth{}
	err := json.Unmarshal(response, &result)
	if err != nil {
		panic(err.Error())
	}
	cdekApi.token = result.Token
	return cdekApi
}

// go api
func (_cdekApi CdekApi) Send(method string, headers map[string][]string, reqType string, body []byte) []byte {
	req, err := http.NewRequest(reqType, cdekApi.baseurl+method, bytes.NewBuffer(body))
	if err != nil {
		panic(err.Error())
	}
	req.Header = headers
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	return result
}

// go get available cities
func (_cdekApi CdekApi) cityList(size int, page int) []cdekCity {
	cityList := []cdekCity{}

	headers := map[string][]string{
		"Authorization": {"Bearer " + cdekApi.token},
	}
	req := "?country_codes=ru&size=" + strconv.Itoa(size) + "&page=" + strconv.Itoa(page)
	response := cdekApi.Send("/location/cities"+req, headers, "GET", []byte(""))
	err := json.Unmarshal(response, &cityList)
	if err != nil {
		panic(err.Error())
	}
	return cityList
}

// go get cdek regions
func (_cdekApi CdekApi) regionList() []cdekRegion {
	regionList := []cdekRegion{}

	headers := map[string][]string{
		"Authorization": {"Bearer " + cdekApi.token},
	}
	req := "?country_codes=RU&size=1000"
	response := cdekApi.Send("/location/regions"+req, headers, "GET", []byte(""))
	err := json.Unmarshal(response, &regionList)
	if err != nil {
		panic(err.Error())
	}
	return regionList
}
