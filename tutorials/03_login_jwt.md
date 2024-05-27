# Login Dengan Menggunakan JWT

Karena pada tutorial sebelumnya kita sudah berhasil membuat API
registrasi user. Sekarang kita lanjutkan ke pembuatan API Login.

## Membuat Struct Login

Buka kembali file `auth_dto.go` didalam folder `dto`, dan buat struct 
baru dengan nama `LoginRequest` dan `LoginResponse`.

```go
  type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}
```

## Membuat Handler Login

Buka file `auth_handler.go` didalam folder `handler` dan tambahkan sebuah
method dengan nama `Login`.

```go
  func (h *authHandler) Login(c *gin.Context) {
	var login dto.LoginRequest

	err := c.ShouldBindJSON(&login)
	if err != nil {
		errorhandler.HandleError(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	result, err := h.service.Login(&login)
	if err != nil {
		errorhandler.HandleError(c, err)
		return
	}

	res := helper.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "Successfully login",
		Data:       result,
	})

	c.JSON(http.StatusOK, res)
}
```

## Membuat Login Repository

Karena di service login ada pengecekan untuk email, maka kita akan
menambahkan sebuah repository untuk login tersebut. buka file 
`auth_repository.go`.

Pada interface tambahkan seperti berikut:

```go
  type AuthRepository interface {
    ....
    GetUserByEmail(email string) (*entity.User, error)
  }
```

Selanjutnya buat sebuah method seperti berikut:

```go
  .....

  func (r *authRepository) GetUserByEmail(email string) (*entity.User, error) {
    var user entity.User
    err := r.db.First(&user, "email = ?", email).Error

    return &user, err
  }
```

## Membuat Helper Generate Token

Buka file `password.go` didalam folder `helper`, dan buat fungsi dengan
nama `VerifyPassword`.

```go
  func VerifyPassword(hashPassword string, password string) error {
    err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
    return err
  }
```

Untuk generate tokennya buat file baru dengan nama `token.go`. Sebelumnya
kita akan menginstall module JWT `go get github.com/golang-jwt/jwt/v4`.
kita lanjutkan ke generate token.

```go
  package helper

import (
	"GO-SOCMED/entity"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// secret key
var mySigningKey = []byte("mysecretkey")

type JWTClaims struct {
	ID int `json:"id"`
	jwt.RegisteredClaims
}

func GenerateToken(user *entity.User) (string, error) {
	claims := JWTClaims{
		user.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 + time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)

	return ss, err
}
```

## Membuat Service Login

Buka file `auth_service.go` didalam folder `service`. tambahkan AuthService Interfacenya.

```go
  type AuthService interface {
    ....
    Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
  }
```

Selanjutnya buat sebuah method di dalam file `auth_service.go`.

```go
  ........

  func (s *authService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	var data dto.LoginResponse

	user, err := s.repository.GetUserByEmail(req.Email)
	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "wrong email or password"}
	}

	if err := helper.VerifyPassword(user.Password, req.Password); err != nil {
		return nil, &errorhandler.NotFoundError{Message: "wrong email or password"}
	}

	token, err := helper.GenerateToken(user)
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	data = dto.LoginResponse{
		ID:    user.ID,
		Name:  user.Name,
		Token: token,
	}

	return &data, nil
}   
```

## Menambahkan Router Login

Buka file `auth_router.go` didalam folder `router` dan tambahkan kode berikut:

```go
  func AuthRouter(api *gin.RouterGroup) {
	  ......
    
    api.POST("/login", authHandler.Login)
}
```

Selanjutnya coba jalankan aplikasinya.
