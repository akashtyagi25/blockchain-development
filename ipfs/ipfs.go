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

type ipfsresponse struct {
	Name string `json:"Name"`
	Hash string `json:"Hash"`
	Size string `json:"Size"`
}

func uploadtoipfs(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", file.Name())
	io.Copy(part, file)
	writer.Close()
	req, _ := http.NewRequest("POST", "http://localhost:5001/api/v0/add", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	respbody, _ := ioutil.ReadAll(resp.Body)
	var ipfsresp ipfsresponse
	json.Unmarshal(respbody, &ipfsresp)
	fmt.Println("file upload to ipfs!")
	fmt.Println("name : ", ipfsresp.Name)
	fmt.Println("hash (cid) : ", ipfsresp.Hash)
	fmt.Println("access url : https://ipfs.io/ipfs/" + ipfsresp.Hash)
}
func main() {
	uploadtoipfs("5thsemminibloodbank.pdf")
}
