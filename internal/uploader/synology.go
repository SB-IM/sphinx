// See official Synology API spec: https://global.download.synology.com/download/Document/Software/DeveloperGuide/Package/FileStation/All/enu/Synology_File_Station_API_Guide.pdf
package uploader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const (
	schema = "http"

	// API names.
	apiAuth     = "SYNO.API.Auth"
	apiDownload = "SYNO.FileStation.Download"
)

type ConfigOptions struct {
	hostname  string
	port      int
	Account   string
	Passwd    string
	StorePath string
}

func FileDownload(apiInfo *APIInfo, config ConfigOptions, files []string, sid string) (string, error) {
	api, ok := apiInfo.Data[apiDownload]
	if !ok {
		return "", errors.New("unsupported api name")
	}

	reqElem := requestElem{
		schema:  schema,
		host:    config.hostname + ":" + strconv.Itoa(config.port),
		apiName: apiDownload,
		version: api.MaxVersion,
		path:    api.Path,
		method:  "download",
		params:  fmt.Sprintf(`path=["%s"]&mode=download"`, strings.Join(files, `","`)),
		sid:     sid,
	}
	url, err := constructURL(reqElem)
	if err != nil {
		return "", err
	}
	fmt.Println(reqElem.params)

	resp, err := http.Get(url.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("http request failed, status code = %d", resp.StatusCode)
	}
	fmt.Println("content type:", resp.Header["Content-Type"])

	path := filepath.Join(config.StorePath + uuid.NewString())
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	fmt.Println("file saved in:", path)

	//Write the bytes to the file
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", err
	}

	return path, nil
}

func AuthLogin(apiInfo *APIInfo, config ConfigOptions) (string, error) {
	api, ok := apiInfo.Data[apiAuth]
	if !ok {
		return "", errors.New("unsupported api name")
	}

	reqElem := requestElem{
		schema:  schema,
		host:    config.hostname + ":" + strconv.Itoa(config.port),
		apiName: apiAuth,
		version: api.MaxVersion,
		path:    api.Path,
		method:  "login",
		params:  "account=" + config.Account + "&passwd=" + config.Passwd + "&session=FileStation&format=sid",
	}
	url, err := constructURL(reqElem)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(url.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("http request failed, status code = %d", resp.StatusCode)
	}

	var auth struct {
		Data struct {
			SID string `json:"sid,omitempty"`
		} `json:"data,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&auth); err != nil {
		return "", err
	}

	return auth.Data.SID, nil
}

func AuthLogout(apiInfo *APIInfo, config ConfigOptions, sid string) error {
	api, ok := apiInfo.Data[apiAuth]
	if !ok {
		return errors.New("unsupported api name")
	}

	reqElem := requestElem{
		schema:  schema,
		host:    config.hostname + ":" + strconv.Itoa(config.port),
		apiName: apiAuth,
		version: api.MaxVersion,
		path:    api.Path,
		method:  "logout",
		params:  "session=FileStation",
		sid:     sid,
	}
	url, err := constructURL(reqElem)
	if err != nil {
		return err
	}

	resp, err := http.Get(url.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("http request failed, status code = %d", resp.StatusCode)
	}

	return nil
}

// APIInfo is synology API infomation in the first step of workflow.
type APIInfo struct {
	Data map[string]APIDetail `json:"data,omitempty"`
}

type APIDetail struct {
	MaxVersion    int    `json:"maxVersion,omitempty"`
	MinVersion    int    `json:"minVersion,omitempty"`
	Path          string `json:"path,omitempty"`
	RequestFormat string `json:"requestFormat,omitempty"`
}

func GetAPIInfo(config ConfigOptions) (*APIInfo, error) {
	requestElem := requestElem{
		schema:  schema,
		host:    config.hostname + ":" + strconv.Itoa(config.port),
		apiName: "SYNO.API.Info",
		version: 1,
		path:    "query.cgi",
		method:  "query",
		params:  "query=" + apiAuth + "," + apiDownload,
	}
	url, err := constructURL(requestElem)
	if err != nil {
		return nil, err
	}
	return Info(url.String())
}

func Info(url string) (*APIInfo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http request failed, status code = %d", resp.StatusCode)
	}

	var apiInfo APIInfo
	if err := json.NewDecoder(resp.Body).Decode(&apiInfo); err != nil {
		return nil, err
	}

	return &apiInfo, nil
}

type requestElem struct {
	schema  string
	host    string
	apiName string
	version int
	path    string
	method  string
	params  string
	sid     string
}

// constructURL constructs a Synology API.
//
// Example:
//
// GET /webapi/<CGI_PATH>?api=<API_NAME>&version=<VERSION>&method=<METHOD>[&<PARAMS>][&_sid=<SID>]
//
// http://localhost:5000/webapi/query.cgi?api=SYNO.API.Info&version=1&method=query&query=all
func constructURL(reqElem requestElem) (*url.URL, error) {
	q, err := url.ParseQuery(reqElem.params)
	if err != nil {
		return nil, err
	}

	q.Set("api", reqElem.apiName)
	q.Set("version", strconv.Itoa(reqElem.version))
	q.Set("method", reqElem.method)
	if reqElem.sid != "" {
		q.Set("_sid", reqElem.sid)
	}

	var u url.URL
	u.Scheme = reqElem.schema
	u.Host = reqElem.host
	u.Path = "webapi/" + reqElem.path
	u.RawQuery = q.Encode()

	return &u, nil
}
