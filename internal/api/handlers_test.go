package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/adettelle/go-url-shortener/internal/config"
	"github.com/adettelle/go-url-shortener/internal/mocks"
	"github.com/carlmjohnson/requests"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

// TODO написать тесты на все хэндлеры
// TODO тесты на взаимодействие с БД и на ошибки !!!

func TestCreateShortAddressPlainText(t *testing.T) {
	// создаём контроллер
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	mockStorage := mocks.NewMockStorager(ctrl)
	cfg, err := config.New()
	require.NoError(t, err)

	handlers := &Handlers{
		repo:   mockStorage,
		config: cfg,
	}

	strBody := "https://practicum.yandex.ru/"
	reqURL := "http://" + cfg.Address + "/"
	id := "qqVjJVf"

	mockStorage.EXPECT().GetShortURLByOriginalURL(strBody).Return("", nil)
	mockStorage.EXPECT().AddOriginalURL(strBody).Return(reqURL+id, nil)

	request, err := http.NewRequest(http.MethodPost, reqURL, strings.NewReader(strBody))
	require.NoError(t, err)

	response := httptest.NewRecorder()

	handlers.CreateShortAddressPlainText(response, request)

	wantHTTPStatus := http.StatusCreated
	require.Equal(t, wantHTTPStatus, response.Code)
}

func TestGetFullAddress(t *testing.T) {
	// создаём контроллер
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	mockStorage := mocks.NewMockStorager(ctrl)

	handlers := &Handlers{
		repo: mockStorage,
	}

	id := "qqVjJVf"
	reqURL := "http://localhost:8080/"
	header := "https://practicum.yandex.ru/"

	mockStorage.EXPECT().GetOriginalURLByShortURL(id).Return(header, nil)

	request, err := http.NewRequest(http.MethodGet, reqURL, nil)
	require.NoError(t, err)
	request.SetPathValue("id", id)
	response := httptest.NewRecorder()

	handlers.GetFullAddress(response, request)

	wantHTTPStatus := http.StatusTemporaryRedirect

	require.Equal(t, wantHTTPStatus, response.Code)
	require.Equal(t, response.Header().Get("Location"), header)
}

func TestCreateShortAddressJson(t *testing.T) {
	// создаём контроллер
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// создаём объект-заглушку
	mockStorage := mocks.NewMockStorager(ctrl)
	// cfg, err := config.New()
	// require.NoError(t, err)
	cfg := &config.Config{Address: "localhost:8080", URLAddress: "http://localhost:8080"}
	handlers := &Handlers{
		repo:   mockStorage,
		config: cfg,
	}

	reqBody := shortAddrCreateRequestDTO{OriginalURL: "https://practicum.yandex.ru/"}
	reqURL := "http://" + cfg.Address + "/api/shorten"
	id := "qqVjJVf"

	mockStorage.EXPECT().GetShortURLByOriginalURL(reqBody.OriginalURL).Return("", nil)
	mockStorage.EXPECT().AddOriginalURL(reqBody.OriginalURL).Return(reqURL+id, nil)

	request, err := requests.
		URL(reqURL).
		Method(http.MethodPost).
		Header("Content-Type", "application/json").
		BodyJSON(&reqBody).
		Request(context.Background())
	require.NoError(t, err)

	wantHTTPStatus := http.StatusCreated
	response := httptest.NewRecorder()
	handlers.CreateShortAddressJSON(response, request)
	require.Equal(t, wantHTTPStatus, response.Code)
}
