package xero

// InvoiceResponse - The object that comes back from the API when we query for invoices.
type InvoiceResponse struct {
	Response
	Invoices []Invoice `json:"Invoices"`
}

// Invoice - As represented in the API. Some fields are ommitted for simplicity.
type Invoice struct {
	Type          string `json:"Type"`
	InvoiceID     string `json:"InvoiceID"`
	InvoiceNumber string `json:"InvoiceNumber"`
	Reference     string `json:"Reference"`
	// Payments      string `json:"Payments"`
	CreditNotes []InvoiceCreditNote `json:"CreditNotes"`
}

// InvoiceCreditNote - For credit notes associated with an invoice. Some fields have been ommitted for simplicity.
/*
   "CreditNoteID": "50e98404-2fba-4031-af67-8ba4bb227c44",
   "CreditNoteNumber": "CR1005",
   "ID": "50e98404-2fba-4031-af67-8ba4bb227c44",
   "HasErrors": false,
   "AppliedAmount": 550.00,
   "DateString": "2020-03-24T00:00:00",
   "Date": "\/Date(1585008000000+0000)\/",
   "LineItems": [],
   "Total": 550.00
*/
type InvoiceCreditNote struct {
	CreditNoteID     string  `json:"CreditNoteID"`
	CreditNoteNumber string  `json:"CreditNoteNumber"`
	ID               string  `json:"ID"`
	HasErrors        bool    `json:"HasErrors"`
	AppliedAmount    float32 `json:"AppliedAmount"`
	DateString       string  `json:"DateString"`
	Date             string  `json:"Date"`
	Total            float32 `json:"Total"`
}
