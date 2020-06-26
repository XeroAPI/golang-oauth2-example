package xero

// "ContactID": "571a2414-81ff-4f8f-8498-d91d83793131",
// "Name": "Bank West",
// "Addresses": [],
// "Phones": [],
// "ContactGroups": [],
// "ContactPersons": [],
// "HasValidationErrors": false

// Contact - Necessary for reading contact information out of the invoice information. Some fields ommitted for
// siplicity.
type Contact struct {
	ContactID string `json:"ContactID"`
	Name      string `json:"Name"`
}
