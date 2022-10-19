package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// Fetches all the non-empty page and stops when reaching an empty page.
// Return the pages. If there is an error on of the pages, it will return that error
// and the contents of all previous pages.
func getAllInvoices(url string) (invoices Invoices, err error) {
	var page Invoices
	pageNum := 0

	// This is just a do-while loop. They look a little weird in go. At least to me...
	for ok := true; ok; ok = len(page) > 0 {
		page, err = getPage(pageNum, url)
		if err != nil {
			return
		}

		// Appending is amortized O(1) in go so this is not too bad.
		invoices = append(invoices, page...)
		pageNum++
	}
	return
}

// Fetches page number 'pageNum' and returns the contents.
// Treats API response of non-200 as error
func getPage(pageNum int, url string) (invoices Invoices, err error) {
	postBody, err := json.Marshal(map[string]int{
		"page": pageNum,
	})
	if err != nil {
		return
	}

	bodyBuffer := bytes.NewBuffer(postBody)
	resp, err := http.Post(url, "application/json", bodyBuffer)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New("get page failed on page " + strconv.Itoa(pageNum) + " with error " + resp.Status)
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&invoices)
	return
}
