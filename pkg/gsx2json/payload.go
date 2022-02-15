package gsx2json

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type MetaData struct {
	CheckSum string `json:"md5,omitempty"`
	Size     int    `json:"bytes,omitempty"`
}

type MetaDataInfo struct {
	Columns    MetaData `json:"columns,omitempty"`
	Rows       MetaData `json:"rows,omitempty"`
	Dictionary MetaData `json:"dict,omitempty"`
}

type DataView struct {
	Columns    map[string][]interface{} `json:"columns,omitempty"`
	Rows       []interface{}            `json:"rows,omitempty"`
	Dictionary map[string]interface{}   `json:"dict,omitempty"`
	Metadata   MetaDataInfo             `json:"meta,omitempty"`
}

type Payload struct {
	Values [][]string `json:"values"`
	View   DataView
}

func NewPayload() *Payload {
	return &Payload{}
}

func (p *Payload) Parse(b []byte, _cfg *Config) error {
	if err := json.Unmarshal(b, p); err != nil {
		return err
	}
	if len(p.Values) == 0 {
		return fmt.Errorf("cannot parse blank spread sheet.")
	}
	rows := make([]interface{}, 0)
	dict := make(map[string]interface{})
	columns := make(map[string][]interface{})
	headings := p.Values[0]
	for i := 1; i < len(p.Values); i++ {
		var pkey = int64(0)
		var row = p.Values[i]
		var queried = len(_cfg.Query) == 0
		newRow := make(map[string]interface{})
		for j, key := range headings {
			key = strings.ToLower(key)
			if strings.HasPrefix(key, "noex") {
				continue
			}
			var value interface{}
			if j < len(row) {
				value = row[j]
			} else {
				value = ""
			}
			if pkey == 0 {
				ivalue, err := strconv.ParseInt(value.(string), 0, 64)
				if ivalue > 0 && err == nil {
					pkey = ivalue
				}
			}
			if len(_cfg.Query) > 0 {
				lkey := strings.ToLower(key)
				lvalue := strings.ToLower(value.(string))
				lquery := strings.ToLower(_cfg.Query)
				queried = queried || strings.Contains(lkey, lquery)
				queried = queried || strings.Contains(lvalue, lquery)
			}
			cast := func() interface{} {
				return value.(string)
			}
			if _cfg.UseInteger {
				if len(value.(string)) == 0 {
					cast = func() interface{} {
						return int64(0)
					}
				} else {
					ivalue, err := strconv.ParseInt(value.(string), 0, 64)
					if err == nil {
						cast = func() interface{} {
							return ivalue
						}
					}
				}
			}
			newRow[key] = cast()
			if queried {
				columns[key] = append(columns[key], cast())
			}
		}
		if queried {
			rows = append(rows, row)
			dict[fmt.Sprintf("%d", pkey)] = row
		}
	}
	if _cfg.ShowColumns && !_cfg.BriefMeta {
		b, err := json.Marshal(columns)
		if err != nil {
			return err
		}
		hash := md5.Sum(b)
		p.View.Columns = columns
		metadata := p.View.Metadata.Columns
		metadata.Size = len(b)
		metadata.CheckSum = hex.EncodeToString(hash[:])
	}
	if _cfg.ShowRows && !_cfg.BriefMeta {
		b, err := json.Marshal(rows)
		if err != nil {
			return err
		}
		hash := md5.Sum(b)
		p.View.Rows = rows
		metadata := p.View.Metadata.Rows
		metadata.Size = len(b)
		metadata.CheckSum = hex.EncodeToString(hash[:])
	}
	if _cfg.ShowDict && !_cfg.BriefMeta {
		b, err := json.Marshal(dict)
		if err != nil {
			return err
		}
		hash := md5.Sum(b)
		p.View.Dictionary = dict
		metadata := p.View.Metadata.Dictionary
		metadata.Size = len(b)
		metadata.CheckSum = hex.EncodeToString(hash[:])
	}
	return nil
}
