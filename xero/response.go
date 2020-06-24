package xero

// Response is a generic top-level set of properties that we get from every resource endpoint.
type Response struct {
	ID           string `json:"id"`
	Status       string `json:"Status"`
	ProviderName string `json:"ProviderName"`
	DateTimeUTC  string `json:"DateTimeUTC"`
}
