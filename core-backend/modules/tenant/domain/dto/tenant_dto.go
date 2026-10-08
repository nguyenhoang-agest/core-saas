package dto

type CreateTenantRequest struct {
	Name        string `json:"name"                       validate:"required,alphanum,min=3,max=100"`
	Description string `json:"description,omitempty"      validate:"max=200"`
	
}
