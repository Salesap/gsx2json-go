package gsx2json

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"gitlab.com/c0b/go-ordered-json"
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
	Columns    *ordered.OrderedMap `json:"columns,omitempty"`
	Rows       []interface{}       `json:"rows,omitempty"`
	Dictionary *ordered.OrderedMap `json:"dict,omitempty"`
	Metadata   MetaDataInfo        `json:"meta,omitempty"`
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
	dictionary := ordered.NewOrderedMap()
	columns := ordered.NewOrderedMap()
	headings := p.Values[0]
	for i := 1; i < len(p.Values); i++ {
		var pkey = int64(0)
		var row = p.Values[i]
		var queried = len(_cfg.Query) == 0
		newRow := ordered.NewOrderedMap()
		for j, key := range headings {
			if strings.HasPrefix(key, "NOEX_") {
				continue
			}
			var value interface{}
			if j >= len(row) {
				row = append(row, "")
			}
			value = row[j]
			if pkey == 0 {
				ivalue, err := strconv.ParseInt(value.(string), 0, 64)
				if ivalue > 0 && err == nil {
					pkey = ivalue
				}
			}
			if len(_cfg.Query) > 0 {
				queried = queried || strings.Contains(key, _cfg.Query)
				queried = queried || strings.Contains(value.(string), _cfg.Query)
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
			newRow.Set(key, cast())
			if queried {
				column := make([]interface{}, 0)
				if last, ok := columns.GetValue(key); ok {
					column = last.([]interface{})
				}
				column = append(column, cast())
				columns.Set(key, column)
			}
		}
		if queried {
			rows = append(rows, newRow)
			dictionary.Set(fmt.Sprintf("%d", pkey), newRow)
		}
	}
	if _cfg.ShowColumns {
		if !_cfg.BriefMeta {
			p.View.Columns = columns
		}
		if b, err := json.Marshal(columns); err == nil {
			bytes := md5.Sum(b)
			hash := hex.EncodeToString(bytes[:])
			p.View.Metadata.Columns.Size = len(b)
			p.View.Metadata.Columns.CheckSum = hash
		}
	}
	if _cfg.ShowRows {
		if !_cfg.BriefMeta {
			p.View.Rows = rows
		}
		if b, err := json.Marshal(rows); err == nil {
			bytes := md5.Sum(b)
			hash := hex.EncodeToString(bytes[:])
			p.View.Metadata.Rows.Size = len(b)
			p.View.Metadata.Rows.CheckSum = hash
		}
	}
	if _cfg.ShowDict {
		if !_cfg.BriefMeta {
			p.View.Dictionary = dictionary
		}
		if b, err := json.Marshal(dictionary); err == nil {
			bytes := md5.Sum(b)
			hash := hex.EncodeToString(bytes[:])
			p.View.Metadata.Dictionary.Size = len(b)
			p.View.Metadata.Dictionary.CheckSum = hash
		}
	}
	return nil
}
