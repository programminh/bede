package bede

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
	Documents []document
}
