package xero

import "errors"

// OrganisationResponse - The response from the /organisation endpoint.
type OrganisationResponse struct {
	ID            string         `json:"id"`
	Status        string         `json:"Status"`
	ProviderName  string         `json:"ProviderName"`
	DateTimeUTC   string         `json:"DateTimeUTC"`
	Organisations []Organisation `json:"Organisations"`
}

// GetOrgByID - Returns an organisation by ID.
func (r *OrganisationResponse) GetOrgByID(id string) (*Organisation, error) {
	for _, org := range r.Organisations {
		if org.OrganisationID == id {
			return &org, nil
		}
	}
	return &Organisation{}, errors.New("Unable to find organisation with ID " + id)
}

// Organisation - As it is received from the API.
type Organisation struct {
	Name                   string                `json:"Name"`
	LegalName              string                `json:"LegalName"`
	PaysTax                bool                  `json:"PaysTax"`
	Version                string                `json:"Version"`
	OrganisationType       string                `json:"OrganisationType"`
	BaseCurrency           string                `json:"BaseCurrency"`
	CountryCode            string                `json:"CountryCode"`
	IsDemoCompany          bool                  `json:"IsDemoCompany"`
	OrganisationStatus     string                `json:"OrganisationStatus"`
	RegistrationNumber     string                `json:"RegistrationNumber"`
	TaxNumber              string                `json:"TaxNumber"`
	FinancialYearEndDay    int                   `json:"FinancialYearEndDay"`
	FinancialYearEndMonth  int                   `json:"FinancialYearEndMonth"`
	SalesTaxBasis          string                `json:"SalesTaxBasis"`
	SalesTaxPeriod         string                `json:"SalesTaxPeriod"`
	DefaultSalesTax        string                `json:"DefaultSalesTax"`
	DefaultPurchasesTax    string                `json:"DefaultPurchasesTax"`
	PeriodLockDate         string                `json:"PeriodLockDate"`
	CreatedDateUTC         string                `json:"CreatedDateUTC"`
	OrganisationEntityType string                `json:"OrganisationEntityType"`
	Timezone               string                `json:"Timezone"`
	ShortCode              string                `json:"ShortCode"`
	OrganisationID         string                `json:"OrganisationID"`
	Edition                string                `json:"Edition"`
	Class                  string                `json:"Class"`
	Addresses              []OrganisationAddress `json:"Addresses"`
	// ExternalLinks and PaymentTerms not included here because this information is not present in the Demo Company™.
	// ExternalLinks: [],
	// PaymentTerms: {}
}

// OrganisationAddress - Addresses associated with an organisation
type OrganisationAddress struct {
	AddressType  string `json:"AddressType"`
	AddressLine1 string `json:"AddressLine1"`
	City         string `json:"City"`
	Region       string `json:"Region"`
	PostalCode   string `json:"PostalCode"`
	Country      string `json:"Country"`
	AttentionTo  string `json:"AttentionTo"`
}

// OrganisationPhone - Phone numbers associated with an organisation
type OrganisationPhone struct {
	PhoneType   string `json:"PhoneType"`
	PhoneNumber string `json:"PhoneNumber"`
}
