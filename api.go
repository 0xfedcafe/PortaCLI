package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"
)

// Fixing golang's hardcode
var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func newBuffer(fileContent []byte, fileParam string, fileName string, contentType string, params map[string]string) (*bytes.Buffer, *multipart.Writer, error) {
	buff := new(bytes.Buffer)
	writer := multipart.NewWriter(buff)
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fileParam), escapeQuotes(fileName)))
	h.Set("Content-Type", contentType)
	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, nil, err
	}
	_, err = part.Write(fileContent)
	if err != nil {
		return nil, nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, nil, err
	}
	return buff, writer, nil
}

func createPortrait(filename string, collage bool, video bool) (string, error) {
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	kind := "1"
	collageVal := ""
	if video {
		kind = "2"
	} else {
		if collage {
			collageVal = "1"
		} else {
			collageVal = "0"
		}
	}
	buff, writer, err := newBuffer(fileBytes, "image", "file.jpg",
		"image/jpeg", map[string]string{
			"wm":      "",
			"collage": collageVal,
			"id":      "",
			"no_crop": "",
			"code":    "",
			"type":    kind,
		})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "http://portraitplus.facefun.ai:8080/Port/MakePort", buff)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", fmt.Sprintf("%s; boundary=%s", writer.FormDataContentType(), writer.Boundary()))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	rBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	rText := string(rBytes)
	if !strings.HasSuffix(rText, ".mp4") && !strings.HasSuffix(rText, ".jpg") {
		return "", errors.New(rText)
	}
	return rText, nil
}

func downloadPortrait(filename string, portrait string) error {
	resp, err := http.Get(fmt.Sprintf("http://portraitplus.facefun.ai:8080/Port/%s", portrait))
	if err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, resp.Body)
	return err
}
