package main

import "fmt"

const fetchUrl = "https://mszaoj3b7r7ha6c7spvlp5xyva0uyubb.lambda-url.eu-west-1.on.aws/"
const userUrl = "https://7by22ubtoqet4e7wgjwo64z25e0wgkgp.lambda-url.eu-west-1.on.aws/"

func main() {
	invoices, err := getAllInvoices(fetchUrl)
	if err != nil {
		// If something went wrong in the process of getting the invoices, we stop execution.
		panic("error in getting invoices: " + err.Error())
	}

	fmt.Println("Successfully fetched all invoices")

	receipts, errList := prepareAndSendInvoices(userUrl, invoices)
	for _, err := range errList {
		fmt.Println(err)
	}

	for _, receipt := range receipts {
		fmt.Println(receipt.UserId, " responded with ", receipt.Msg)
	}

}
