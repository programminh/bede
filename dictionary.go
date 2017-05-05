package bede

import (
	"encoding/xml"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type document struct {
	XMLName      string `xml:"DOCUMENT"`
	DocumentName string `xml:"DOCUMENT_NAME"`
	Version      string `xml:"VERSION"`
	Data         string `xml:"DATA,omitempty"`
	Length       int    `xml:"LENGTH,omitempty"`
	PathToFile   string `xml:"PATH_TO_FILE,omitempty"`
}

type dictionary struct {
	XMLName   string `xml:"DICTIONARY"`
	mu        sync.Mutex
	Documents []document
}

func (dict *dictionary) Add(d document) {
	dict.mu.Lock()
	defer dict.mu.Unlock()

	dict.Documents = append(dict.Documents, d)
}

// GenDict creates client and server dictionaries by parsing the src folder
func GenDict(src, clientDst, serverDst string) (err error) {
	var (
		clientFile, serverFile *os.File
	)

	if clientFile, err = os.Open(filepath.Join(clientDst, "client.xml")); err != nil {
		return
	}

	if serverFile, err = os.Open(filepath.Join(serverDst, "server.xml")); err != nil {
		return
	}

	wg := sync.WaitGroup{}
	dict := dictionary{}

	// Walk the directory, read all the files asynchronously and add it to the dictionary
	walker := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		wg.Add(1)
		go addDocument(path, &dict, &wg)

		return nil
	}

	filepath.Walk(src, walker)
	wg.Wait()

	// Write the server file first since we need all the data
	if err = xml.NewEncoder(serverFile).Encode(dict); err != nil {
		return
	}

	return nil
}

// writeClient outputs the dictionary to client.xml
// client.xml does not contain the real path to the file therefore we must strip it
func writeClient(dict *dictionary, f *os.File) (err error) {
	for _, d := range dict.Documents {
		d.PathToFile = ""
	}

	return xml.NewEncoder(f).Encode(dict)
}

func addDocument(path string, dict *dictionary, wg *sync.WaitGroup) {
	var (
		err error
		f   *os.File
		d   document
	)

	defer wg.Done()

	// Just skip the file if we can't open it
	if f, err = os.Open(path); err != nil {
		log.Println(err)
		return
	}

	// Just skip if the file is malformed
	if err = xml.NewDecoder(f).Decode(&d); err != nil {
		log.Println(err)
		return
	}

	d.Length = len(d.Data)
	d.Data = ""
	d.PathToFile = f.Name()

	dict.Add(d)
}
