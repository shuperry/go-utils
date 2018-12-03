package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RequestUtil struct{}

func (ru *RequestUtil) getClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true}, // 支持 https.
		DisableCompression: true,
		DisableKeepAlives:  true,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   15 * time.Second,
	}

	return client
}

func (ru *RequestUtil) Get(urlStr string) ([]byte, error) {
	log.Println("into request-util Get function with url:", urlStr)

	u, _ := url.Parse(urlStr)
	qs := u.Query()
	u.RawQuery = qs.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Panic("into request-util Get function with new request err:", err)
		return nil, err
	}

	client := ru.getClient()

	resp, err := client.Do(req)
	if err != nil {
		log.Panic("into request-util Get function with send request err:", err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Panic("into request-util Get function with read body from response err:", err)
		return nil, err
	}

	return body, nil
}

func (ru *RequestUtil) MultipartRequest(
	method string,
	urlStr string,
	header map[string]string,
	params map[string]string,
	queryString map[string]string,
	filePaths []string,
	fileFieldName string,
) ([]byte, error) {
	if method != "POST" && method != "PATCH" && method != "PUT" {
		log.Panic("request-util MultipartRequest method only support POST、PATCH、PUT")
	}

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	var fileWriter io.Writer
	var fd *os.File

	for _, filePath := range filePaths {
		fileWriter, _ = bodyWriter.CreateFormFile(fileFieldName, filepath.Base(filePath))
		fd, _ = os.Open(filePath)
		defer fd.Close()

		io.Copy(fileWriter, fd)
	}

	if len(params) != 0 {
		for k, v := range params {
			bodyWriter.WriteField(k, v)
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	req, err := http.NewRequest(method, urlStr, nil)
	if err != nil {
		log.Panic("into request-util MultipartRequest function with new request err:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	for k, v := range header {
		req.Header.Set(k, v)
	}

	values := req.URL.Query()
	for k, v := range queryString {
		values.Add(k, v)
	}
	req.URL.RawQuery = values.Encode()

	req.Body = ioutil.NopCloser(bodyBuf)

	client := ru.getClient()
	response, err := client.Do(req)
	if err != nil {
		log.Println("into request-util MultipartRequest function with send request err:", err)
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("into request-util MultipartRequest function with read body from response err:", err)
		return nil, err
	}

	return body, nil
}

func (ru *RequestUtil) FormRequest(
	method string,
	urlStr string,
	header map[string]string,
	params map[string][]string,
	queryString map[string]string,
) ([]byte, error) {
	if method != "POST" && method != "PATCH" && method != "PUT" {
		log.Panic("request-util MultipartRequest method must be one of POST、PATCH、PUT")
	}

	newParams := ""
	values := url.Values{}
	for k, v := range params {
		log.Println("param =", strings.Replace(strings.Trim(fmt.Sprint(v), "[]"), " ", ",", -1))

		values.Add(k, strings.Replace(strings.Trim(fmt.Sprint(v), "[]"), " ", ",", -1))
	}
	newParams = values.Encode()

	req, err := http.NewRequest(method, urlStr, strings.NewReader(newParams))
	if err != nil {
		log.Panic("into request-util Post function with new request err:", err)
		return nil, err
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	qsMap := req.URL.Query()
	for k, v := range queryString {
		qsMap.Add(k, v)
	}
	req.URL.RawQuery = qsMap.Encode()

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := ru.getClient()
	response, err := client.Do(req)
	if err != nil {
		log.Println("into request-util Post function with send request err:", err)
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("into request-util Post function with read body from response err:", err)
		return nil, err
	}

	return body, nil
}

func (ru *RequestUtil) JsonRequest(
	method string,
	urlStr string,
	header map[string]string,
	params map[string]interface{},
	queryString map[string]string,
) ([]byte, error) {
	if method != "POST" && method != "PATCH" && method != "PUT" {
		log.Panic("request-util MultipartRequest method only support POST、PATCH、PUT")
	}

	jsonByte, _ := json.Marshal(params)

	req, err := http.NewRequest(method, urlStr, bytes.NewBuffer(jsonByte))
	if err != nil {
		log.Panic("into request-util JsonRequest function with new request err:", err)
		return nil, err
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	qsMap := req.URL.Query()
	for k, v := range queryString {
		qsMap.Add(k, v)
	}
	req.URL.RawQuery = qsMap.Encode()

	req.Header.Set("Content-Type", "application/json")

	client := ru.getClient()
	response, err := client.Do(req)
	if err != nil {
		log.Println("into request-util JsonRequest function with send request err:", err)
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("into request-util JsonRequest function with read body from response err:", err)
		return nil, err
	}

	return body, nil
}
