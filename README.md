adv_db_server
====

## Overview

sharing advertiser datas

## Build

### require

- go 1.5 or later

```
$ go get -u github.com/keita0q/adv_db_server
$ cd /path/to/adv_db_server
$ go build
```

## SetUp

以下のような設定ファイルを作成してください。

```
{
  "context_path":"dsp",
  "port":8080,
  "save_path" :"PATH/TO/DATABASE"
}
```

## Usage

```
./adv_db_server < -c [ path to config.json ] >
```

## REST API

### POST
