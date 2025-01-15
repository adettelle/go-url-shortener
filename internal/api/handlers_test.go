package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/adettelle/go-url-shortener/internal/config"
	"github.com/adettelle/go-url-shortener/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestPostShortPath(t *testing.T) {
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

	mockStorage.EXPECT().AddPath(strBody).Return(reqURL+id, nil)

	request, err := http.NewRequest(http.MethodPost, reqURL, strings.NewReader(strBody))
	require.NoError(t, err)

	response := httptest.NewRecorder()

	handlers.PostShortPath(response, request)

	wantHTTPStatus := http.StatusCreated
	require.Equal(t, wantHTTPStatus, response.Code)
}

func TestGetID(t *testing.T) {
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

	mockStorage.EXPECT().GetPath(id).Return(header, nil)

	request, err := http.NewRequest(http.MethodGet, reqURL, nil)
	require.NoError(t, err)
	request.SetPathValue("id", id)
	response := httptest.NewRecorder()

	handlers.GetID(response, request)

	wantHTTPStatus := http.StatusTemporaryRedirect

	require.Equal(t, wantHTTPStatus, response.Code)
	require.Equal(t, response.Header().Get("Location"), header)
}
