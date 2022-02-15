package gsx2json

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	SPREADSHEET_HOST       = "sheets.googleapis.com"
	SPREADSHEET_URI_FORMAT = "/v4/spreadsheets/%s/values/%s?key=%s"
)

type SpreaSheet struct {
}

func Request(id *Identifier) ([]byte, error) {
	var b []byte
	uri := fmt.Sprintf(SPREADSHEET_URI_FORMAT,
		id.SheetId, id.SheetName,
		id.ApiKey)
	url := "https://" + SPREADSHEET_HOST + uri
	res, err := http.Get(url)
	if err != nil {
		return b, err
	}
	defer res.Body.Close()
	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return b, err
	}
	return b, err
}
