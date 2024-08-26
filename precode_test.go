package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Тест 1: Корректный запрос, статус 200, тело ответа не пустое
func TestMainHandlerValidRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=moscow", nil)
	responseRecorder := httptest.NewRecorder()

	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Проверяем с использованием require, так как код 200 критичен для дальнейших проверок
	require.Equal(t, http.StatusOK, responseRecorder.Code, "Ожидается код 200")

	// Проверка, что тело ответа не пустое, важна для понимания, что сервис работает правильно
	require.NotEmpty(t, responseRecorder.Body.String(), "Ответ должен быть непустым")

	// Делаем assert для дальнейшей проверки содержания
	expectedResponse := "Мир кофе,Сладкоежка"
	assert.Equal(t, expectedResponse, responseRecorder.Body.String(), "Ожидается два кафе")
}

// Тест 2: Некорректный город, статус 400 и ошибка "wrong city value"
func TestMainHandlerWrongCity(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=paris", nil)
	responseRecorder := httptest.NewRecorder()

	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Код 400 — это критическая проверка, если она не пройдёт, тест смысла не имеет
	require.Equal(t, http.StatusBadRequest, responseRecorder.Code, "Ожидается код 400")

	// Делаем assert для проверки текста ошибки
	assert.Equal(t, "wrong city value", responseRecorder.Body.String(), "Ожидается ошибка 'wrong city value'")
}

// Тест 3: Если count больше, чем есть всего кафе, возвращаются все кафе
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=10&city=moscow", nil) // count больше, чем доступно
	responseRecorder := httptest.NewRecorder()

	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Проверка на критичный код ответа
	require.Equal(t, http.StatusOK, responseRecorder.Code, "Ожидается код 200")

	// Проверка, что возвращены все кафе
	expectedResponse := "Мир кофе,Сладкоежка,Кофе и завтраки,Сытый студент"
	assert.Equal(t, expectedResponse, responseRecorder.Body.String(), "Ожидаются все доступные кафе")
}
