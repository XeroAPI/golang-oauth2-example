package xero

// InvoiceResponse - The object that comes back from the API when we query for invoices.
type InvoiceResponse struct {
	Response
	Invoices []Invoice `json:"Invoices"`
}

// Invoice - As represented in the API. Some fields are ommitted for simplicity.
type Invoice struct {
	Type           string  `json:"Type"`
	InvoiceID      string  `json:"InvoiceID"`
	InvoiceNumber  string  `json:"InvoiceNumber"`
	Reference      string  `json:"Reference"`
	AmountDue      float32 `json:"AmountDue"`
	AmountPaid     float32 `json:"AmountPaid"`
	AmountCredited float32 `json:"AmountCredited"`
	Status         string  `json:"Status"`
	Contact        Contact `json:"Contact"`
	Total          float32 `json:"Total"`
}
