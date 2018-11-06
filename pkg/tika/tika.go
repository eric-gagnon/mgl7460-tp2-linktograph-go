package tika

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func disabledGetTikaContentAndMetaForFile(filepath string) {

	// Open the file
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalln(err)
	}
	// Close the file later
	defer file.Close()

	// Buffer to store our request body as bytes
	var requestBody bytes.Buffer

	// Create a multipart writer
	multiPartWriter := multipart.NewWriter(&requestBody)

	// Initialize the file field
	fileWriter, err := multiPartWriter.CreateFormFile("file_field", "name")
	if err != nil {
		log.Fatalln(err)
	}

	// Copy the actual file content to the field field's writer
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		log.Fatalln(err)
	}

	// Populate other fields
	fieldWriter, err := multiPartWriter.CreateFormField("normal_field")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = fieldWriter.Write([]byte("Value"))
	if err != nil {
		log.Fatalln(err)
	}

	// We completed adding the file and the fields, let's close the multipart writer
	// So it writes the ending boundary
	multiPartWriter.Close()

	// By now our original request body should have been populated, so let's just use it with our custom request
	url := "http://localhost:9998/rmeta/text"
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		log.Fatalln(err)
	}
	// We need to set the content type from the writer, it includes necessary boundary as well
	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())
	req.Header.Add("Accept", "application/json")

	// Do the request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}

	json.NewDecoder(response.Body).Decode(&result)

	log.Println(result)
}

func GetTikaContentAndMetaForFile(filepath string) {
	// Voir https://wiki.apache.org/tika/TikaJAXRS
	// curl -H "Accept: application/json" -T 0fd9a51e2667aa7244f81b7ad68841297d2fe82f http://localhost:9998/rmeta/text > meta.json
	// Exemples :
	// http://www.codershood.info/2017/06/25/http-curl-request-golang/.
	// https://gist.github.com/ebraminio/576fdfdff425bf3335b51a191a65dbdb
	// https://stackoverflow.com/questions/20205796/golang-post-data-using-the-content-type-multipart-form-data
	// https://stackoverflow.com/questions/39761910/how-can-you-upload-files-as-a-stream-in-go
	// http://polyglot.ninja/golang-making-http-requests/
	// Scrapfiletocache(url, tika cache folder).

	fmt.Println("test")

	// todo: put this in configs.
	url := "http://localhost:9998/rmeta/text"

	filepayload, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer filepayload.Close()

	req, _ := http.NewRequest("PUT", url, filepayload)
	req.Header.Add("Accept", "application/json")


	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var result map[string]interface{}

	json.NewDecoder(res.Body).Decode(&result)

	log.Println(result)

	fmt.Println(string(body))
}

func GetPersons() {
	// ...
}
