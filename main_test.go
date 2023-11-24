package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/cocoza4/data_microservices/models"
	"github.com/stretchr/testify/assert"
)

var PORT = 8080
var URI = fmt.Sprintf("http://app-test:%d", PORT)

var (
	test_user       = os.Getenv("SECRET_USER")
	test_password   = os.Getenv("SECRET_PASSWORD")
	test_secret_key = os.Getenv("SECRET_KEY")
)

func request(url, requestType string, data []byte) (*http.Response, []byte) {
	client := &http.Client{
		Transport: &http.Transport{},
	}

	req, err := http.NewRequest(requestType, url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err.Error())
	}

	req.SetBasicAuth(test_user, test_password)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp, body
}

func init() {
	data1 := []byte(`{
		"name": "p1",
		"description": "desc1"
	}`)
	data2 := []byte(`{
		"name": "p2",
		"description": "desc2"
	}`)
	data3 := []byte(`{
		"name": "p3",
		"description": "desc3"
	}`)

	url := fmt.Sprintf("%v/v1/products", URI)
	request(url, http.MethodPost, []byte(data1))
	request(url, http.MethodPost, []byte(data2))
	request(url, http.MethodPost, []byte(data3))
}

func TestGetProducts(t *testing.T) {
	url := fmt.Sprintf("%v/v1/products", URI)
	resp, body := request(url, http.MethodGet, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var products []models.Product
	err = json.Unmarshal(body, &products)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(products))
}

func TestGetProduct(t *testing.T) {
	url := fmt.Sprintf("%v/v1/products/%v", URI, "p1")
	resp, body := request(url, http.MethodGet, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var product models.Product
	err = json.Unmarshal(body, &product)
	assert.Nil(t, err)
	fmt.Println(product)
	assert.Equal(t, "p1", product.Name)
}

func TestGetLatest(t *testing.T) {
	url := fmt.Sprintf("%v/v1/products/latest", URI)
	resp, body := request(url, http.MethodGet, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var product models.Product
	err = json.Unmarshal(body, &product)
	assert.Nil(t, err)
	assert.Equal(t, "p3", product.Name)
}

func TestCreateProduct(t *testing.T) {
	data := []byte(`{
		"name": "mock",
		"description": "mock"
	}`)

	url := fmt.Sprintf("%v/v1/products", URI)
	resp, body := request(url, http.MethodPost, []byte(data))
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, `{"message":"success"}`, string(body))
}

func TestGetIndexes(t *testing.T) {
	url := fmt.Sprintf("%v/v1/products/indexes", URI)
	resp, body := request(url, http.MethodGet, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, `[{"name":"_id_"}]`, string(body))
}

func TestCreateIndex(t *testing.T) {
	url := fmt.Sprintf("%v/v1/products/indexes?field=name", URI)
	resp, body := request(url, http.MethodPost, nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, `{"message":"success"}`, string(body))
}

func TestCreateIndex_noFieldSpecified(t *testing.T) {
	url := fmt.Sprintf("%v/v1/products/indexes", URI)
	resp, body := request(url, http.MethodPost, nil)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"message":"'field' can't be empty"}`, string(body))
}

func TestRequest_noCredentials(t *testing.T) {
	url := fmt.Sprintf("%v/v1/products", URI)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error", err.Error())
	}

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.Equal(t, `{"message":"Authentication required"}`, string(body))
}

func TestRequest_invalidCredentials(t *testing.T) {
	url := fmt.Sprintf("%v/v1/products", URI)
	client := &http.Client{
		Transport: &http.Transport{},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	req.SetBasicAuth("xxx", "yyy") // invalid credentials

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.Equal(t, `{"message":"Invalid credentials"}`, string(body))
}

func TestLogin_invalidCredentials(t *testing.T) {
	url := fmt.Sprintf("%v/v1/login", URI)
	client := &http.Client{
		Transport: &http.Transport{},
	}

	data := []byte(`{
		"username": "xxx",
		"password": "yyy"
	}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"message":"user or password is incorrect"}`, string(body))
}

func TestLogin(t *testing.T) {
	data := []byte(fmt.Sprintf(`{
		"username": "%s",
		"password": "%s"
	}`, test_user, test_password))

	url := fmt.Sprintf("%v/v1/login", URI)
	resp, body := request(url, http.MethodPost, []byte(data))
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var respObj map[string]string
	json.Unmarshal([]byte(body), &respObj)
	_, ok := respObj["token"]
	assert.True(t, ok)
}
