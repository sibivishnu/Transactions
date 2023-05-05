package controllers

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/sibivishnu/Transactions/models"
)

var transactions = make([]*models.Transaction, 0)
var transactionsMutex = &sync.Mutex{}

var location *models.Location
var locationMutex = &sync.Mutex{}

func PostTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()
	if transaction.Timestamp.After(now) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if now.Sub(transaction.Timestamp).Seconds() > 60 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	transactionsMutex.Lock()
	transactions = append(transactions, &transaction)
	transactionsMutex.Unlock()

	w.WriteHeader(http.StatusCreated)
}

func DeleteTransactions(w http.ResponseWriter, r *http.Request) {
	transactionsMutex.Lock()
	transactions = make([]*models.Transaction, 0)
	transactionsMutex.Unlock()
	w.WriteHeader(http.StatusNoContent)
}

func getStatistics(w http.ResponseWriter, r *http.Request) {

	locationMutex.Lock()
	if location != nil {
		if r.Header.Get("X-End-User-City") != location.City {
			w.WriteHeader(http.StatusUnauthorized)
			locationMutex.Unlock()
			return
		}
	}
	locationMutex.Unlock()

	// Calculate statistics in O(1)
	now := time.Now().UTC()
	var sum, max, min float64
	var count int64
	transactionsMutex.Lock()
	for _, transaction := range transactions {
		if now.Sub(transaction.Timestamp).Seconds() <= 60 {
			sum += transaction.Amount
			count++
			if transaction.Amount > max {
				max = transaction.Amount
			}
			if count == 1 || transaction.Amount < min {
				min = transaction.Amount
			}
		}
	}
	transactionsMutex.Unlock()

	var avg float64
	if count > 0 {
		avg = sum / float64(count)
	}

	stats := &models.Statistics{
		Sum:   sum,
		Avg:   avg,
		Max:   max,
		Min:   min,
		Count: count,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func SetLocation(w http.ResponseWriter, r *http.Request) {
	var loc models.Location
	if err := json.NewDecoder(r.Body).Decode(&loc); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	locationMutex.Lock()
	location = &loc
	locationMutex.Unlock()

	w.WriteHeader(http.StatusCreated)
}

func ResetLocation(w http.ResponseWriter, r *http.Request) {
	locationMutex.Lock()
	location = nil
	locationMutex.Unlock()

	w.WriteHeader(http.StatusNoContent)
}
