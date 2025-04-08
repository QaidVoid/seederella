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
	case "name":
		return gofakeit.Name(), nil
	case "email":
		return gofakeit.Email(), nil
	case "sentence":
		return gofakeit.Sentence(5), nil
	case "uuid":
		return gofakeit.UUID(), nil
	case "int":
		return gofakeit.Int64(), nil
	case "bool":
		return gofakeit.Bool(), nil
	default:
		return nil, fmt.Errorf("unsupported faker field: %s", field)
	}
}
