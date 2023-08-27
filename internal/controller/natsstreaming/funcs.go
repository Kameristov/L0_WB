package natsstreaming

import (
	"L0_EVRONE/internal/aggregate"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

func Serialization(data []byte) (aggregate.Order, error) {
	
	orderData, err := aggregate.NewOrder(data)
	return orderData, err
}

func Validation(data []byte) error {

	schemaLoader := gojsonschema.NewStringLoader(JsonScheme)
	documentLoader := gojsonschema.NewStringLoader(string(data[:]))

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
    if err != nil {
        return fmt.Errorf("gojsonschema.Validate error: %v", err)
    }
	if !result.Valid() {
        return fmt.Errorf("JSON not valid")
    } 
	return nil
}
