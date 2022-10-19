package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
)

// A more end to end style test of the preparation and sending of the invoices.
func TestPrepareAndSendInvoices(t *testing.T) {
	testInvoices := Invoices{
		{TransactionId: 1234, LicensePlate: "AA12345", CheckInDate: "2022-08-26T15:50:24+02:00", CheckOutDate: "2022-08-26T15:50:24+02:00", Price: 3909, BusinessUserId: 4321, Type: 0},
		{TransactionId: 1235, LicensePlate: "AA12345", CheckInDate: "2022-08-26T15:50:24+02:00", CheckOutDate: "2022-08-26T15:50:24+02:00", Price: 3909, BusinessUserId: 4321, Type: 0},
		{TransactionId: 1236, LicensePlate: "AA12345", CheckInDate: "2022-08-26T15:50:24+02:00", CheckOutDate: "2022-08-26T15:50:24+02:00", Price: 3909, BusinessUserId: 5321, Type: 1},
		{TransactionId: 1237, LicensePlate: "AA12345", CheckInDate: "2022-08-26T15:50:24+02:00", CheckOutDate: "2022-08-26T15:50:24+02:00", Price: 3909, BusinessUserId: 5321, Type: 1},
	}

	// set up test server
	r := mux.NewRouter()
	r.HandleFunc("/{userId}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := mux.Vars(r)["userId"]

		if userId == strconv.FormatInt(4321, 10) {
			// If user with json type
			var invoices JSONInvoices
			err := json.NewDecoder(r.Body).Decode(&invoices)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// if the recieved match what we expect.. return OK
			if invoices.Nonpaids[0].TransactionId == testInvoices[0].TransactionId && invoices.Nonpaids[1].TransactionId == testInvoices[1].TransactionId {
				w.WriteHeader(http.StatusOK)
				return
			}

		} else if userId == strconv.FormatInt(5321, 10) {
			// If user with xml type
			var invoices XMLInvoices
			err := xml.NewDecoder(r.Body).Decode(&invoices)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// if the recieved match what we expect.. return OK
			if invoices.Invoices[0].TransactionId == testInvoices[2].TransactionId && invoices.Invoices[1].TransactionId == testInvoices[3].TransactionId {
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		w.WriteHeader(http.StatusBadRequest)
	}))

	s := httptest.NewServer(r)
	defer s.Close()

	// Check that we send error free
	_, errList := prepareAndSendInvoices(s.URL+"/", testInvoices)

	if len(errList) > 0 {
		for _, err := range errList {
			fmt.Println(err)
		}
		t.Error("non-empty error list")
	}
}
