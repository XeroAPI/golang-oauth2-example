package xero

// Connection - Represents an organisation as it appears in the /connections endpoint
type Connection struct {
	ID             string `json:"id"`
	TenantID       string `json:"tenantId"`
	TenantType     string `json:"tenantType"`
	TenantName     string `json:"tenantName"`
	CreateDateUTC  string `json:"createdDateUtc"`
	UpdatedDateUTC string `json:"updatedDateUtc"`
}
