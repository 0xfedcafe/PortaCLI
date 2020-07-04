package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"
)

var helloMessage string = ` 
                                 __                 ______   __        ______ 
                                /  |               /      \ /  |      /      |
  ______    ______    ______   _$$ |_     ______  /$$$$$$  |$$ |      $$$$$$/ 
 /      \  /      \  /      \ / $$   |   /      \ $$ |  $$/ $$ |        $$ |  
/$$$$$$  |/$$$$$$  |/$$$$$$  |$$$$$$/    $$$$$$  |$$ |      $$ |        $$ |  
$$ |  $$ |$$ |  $$ |$$ |  $$/   $$ | __  /    $$ |$$ |   __ $$ |        $$ |  
$$ |__$$ |$$ \__$$ |$$ |        $$ |/  |/$$$$$$$ |$$ \__/  |$$ |_____  _$$ |_ 
$$    $$/ $$    $$/ $$ |        $$  $$/ $$    $$ |$$    $$/ $$       |/ $$   |
$$$$$$$/   $$$$$$/  $$/          $$$$/   $$$$$$$/  $$$$$$/  $$$$$$$$/ $$$$$$/ 
$$ |                                                                          
$$ |                                                                          
$$/                          
Thank you for using PortaCLI.
Usage:
	-f = Path to your photo/video
	-c = (Y/n) Enables/disables collage mode     | Default: n
	-v = (Y/n) Select video/photo mode           | Default: n
	-o = Change output destination with filename | Default: img.jpg/vid.mp4 in your current folder
Example usages: 
		portacli -f /home/user/myphoto.jpg -c Y -o /home/user/outphoto.jpg
		portacli -f /home/user/myphoto.jpg
Authors:
	godande
	MaxUNof
Donate BTC: 36MQNEv8vkXgVuTa8HS1aJYzFTsuCmwNBK
`

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func newfileUploadRequest(uri string, params map[string]string, paramName, filepath string) (*http.Request, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	file.Close()
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(paramName), escapeQuotes("file.jpg")))
	h.Set("Content-Type", "image/jpeg")
	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	returnData, err := http.NewRequest("POST", uri, body)
	if returnData == nil {
		log.Fatal(err)
	}
	returnData.Header.Set("Content-Type", writer.FormDataContentType()+"; boundary="+writer.Boundary())
	return returnData, err
}

func main() {
	var filepath string
	var isCollage string
	var isVideo string
	var output string
	flag.StringVar(&filepath, "f", "", "Choose file using filepath\n")
	flag.StringVar(&isCollage, "c", "n", "(Y/n) Enables/disables collage mode\n")
	flag.StringVar(&output, "o", "img.jpg", "Select output path\n")
	flag.StringVar(&isVideo, "v", "n", "Choose photo/video mode\n")
	flag.Parse()
	fmt.Println(helloMessage)
	if len(filepath) == 0 {
		fmt.Println("Filepath can't be empty")
		return
	}
	sendData := map[string]string{
		"wm":      "",
		"collage": "0",
		"id":      "",
		"no_crop": "",
		"code":    "",
		"type":    "1",
	}
	if isCollage == "Y" {
		sendData["collage"] = "1"
	}
	if strings.HasSuffix(filepath, "mp4") || isVideo == "Y" {
		if output == "img.jpg" {
			output = "vid.mp4"
		}
		sendData["type"] = "2"
	}
	request, err := newfileUploadRequest("http://portraitplus.facefun.ai:8080/Port/MakePort", sendData, "image", filepath)
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if len(body) != 0 {
		f, err := os.Create(output)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		var bodyContent []byte
		response, err := http.Get("http://portraitplus.facefun.ai:8080/Port/" + string(body))
		if err != nil {
			log.Fatal(err)
			return
		}
		bodyContent, err = ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
			return
		}
		_, err = f.Write(bodyContent)
		if err != nil {
			log.Fatal(err)
			return
		}
	} else {
		fmt.Println("Error occurred, try again")
		log.Fatal(err)
	}
}
