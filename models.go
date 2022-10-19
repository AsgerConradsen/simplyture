package main

import "encoding/xml"

type Invoices []struct {
	TransactionId  uint64 `json:"transactionId"`
	LicensePlate   string `json:"licensePlate"`
	CheckInDate    string `json:"checkInDate"`
	CheckOutDate   string `json:"checkOutDate"`
	Price          uint64 `json:"price"`
	BusinessUserId uint64 `json:"businessUserId"`
	Type           uint64 `json:"type"`
}

type InvoiceDTO struct {
	TransactionId  uint64 `json:"id" xml:"UniqueId"`
	LicensePlate   string `json:"vrm" xml:"RegistrationNumber"`
	CheckInDate    string `json:"eventStartDate" xml:"DriveInDate"`
	CheckOutDate   string `json:"eventEndDate" xml:"DriveOutDate"`
	Price          uint64 `json:"price" xml:"Price"`
	BusinessUserId uint64 `json:"facilityId" xml:"CarparkId"`
}

type JSONInvoices struct {
	Nonpaids []InvoiceDTO `json:"nonpaids"`
}

type XMLInvoices struct {
	XMLName  xml.Name
	Invoices []InvoiceDTO `xml:"Invoice"`
}

type Receipt struct {
	UserId uint64
	Msg    string
}
