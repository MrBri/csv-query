package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

type Response struct {
	Result []interface{} `json:"result"`
	Status string    `json:"status"`
	Time   string    `json:"time"`
}

var url = "http://csv:3333"

func main() {
	testUpload()
	testDateQuery()
	testLimitQuery()
	testWeatherTypeQuery()
	testWeatherTypeAndLimitQuery()
}

func testUpload() {
	fmt.Println("testing upload of csv...")

	// Open the CSV file
	file, err := os.Open("./seattle-weather.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Prepare a new form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the CSV file to the form
	part, err := writer.CreateFormFile("csvfile", "filename.csv")
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		panic(err)
	}

	// Send a POST request with the file
	req, err := http.NewRequest("POST", url + "/upload", body)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Print the response from the server
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBody))
	// res, err := http.Post()
}
func testDateQuery() {
	fmt.Println("testing date query...")
	res, err := http.Get(url + "/query?date=2012-02-06")
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println(res.Header["Content-Type"])
	fmt.Println(string(body))
}


func testLimitQuery() {
	fmt.Println("testing limit query...")
	res, err := http.Get(url + "/query?limit=5")
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}

	var resp []Response
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.Header["Content-Type"])
	fmt.Printf("response length: %d\n", len(resp[0].Result))
}

func testWeatherTypeQuery() {
	fmt.Println("testing weather type query...")
	res, err := http.Get(url + "/query?weather=rain")
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}

	var resp []Response
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.Header["Content-Type"])
	fmt.Printf("response length: %d\n", len(resp[0].Result))
}

func testWeatherTypeAndLimitQuery() {
	fmt.Println("testing weather type and limit query...")
	res, err := http.Get(url + "/query?weather=rain&limit=20")
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}

	var resp []Response
	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		panic(err)
	}
	fmt.Println(res.Header["Content-Type"])
	fmt.Printf("response length: %d\n", len(resp[0].Result))
}
