package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Takes a url and a slice of invoices
// Sorts them by user id, prepare them for sending (by converting to DTO object) and sends them to each user
// IMPORTANT: does not stop running on errors in the send process. If sending fails for 1 user, it is noted and it moves on.
// returns a receipt for all the succesful sends and an error for each unsuccessful send.
func prepareAndSendInvoices(url string, invoices Invoices) (receipts []Receipt, errorList []error) {

	// Use a map to sort by user id. Key is user id, val is that users invoices
	invoiceUserMap := make(map[uint64]Invoices)
	for _, invoice := range invoices {
		val, exists := invoiceUserMap[invoice.BusinessUserId]
		if exists {
			// I am assuming there is only one invoice return type per users. Seems to be the case in your test data and also makes sense. This will return error if that doesn't hold.
			// This is boundary safe btw. We only init with non-empty slice.
			if val[0].Type != invoice.Type {
				errorList = append(errorList, errors.New("user with id "+strconv.FormatUint(invoice.BusinessUserId, 10)+" had multiple invoice types."))
				return
			}

			invoiceUserMap[invoice.BusinessUserId] = append(val, invoice)
		} else {
			invoiceUserMap[invoice.BusinessUserId] = Invoices{invoice}
		}
	}

	// For each user, convert all of their invoices to our invoice DTO object. Then send them their invoices in their requested format.
	// The DTO object has different naming for JSON and XML so it marshals according to the task.
	for _, val := range invoiceUserMap {
		//Convert to DTO objects
		var dtos []InvoiceDTO
		for _, invoice := range val {
			dtos = append(dtos, InvoiceDTO{
				TransactionId:  invoice.TransactionId,
				LicensePlate:   invoice.LicensePlate,
				CheckInDate:    invoice.CheckInDate,
				CheckOutDate:   invoice.CheckOutDate,
				Price:          invoice.Price,
				BusinessUserId: invoice.BusinessUserId,
			})
		}

		// Send in the specified type
		invoiceType := val[0].Type
		if invoiceType == 0 {
			// JSON
			jsonInvoices := JSONInvoices{Nonpaids: dtos}
			jsonBody, err := json.Marshal(jsonInvoices)
			if err != nil {
				errorList = append(errorList, err)
				continue
			}

			receipt, err := sendInvoices(url, jsonBody, dtos[0].BusinessUserId, "application/json")
			if err != nil {
				errorList = append(errorList, err)
			}

			receipts = append(receipts, receipt)
		} else if invoiceType == 1 {
			// XML
			xmlInvoices := XMLInvoices{XMLName: xml.Name{Local: "Invoices"}, Invoices: dtos}
			xmlBody, err := xml.Marshal(xmlInvoices)
			if err != nil {
				errorList = append(errorList, err)
				continue
			}

			receipt, err := sendInvoices(url, xmlBody, dtos[0].BusinessUserId, "application/xml")
			if err != nil {
				errorList = append(errorList, err)
			}

			receipts = append(receipts, receipt)
		} else {
			// Unknown type
			errorList = append(errorList, errors.New("invalid invoice type"))
		}
	}
	return
}

// Takes an URL, body containing invoices, user id and contenttype (xml or json).
// Sends the body to the url with user id as param in the specified format.
// Returns receipt on success and non-nil error if something went wrong.
// Treats API response other than 200 as error
func sendInvoices(url string, invoicesBody []byte, userId uint64, contentType string) (receipt Receipt, err error) {
	bodyBuffer := bytes.NewBuffer(invoicesBody)
	resp, err := http.Post(url+strconv.FormatInt(int64(userId), 10), contentType, bodyBuffer)
	if err != nil {
		return
	}

	// Treat non-success codes as error
	if resp.StatusCode != 200 {
		err = errors.New("send operation for user: " + strconv.FormatInt(int64(userId), 10) + " failed with " + resp.Status)
		return
	}

	// Read body to check for receipt
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	receipt = Receipt{UserId: userId, Msg: string(body)}

	return
}
