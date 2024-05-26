# Membuat Registrasi User

Kali ini kita akan lanjutkan membuat modul Login dan Registrasi User.

## Membuat Tabel

Buat tabel `users` di dalam database yang sudah kita buat pada tuturial sebelumnya.
adapun field-field nya adalah sebagai berikut:

1. id         int     auto_increment
2. name       varchar 100
3. email      varchar 100
4. password   varchar 255
5. gender     enum 'male';'female'
6. created_at timestamp
7. updated_at timestamp

## Membuat Entity

Buat folder baru dengan nama `entity`, selanjutnya buat file `user.go` didalam folder
tersebut.

```go
  package entity

import "time"

type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	Gender    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
```

## Membuat file Auth

Buat file `auth_dto.go` didalam folder `dto`.

```go
  package dto

type RegisterRequest struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirm"`
	Gender               string `json:"gender"`
}
```

## Repository Pattern

Karena di Rest API nya ini kita akan menerapkan repository pattern, maka kita
akan membuat 4 buah folder yaitu:

- **reporitory** folder ini berisikan fungsi untuk komunikasi ke database
- **handler** untuk menerima request.
- **service** yang berisikan logik-logik disini.
- **router** untuk list-list api nya.

### Membuat file Auth Repository

Buat file `auth_repository.go` di dalam folder `repository`.

```go
 package repository

import (
	"GO-SOCMED/entity"

	"gorm.io/gorm"
)

type AuthRepository interface {
	EmailExist(email string) bool
	Register(req *entity.User) error
}

/*
	auth menggunakan huruf kecil berarti

hanya bisa di panggil di dalam file
auth_repository
*/
type authRepository struct {
	db *gorm.DB
}

/* fungsi constractor yang diawali dgn New */
func NewAuthRepository(db *gorm.DB) *authRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) Register(user *entity.User) error {
	err := r.db.Create(&user).Error

	return err
}

func (r *authRepository) EmailExist(email string) bool {
	var user entity.User
	err := r.db.First(&user, "email = ?", email).Error

	return err == nil
} 
```

### Membuat file Auth Service

Buat file `auth_service.go` di dalam folder `service`.

```go
package service

import (
	"GO-SOCMED/dto"
	"GO-SOCMED/entity"
	"GO-SOCMED/errorhandler"
	"GO-SOCMED/helper"
	"GO-SOCMED/repository"
)

type AuthService interface {
	Register(req *dto.RegisterRequest) error
}

type authService struct {
	repository repository.AuthRepository
}

func NewAuthService(r repository.AuthRepository) *authService {
	return &authService{
		repository: r,
	}
}

func (s *authService) Register(req *dto.RegisterRequest) error {
	if emailExist := s.repository.EmailExist(req.Email); emailExist {
		return &errorhandler.BadRequestError{Message: "email already registered"}
	}

	if req.Password != req.PasswordConfirmation {
		return &errorhandler.BadRequestError{Message: "password not match"}
	}

	passwordHash, err := helper.HashPassword(req.Password)
	if err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}

	}

	user := entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: passwordHash,
		Gender:   req.Gender,
	}

	if err := s.repository.Register(&user); err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	return nil
} 
```

### Membuat file Auth Handler

Buat file `auth_handler.go` di dalam folder `handler`.

```go
 package handler

import (
	"GO-SOCMED/dto"
	"GO-SOCMED/errorhandler"
	"GO-SOCMED/helper"
	"GO-SOCMED/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *authHandler {
	return &authHandler{
		service: s,
	}
}

// method
func (h *authHandler) Register(c *gin.Context) {
	var register dto.RegisterRequest

	// model binding and validation
	if err := c.ShouldBindJSON(&register); err != nil {
		errorhandler.HandleError(c, &errorhandler.BadRequestError{Message: err.Error()})
		// agar code dibawah tidak tereksikusi tambahkan return
		return
	}

	if err := h.service.Register(&register); err != nil {
		errorhandler.HandleError(c, err)
		return
	}

	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusCreated,
		Message:    "Register successfully, please login",
	})

	c.JSON(http.StatusCreated, res)
} 
```

### Membuat Helper Baru

Disini kita akan membuat helper baru, untuk hashing password.
Buat file `password.go` didalam folder `helper`.

```go
  package helper

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(passwordHash), err
}
```

### Membuat Router

But file `auth_router.go` di dalam folder `router`.

```go
  package router

import (
	"GO-SOCMED/config"
	"GO-SOCMED/handler"
	"GO-SOCMED/repository"
	"GO-SOCMED/service"

	"github.com/gin-gonic/gin"
)

func AuthRouter(api *gin.RouterGroup) {
	authRepository := repository.NewAuthRepository(config.DB)
	authService := service.NewAuthService(authRepository)
	authHandler := handler.NewAuthHandler(authService)

	api.POST("/register", authHandler.Register)
}
```

## Register File Router

Daftarkan file router ke dalam file `main.go`

```go
  package main

import (
	"GO-SOCMED/config"
	"GO-SOCMED/router"
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

	// router
	router.AuthRouter(api)

	app.Run(fmt.Sprintf(":%v", config.ENV.PORT))
}
```


## Coba di Postmen
