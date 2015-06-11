package gsafeeder

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func Upload(gsa, file string) {
    url :=  fmt.Sprintf("http://%s:19900/xmlfeed", gsa)
    h := getMetadata(file)
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	f, err := os.Open(file)
	check(err)

	fw, err := w.CreateFormFile("data", file)
	check(err)

	_, err = io.Copy(fw, f)
	check(err)

	fw, err = w.CreateFormField("datasource")
	check(err)

	_, err = fw.Write([]byte(h.Datasource))
	check(err)

	fw, err = w.CreateFormField("feedtype")
	check(err)

	_, err = fw.Write([]byte(h.Feedtype))
	check(err)

	w.Close()

    req, err := http.NewRequest("POST", url, &b)
	check(err)

	req.Header.Set("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	check(err)

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	} else {
		contents, err := ioutil.ReadAll(res.Body)
		check(err)
		fmt.Printf("Response: %v\n", string(contents))
	}
	return
}

type Header struct {
    Datasource string `xml:"datasource"`
    Feedtype string `xml:"feedtype"`
}

func getMetadata(inputFile string)(header Header){
    xmlFile, err := os.Open(inputFile)
    check(err)
    defer xmlFile.Close()

    var h Header
    decoder := xml.NewDecoder(xmlFile)
    var inElement string
    for {
        // Read tokens from the XML document in a stream.
        t, _ := decoder.Token()
        if t == nil {
            break
        }
        // Inspect the type of the token just read.
        switch se := t.(type) {
        case xml.StartElement:
            // If we just read a StartElement token
            inElement = se.Name.Local
            // ...and its name is "page"
            if inElement == "header" {
                decoder.DecodeElement(&h, &se)
                break
                }
        default:
        }
    }
    return h
}
