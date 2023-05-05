package controllers_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sibivishnu/Transactions/controllers"
	"github.com/sibivishnu/Transactions/models"
	"github.com/sibivishnu/Transactions/routes"
	"github.com/stretchr/testify/assert"
)

func TestAPI(t *testing.T) {
	router := routes.RegisterRoutes()

	t.Run("PostTransaction", func(t *testing.T) {
		transaction := models.Transaction{
			Amount:    12.3,
			Timestamp: time.Now().UTC(),
		}
		body, _ := json.Marshal(transaction)

		req, _ := http.NewRequest("POST", "/transactions", bytes.NewReader(body))
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
	})

	t.Run("PostTransactionWithFutureTimestamp", func(t *testing.T) {
		transaction := models.Transaction{
			Amount:    12.3,
			Timestamp: time.Now().UTC().Add(time.Minute),
		}
		body, _ := json.Marshal(transaction)

		req, _ := http.NewRequest("POST", "/transactions", bytes.NewReader(body))
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
	})

	t.Run("PostTransactionWithOldTimestamp", func(t *testing.T) {
		transaction := models.Transaction{
			Amount:    12.3,
			Timestamp: time.Now().UTC().Add(-61 * time.Second),
		}
		body, _ := json.Marshal(transaction)

		req, _ := http.NewRequest("POST", "/transactions", bytes.NewReader(body))
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNoContent, resp.Code)
	})

	t.Run("DeleteTransactions", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/transactions", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNoContent, resp.Code)
	})

	t.Run("GetStatistics", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/statistics", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var stats models.Statistics
		json.Unmarshal(resp.Body.Bytes(), &stats)
		assert.GreaterOrEqual(t, stats.Count, int64(0))
	})

	t.Run("SetLocation", func(t *testing.T) {
		location := models.Location{
			City: "San Francisco",
		}
		body, _ := json.Marshal(location)

		req, _ := http.NewRequest("POST", "/location", bytes.NewReader(body))
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
	})

	t.Run("ResetLocation", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/location/reset", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNoContent, resp.Code)
	})

	t.Run("GetStatisticsWithLocationAccessDenied", func(t *testing.T) {
		location := models.Location{
			City: "San Francisco",
		}
		body, _ := json.Marshal(location)

		req, _ := http.NewRequest("POST", "/location", bytes.NewReader(body))
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		req, _ = http.NewRequest("GET", "/statistics", nil)
		req.Header.Set("X-End-User-City", "New York")
		resp = httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("GetStatisticsWithLocationAccessAllowed", func(t *testing.T) {
		location := models.Location{
			City: "San Francisco",
		}
		body, _ := json.Marshal(location)

		req, _ := http.NewRequest("POST", "/location", bytes.NewReader(body))
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		req, _ = http.NewRequest("GET", "/statistics", nil)
		req.Header.Set("X-End-User-City", "San Francisco")
		resp = httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var stats models.Statistics
		json.Unmarshal(resp.Body.Bytes(), &stats)
		assert.GreaterOrEqual(t, stats.Count, int64(0))
	})

}
