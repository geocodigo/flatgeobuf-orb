package properties

import (
	"bytes"
	"fmt"

	"github.com/gogama/flatgeobuf/flatgeobuf"
	"github.com/paulmach/orb/geojson"
)

func FromFlatGeobuf(b []byte, schema flatgeobuf.Schema) (geojson.Properties, error) {
	propReader := flatgeobuf.NewPropReader(bytes.NewReader(b))

	pvs, err := propReader.ReadSchema(schema)
	if err != nil {
		return geojson.Properties{}, fmt.Errorf("failed to read propetries: %w", err)
	}

	props := make(geojson.Properties)
	for _, pv := range pvs {
		key := string(pv.Col.Name())
		props[key] = pv.Value
	}

	return props, nil
}
