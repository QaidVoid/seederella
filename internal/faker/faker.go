package faker

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
)

func Init() {
	gofakeit.Seed(0) // seed once at app start
}

func generateParagraph(params []string) string {
	paragraphCount := 2
	sentenceCount := 5
	wordCount := 10
	separator := " "

	if len(params) > 0 {
		if val, err := strconv.Atoi(params[0]); err == nil {
			paragraphCount = val
		}
	}
	if len(params) > 1 {
		if val, err := strconv.Atoi(params[1]); err == nil {
			sentenceCount = val
		}
	}
	if len(params) > 2 {
		if val, err := strconv.Atoi(params[2]); err == nil {
			wordCount = val
		}
	}
	if len(params) > 3 {
		separator = params[3]
	}
	return gofakeit.Paragraph(paragraphCount, sentenceCount, wordCount, separator)
}

func Generate(fieldSpec string) (any, error) {
	parts := strings.Split(fieldSpec, ":")
	field := parts[0]

	var params []string
	if len(parts) > 1 {
		params = strings.Split(parts[1], ";")
	}

	switch field {
	case "first_name":
		return gofakeit.FirstName(), nil
	case "last_name":
		return gofakeit.LastName(), nil
	case "name":
		return gofakeit.Name(), nil
	case "word":
		return gofakeit.Word(), nil
	case "username":
		return gofakeit.Username(), nil
	case "email":
		return gofakeit.Email(), nil
	case "phone":
		return gofakeit.Phone(), nil
	case "credit_card":
		return gofakeit.CreditCardNumber(nil), nil
	case "address":
		return gofakeit.Address().Address, nil
	case "city":
		return gofakeit.City(), nil
	case "country":
		return gofakeit.Country(), nil
	case "zip":
		return gofakeit.Zip(), nil
	case "sentence":
		count := 5
		if len(params) > 0 {
			if val, err := strconv.Atoi(params[0]); err == nil {
				count = val
			}
		}
		return gofakeit.Sentence(count), nil
	case "paragraph":
		return generateParagraph(params), nil
	case "uuid":
		return gofakeit.UUID(), nil
	case "int":
		return gofakeit.Int64(), nil
	case "float":
		return gofakeit.Float64(), nil
	case "bool":
		return gofakeit.Bool(), nil
	case "date_future":
		return gofakeit.FutureDate(), nil
	case "date_past":
		return gofakeit.PastDate(), nil
	default:
		return nil, fmt.Errorf("unsupported faker field: %s", fieldSpec)
	}
}
