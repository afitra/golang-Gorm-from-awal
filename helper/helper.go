package helper

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func ApiResponse(Message string, Code int, Status string, Data interface{}) Response {
	meta := Meta{
		Message,
		Code,
		Status,
	}

	jsonResponse := Response{
		Meta: meta,
		Data: Data,
	}

	return jsonResponse
}

func FormatValidationError(err error) []string {

	var errors []string
	if err.Error() == "EOF" {
		errors = append(errors, "error input data")
		// return ["error input data"]
		return errors
	}
	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}
	return errors
}
func GoDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
