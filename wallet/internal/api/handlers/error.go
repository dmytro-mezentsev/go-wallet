package handlers

import "fmt"

func ErrorResponse(errorMessage string) string {
	return fmt.Sprintf(`{"error": "%s"}`, errorMessage)
}

func InternalErrorResponse() string {
	return ErrorResponse("Internal server error")
}
