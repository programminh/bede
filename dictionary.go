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
	Documents []document `xml:"DOCUMENT"`
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

	if clientFile, err = os.Create(clientDst); err != nil {
		return
	}

	if serverFile, err = os.Create(serverDst); err != nil {
		return
	}

	defer func() {
		serverFile.Close()
		clientFile.Close()
	}()

	wg := sync.WaitGroup{}
	dict := dictionary{}

	// Walk the directory, read all the files asynchronously and add it to the dictionary
	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		wg.Add(1)
		go addDocument(path, &dict, &wg)

		return nil
	}

	if err = filepath.Walk(src, walker); err != nil {
		return
	}

	wg.Wait()

	// Write the server file first since we need all the data
	if err = xml.NewEncoder(serverFile).Encode(dict); err != nil {
		return
	}

	return writeClient(&dict, clientFile)
}

// writeClient outputs the dictionary to client.xml
// client.xml does not contain the real path to the file therefore we must strip it
func writeClient(dict *dictionary, f *os.File) (err error) {
	for i, _ := range dict.Documents {
		dict.Documents[i].PathToFile = ""
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
		log.Println("can't open:", err)
		return
	}

	// Just skip if the file is malformed
	if err = xml.NewDecoder(f).Decode(&d); err != nil {
		log.Println("can't decode", f.Name(), ":", err)
		return
	}

	// Let's assume that data is just a string and we want to know the length of the string
	d.Length = len(d.Data)
	// Set data to empty string so we don't write it out to the dictionnary
	d.Data = ""
	d.PathToFile = f.Name()

	dict.Add(d)
}
