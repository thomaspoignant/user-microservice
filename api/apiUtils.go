package api

const NotFound string = "NOT_FOUND"
const Success string = "SUCCESS"

type ApiErrorResponse struct {
	Error string `json:"error"`
}
