package parser

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/CardInfoLink/log"

	"github.com/wonsikin/sank/downloader"
)

// Parser parser the file and output de result file
func Parser(input, output string) (err error) {
	file, err := ReadFile(input)
	if err != nil {
		return err
	}
	defer file.Close()

	// header, body, err :=
	header, body, err := ReadContent(file)
	if err != nil {
		return err
	}

	// result := make([][]string, 0)

	result := append(body[:0], append([][]string{header}, body[0:]...)...)
	// the index of the remote URL field
	idxRU, fieldRU := -1, "remoteURL"
	// the index of the image name field
	idxIN, fieldIN := -1, "imgName"
	// the index of the new URL field
	idxU, fieldURL := -1, "url"

	for i, field := range header {
		if field == fieldRU {
			idxRU = i
		}
		if field == fieldIN {
			idxIN = i
		}
		if field == fieldURL {
			idxU = i
		}
	}

	if _, err1 := os.Stat(output); os.IsNotExist(err1) {
		os.Mkdir(output, os.ModePerm)
	}

	var wg sync.WaitGroup

	for i := 1; i < len(result); i++ {
		wg.Add(1)
		row := result[i]
		ru := row[idxRU]
		name := row[idxIN]

		go func(url, name string, row []string) {
			defer wg.Done()

			fileName, err1 := downloader.DownloadFile(output, name, ru)
			if err1 != nil {
				log.Errorf("error caught when download file(%s): %s", ru, err1)
				return
			}

			row[idxU] = fileName
		}(ru, name, row)
	}

	wg.Wait()

	nFile, err := os.Create(input)
	if err != nil {
		log.Errorf("error caught create file(%s): %s", input, err)
		return err
	}

	w := csv.NewWriter(nFile)

	w.WriteAll(result)
	w.Flush()

	if err := w.Error(); err != nil {
		log.Errorf("error caught write file(%s): %s", input, err)
		return err
	}
	return nil
}

// ReadFile read file by input
func ReadFile(input string) (file *os.File, err error) {
	goPaths := filepath.SplitList("GOPATH")
	if len(goPaths) == 0 {
		return nil, errors.New("GOPATH environment variable is not set or empty")
	}

	goRoot := runtime.GOROOT()
	if goRoot == "" {
		return nil, errors.New("GOROOT environment variable is not set or empty")
	}

	absPath, err := filepath.Abs(input)
	if err != nil {
		return nil, err
	}

	file, err = os.Open(absPath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// ReadContent read the content of the csv file , handler func(string)
func ReadContent(file *os.File) (header []string, body [][]string, err error) {
	r := csv.NewReader(file)
	// 逐行读取
	records, err := r.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	header = records[0]
	body = records[1:]

	return header, body, nil
}

func jsonMarshal(v interface{}, safeEncoding bool) (data []byte, err error) {
	data, err = json.Marshal(v)

	if safeEncoding {
		data = bytes.Replace(data, []byte("\\u0026"), []byte("&"), -1)
		data = bytes.Replace(data, []byte("\\u003c"), []byte("<"), -1)
		data = bytes.Replace(data, []byte("\\u003e"), []byte(">"), -1)
	}

	return data, err
}
