package xero

// Organisation - Represents an organisation as it appears in the /connections endpoint
type Organisation struct {
	ID             string `json:"id"`
	TenantID       string `json:"tenantId"`
	TenantType     string `json:"tenantType"`
	TenantName     string `json:"tenantName"`
	CreateDateUTC  string `json:"createdDateUtc"`
	UpdatedDateUTC string `json:"updatedDateUtc"`
}
