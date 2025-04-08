package faker

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
)

func Init() {
	gofakeit.Seed(0) // seed once at app start
}

func Generate(field string) (any, error) {
	switch field {
	case "first_name":
		return gofakeit.FirstName(), nil
	case "last_name":
		return gofakeit.LastName(), nil
	case "name":
		return gofakeit.Name(), nil
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
		return gofakeit.Sentence(5), nil
	case "paragraph":
		return gofakeit.Paragraph(2, 5, 10, " "), nil
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
		return nil, fmt.Errorf("unsupported faker field: %s", field)
	}
}
