# Struktur Folder dan Konfigurasi Project

Pada project kali ini kita akan belajar golang, dengan membuat Rest API
Media Sosial (Twitter Clone). Kita akan menggunakan Framework GIN dan kita
juga akan menerapkan Repository Pattern.

## Setup Project

Buat sebuah folder baru dengan nama `GO-SOCMED`, folder ini akan menjadi
project kita nantinya. Selanjutnya kita akan melakukan inisialisasi project
go dengan mengetikkan `go mod init go-socmed`.

## Struktur Folder Project

Selanjutnya buat struktur projectnya yang terdiri dari folder:

- **config** untuk menyimpan konfigurasi database dan env.
- **dto** (Data Transfer Objek)
- **errorhandler** untuk memuat dan menyimpan custome error.
- **helper** untuk fungsi-fungsi yang dipakai berkali-kali.
- **main.go** untuk file utamanya.

## Instalasi GIN Framework

```bash
  go get -u github.com/gin-gonic/gin
```

buka file `main.go` dan panggil gin nya.

```go
  app := gin.Default()

  app.Get("/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "pong"
    })
  })

  app.run(":5000")
```

## Installasi Depedensi

- Viper => untuk konfigurasi managemen

```bash
  go get github.com/spf13/viper
```

- GORM => untuk database ORM di Go

```bash
  go get -u gorm.io/gorm
```

```bash
  go get -u gorm.io/driver/mysql
```

## Konfigurasi Koneksi Database dengan GORM

Selanjutnya buat file `.env`

```yaml
  PORT=5000
  DB_USERNAME="root"
  DB_PASSWORD="supersecretpassword"
  DB_URL="127.0.0.1:3306"
  DB_DATABASE="go_socmed"
```

Kemudian buat sebuah file `config.go` di dalam folder `config`.

```go
  package config

import "github.com/spf13/viper"

type Config struct {
	PORT        string
	DB_USERNAME string
	DB_PASSWORD string
	DB_URL      string
	DB_DATABASE string
}

var ENV *Config

func LoadConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&ENV); err != nil {
		panic(err)
	}
}
```

Selanjutnya buat file `database.go` di dalam folder `config`:

```go
  package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadDB() {
	connectionStr := fmt.Sprintf("%v:%v@tcp(%v)/%v?%v", ENV.DB_USERNAME, ENV.DB_PASSWORD, ENV.DB_URL, ENV.DB_DATABASE, "charset=utf8mb4&parseTime=true&loc=Asia%2FJakarta")
	db, err := gorm.Open(mysql.Open(connectionStr), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
}
```

## Data Transfer Object (DTO)

### Pagination

Buat file `paginate_dto.go` didalam folder `dto`:

```go
  package dto

type Paginate struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}
```

### Params
Selanjutnya kita akan membuat rapi response datanya, buat file `params_dto.go`
didalam folder `dto`:

```go
  package dto

type ResponseParams struct {
	StatusCode int
	Message    string
	Paginate   *Paginate
	Data       any
}
```

## Membuat Helper

### Response
Sekarang kita akan membuat helper untuk custome reponsenya, buat file `response.go`
didalam folder `helper`:

```go
  package helper

import "GO-SOCMED/dto"

type ResponseWithData struct {
	Code     int           `json:"code"`
	Status   string        `json:"status"`
	Message  string        `json:"message"`
	Paginate *dto.Paginate `json:"paginate,omitempty"`
	Data     any           `json:"data"`
}

type ResponseWithoutData struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Response(params dto.ResponseParams) any {
	var response any
	var status string

	if params.StatusCode >= 200 && params.StatusCode <= 299 {
		status = "success"
	} else {
		status = "failed"
	}

	if params.Data != nil {
		response = &ResponseWithData{
			Code:     params.StatusCode,
			Status:   status,
			Message:  params.Message,
			Paginate: params.Paginate,
			Data:     params.Data,
		}
	} else {
		response = &ResponseWithoutData{
			Code:    params.StatusCode,
			Status:  status,
			Message: params.Message,
		}
	}

	return response
}
```

## Membuat Handler

### Membuat Tipe Data Baru untuk Handler

Buat file `types.go` di dalam folder `errorhandler`:

```go
  package errorhandler

type NotFoundError struct {
	Message string
}

type BadRequestError struct {
	Message string
}

type InternalServerError struct {
	Message string
}

type UnauthorizedError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func (e *BadRequestError) Error() string {
	return e.Message
}

func (e *InternalServerError) Error() string {
	return e.Message
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}
```

Lalu buat file `error_handler.go` di dalam folder `errorhandler`:

```go
  package errorhandler

import (
	"GO-SOCMED/dto"
	"GO-SOCMED/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	var statuscode int

	switch err.(type) {
	case *NotFoundError:
		statuscode = http.StatusNotFound
	case *BadRequestError:
		statuscode = http.StatusBadRequest
	case *InternalServerError:
		statuscode = http.StatusInternalServerError
	case *UnauthorizedError:
		statuscode = http.StatusUnauthorized
	}

	response := helper.Response(dto.ResponseParams{
		StatusCode: statuscode,
		Message:    err.Error(),
	})

	c.JSON(statuscode, response)
}
```

## Memanggil Konfigurasi di dalam file Main.go

Buka file `main.go` dan ubah code seperti berikut:

```go
  package main

import (
	"GO-SOCMED/config"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	config.LoadDB()

	app := gin.Default()
	api := app.Group("/api")

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	app.Run(fmt.Sprintf(":%v", config.ENV.PORT))
}
```

