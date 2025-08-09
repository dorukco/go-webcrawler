package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.LoadHTMLGlob("../templates/*")

	router.GET("/", IndexHandler)
	router.POST("/submit", SubmitHandler)

	return router
}

func TestIndexHandler(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestSubmitHandler_EmptyInput(t *testing.T) {
	router := setupTestRouter()

	form := url.Values{}
	form.Add("text_input", "")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/submit", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if !strings.Contains(w.Body.String(), "Please enter a URL!") {
		t.Error("Expected error message not found in response")
	}
}

func TestSubmitHandler_InvalidURL(t *testing.T) {
	router := setupTestRouter()

	form := url.Values{}
	form.Add("text_input", "invalid-url")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/submit", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if !strings.Contains(w.Body.String(), "Please enter a valid URL") {
		t.Error("Expected invalid URL error message not found in response")
	}
}

func TestSubmitHandler_ValidURL(t *testing.T) {
	router := setupTestRouter()

	form := url.Values{}
	form.Add("text_input", "https://doruk.com")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/submit", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}
