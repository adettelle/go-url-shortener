package api

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/adettelle/go-url-shortener/internal/config"
	"github.com/adettelle/go-url-shortener/internal/storage/custrepo"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type ICustomerRepo interface {
	AddCustomer(ctx context.Context, login, masterPassword string) (int, error)
	// GetCustomerByLogin(ctx context.Context, login string) (*repo.CustomerGetByLogin, error)
	VerifyUser(ctx context.Context, login string, pass string) (bool, int)
}

type CustomerHandlers struct {
	CustomerRepo ICustomerRepo
	// JwtRepo      IJwtRepo
	SignKey []byte
	Config  *config.Config
}

func NewCustomerHandlers(
	customerRepo ICustomerRepo,
	// jwtRepo IJwtRepo,
	signKey []byte, cfg *config.Config) *CustomerHandlers {
	return &CustomerHandlers{
		CustomerRepo: customerRepo,
		// JwtRepo:      jwtRepo,
		SignKey: signKey,
		Config:  cfg,
	}
}

type CustomerRegistrationRequestDTO struct {
	Login    string `json:"login" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3"`
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func (ch *CustomerHandlers) RegisterCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var buf bytes.Buffer
	var customer CustomerRegistrationRequestDTO

	// читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		errlog.Error("error in reading body", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &customer); err != nil {
		errlog.Error("error in unmarshalling json:", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = validate.Struct(customer)
	if err != nil {
		errlog.Error("error in validating:", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// получаем хеш пароля
	hashedPassword := sha256.Sum256([]byte(customer.Password))
	hashStringPassword := hex.EncodeToString(hashedPassword[:]) // дополнительно кодируем пароль

	_, err = ch.CustomerRepo.AddCustomer(context.Background(), customer.Login, hashStringPassword)
	if err != nil {
		if custrepo.IsCustomerExistsErr(err) {
			errlog.Error("error in registering user:", zap.Error(err))
			w.WriteHeader(http.StatusConflict)
			return
		}
		errlog.Error("error in adding user:", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("customer: ", customer)

	w.WriteHeader(http.StatusOK)
}

type authRequestDTO struct {
	Login    string `json:"login" validate:"required,email"`
	Password string `json:"pwd" validate:"required,min=3"`
}

// `json:"password" validate:"required,min=3"`
// Login происходит по логину и паролю (password)
func (ch *CustomerHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var auth authRequestDTO

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		errlog.Error("error in reading body:", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &auth); err != nil {
		errlog.Error("error in unmarshalling json:", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = validate.Struct(auth)
	if err != nil {
		errlog.Error("error in validating:", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("111111111111111", ch.CustomerRepo)
	log.Println("2222222222222", auth.Login, auth.Password)

	ok, custID := ch.CustomerRepo.VerifyUser(context.Background(), auth.Login, auth.Password)
	log.Println("!!!!!!!!!!!!!!", custID)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized) // неверная пара логин/пароль
		return
	}

	// token, err := jwt.GenerateJwtToken(ch.SignKey, auth.Login, custID)
	// if err != nil {
	// 	errlog.Error("error in generating token:", zap.Error(err))
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// w.Header().Set("Authorization", "Bearer "+token)

	custCookie := http.Cookie{Name: "custCookie", Value: strconv.Itoa(custID), Path: "/"}
	// log.Println("::::::::::::::::::", custCookie.Name, custCookie.Value)
	http.SetCookie(w, &custCookie)
}
