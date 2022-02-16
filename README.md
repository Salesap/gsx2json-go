[![GitHub license](https://img.shields.io/github/license/deflinhec/gsx2json-go.svg)](https://github.com/deflinhec/gsx2json-go/blob/master/LICENSE) 
[![GitHub release](https://img.shields.io/github/release/deflinhec/gsx2json-go.svg)](https://github.com/deflinhec/gsx2json-go/releases/)
[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://github.com/deflinhec/gsx2json-go/graphs/commit-activity)
# GSX2JSON-GO - Google Spreadsheet to JSON API Go service.

Inspired by [55sketch/gsx2Json](https://github.com/55sketch/gsx2json), preserve all functinality of origin and implement additional features.

## :speech_balloon: About
After Google SpreadSheet API v4 is out, it's pretty easy to get a nice decent JSON format directly from. This API build on top of v4 API, downloading spread sheet and generates diffirent data views, and also make sure all rows have the exact same length according to the heading. Additional feature, such as `file cache` and `md5 checksum`.

This is a backport version of 
[gsx2jsonpp](https://github.com/deflinhec/gsx2jsonpp), which some of my project relies on those addtional features. [gsx2jsonpp](https://github.com/deflinhec/gsx2jsonpp) might not be mantain in the future, since uriparser doesn't work properly with UTF-8 URI in some circumstance.



 
## :whale: Using docker image

Make sure [docker engine](https://www.docker.com/products/docker-desktop) has already install in your operating-system.

In this example below I'm going to use `5000` as port, and output log file under `bin/volume` directory.

- Launch with a remote image

    ```
    docker pull deflinhec/gsx2json-go:latest
    ```
    
    ```
    docker run -it -d --rm deflinhec/gsx2json-go --port 5000
    ```

- Launch with a local image
    
    Follow instructions below, :toolbox: Build from source.

    ```
    docker build -t gsx2json-go .
    ```
    
    ```
    docker run -it -d --rm gsx2json-go --port 5000
    ```

    Avaliable arguments: 
    - -p, --port
    - --host 

    - --cache (file|memory)

      Cache mode is disabled by default, this feature allows client
      to query on specific data version. When cache is configure with
      file mode, cache file will be preserve under `.cache/` folder. 

    - --ssl

      SSL mode is disabled by default, if you prefer using SSL mode
      with your certification and private key. Copy your files into
      `cert/` and rename as `cert.pem, key.pem`.

After launched, gsx2json-go should be accessable in your browser [localhost:5000](http://localhost:5000/version).

## :dart: Spreadsheet rule

- Column name begin with `NOEX_` will not export to the final result.

- Make sure to add a left most column represents as an unique integer key, require for dicitionay data view.

### :memo: API

- GET `api`

    - **id (required):** 
    
        The ID of your document. This is the big long aplha-numeric code in the middle of your document URL.

    - **sheet (required):** 
    
        The name of the individual sheet you want to get data from.

    - **api_key (required):** 
    
        Generate [Google API key](https://developers.google.com/sheets/api/guides/authorizing#APIKey) through this guide.

    - **q (optional):** 
    
        A simple query string. This is case insensitive and will add any row containing the string in any cell to the filtered result.

    - **integers (optional - default: true)**: 
    
        Setting 'integers' to false will return numbers as a string.

    - **dict (optional - default: true)**: 
    
        Setting 'dict' to false will return rows and columns view.

    - **rows (optional - default: true)**: 
    
        Setting 'rows' to false will return dictionary and columns view.

    - **columns (optional - default: true)**: 
    
        Setting 'columns' to false will return dictionary and rows view.

    - **meta (optional - default: false)**: 
    
        Setting 'meta' to true will return only meta data.

    - **pretty (optional - default: false)**: 
    
        Pretty print the result if sets to true.

- GET `cache`

    Get list of cache.

- DELETE `cache`

    Remove all cache data.

- Get `version`

    Show build version.

## :bookmark: Example Response

[Example spreadsheet](https://docs.google.com/spreadsheets/d/1-DGS8kSiBrPOxvyM1ISCxtdqWt-I7u1Vmcp-XksQ1M4/edit#gid=0)

There are four sections to the returned data.

- Columns (containing the names of each column)
- Dictionary (used left most column as primary key)
- Rows (containing each row of data as an object)
- Meta (contains short brief of target data)

```
{
 "columns": {
  "key": [
   1,
   2,
   3,
   4
  ],
  "column1": [
   "1b",
   "2b",
   "3b",
   "4b"
  ],
  "column2": [
   11,
   22,
   33,
   44
  ]
 },
 "rows": [
  {
   "key": 1,
   "column1": "1b",
   "column2": 11
  },
  {
   "key": 2,
   "column1": "2b",
   "column2": 22
  },
  {
   "key": 3,
   "column1": "3b",
   "column2": 33
  },
  {
   "key": 4,
   "column1": "4b",
   "column2": 44
  }
 ],
 "dict": {
  "1": {
   "key": 1,
   "column1": "1b",
   "column2": 11
  },
  "2": {
   "key": 2,
   "column1": "2b",
   "column2": 22
  },
  "3": {
   "key": 3,
   "column1": "3b",
   "column2": 33
  },
  "4": {
   "key": 4,
   "column1": "4b",
   "column2": 44
  }
 },
 "meta": {
  "columns": {
   "md5": "EAC2F0EF3EA62CEEDD3B65B627B06CBA",
   "bytes": 73
  },
  "rows": {
   "md5": "7767981744A818A7574B4A0B8EBE1C25",
   "bytes": 153
  },
  "dict": {
   "md5": "76C73EAEAFC8BA2ACD890C50E20C1613",
   "bytes": 169
  }
 }
}

```

## :coffee: Donation

If you like my project and also appericate for the effort. Don't hesitate to [buy me a coffee](https://ko-fi.com/deflinhec)😊.
