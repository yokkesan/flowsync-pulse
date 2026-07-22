package company

type ErrorResponse struct {
	Message string `json:"message"`
}

type CreateRequest struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
	Slug string `json:"slug" binding:"required,min=2,max=100"`
}

type CreateResponse struct {
	CompanyID uint64 `json:"company_id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
}
