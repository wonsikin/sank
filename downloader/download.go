package downloader

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	mimeType = map[string]string{"image/jpeg": "jpg"}
)

// DownloadFile download file to specify foler
func DownloadFile(folder, fileName, url string) (file string, err error) {
	// return fileName + ".jpg", nil
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ct := http.DetectContentType(data)
	output := folder + fileName
	if _, ok := mimeType[ct]; ok {
		output += "." + mimeType[ct]
	}

	out, err := os.Create(output)
	if err != nil {
		return "", err
	}
	defer out.Close()

	content := bytes.NewReader(data)

	_, err = io.Copy(out, content)
	if err != nil {
		return "", err
	}

	return fileName + "." + mimeType[ct], nil
}
