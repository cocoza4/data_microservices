package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/cocoza4/data_microservices/models"
	"github.com/stretchr/testify/assert"
)

var PORT = 8080
var URI = fmt.Sprintf("http://app-test:%d", PORT)

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
	http.Post(url, "application/json", bytes.NewBuffer(data1))
	http.Post(url, "application/json", bytes.NewBuffer(data2))
	http.Post(url, "application/json", bytes.NewBuffer(data3))

}

func TestGetProducts(t *testing.T) {
	url := fmt.Sprintf("%v/v1/products", URI)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error", err.Error())
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	var products []models.Product
	err = json.Unmarshal(body, &products)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(products))
}

func TestGetProduct(t *testing.T) {
	url := fmt.Sprintf("%v/v1/products/%v", URI, "p1")
	resp, _ := http.Get(url)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	var product models.Product
	err = json.Unmarshal(body, &product)
	assert.Nil(t, err)
	fmt.Println(product)
	assert.Equal(t, "p1", product.Name)
}

func TestGetLatest(t *testing.T) {
	url := fmt.Sprintf("%v/v1/products/latest", URI)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error", err.Error())
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

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
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(data))
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	assert.Equal(t, `{"message":"success"}`, string(body))
}

func TestGetIndexes(t *testing.T) {
	url := fmt.Sprintf("%v/v1/products/indexes", URI)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error", err.Error())
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	assert.Equal(t, `[{"name":"_id_"}]`, string(body))
}

func TestCreateIndex(t *testing.T) {
	url := fmt.Sprintf("%v/v1/products/indexes?field=name", URI)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte("")))
	if err != nil {
		log.Println("Error", err.Error())
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	assert.Equal(t, `{"message":"success"}`, string(body))
}

func TestCreateIndex_noFieldSpecified(t *testing.T) {
	url := fmt.Sprintf("%v/v1/products/indexes", URI)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte("")))
	if err != nil {
		log.Println("Error", err.Error())
	}
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	assert.Equal(t, `{"message":"'field' can't be empty"}`, string(body))
}
