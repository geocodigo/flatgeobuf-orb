package properties

import (
	"bytes"
	"fmt"

	"github.com/geocodigo/flatgeobuf-orb/flatbuffers"
	"github.com/gogama/flatgeobuf/flatgeobuf"
	"github.com/paulmach/orb/geojson"
)

func FromGeoJSON(props geojson.Properties, schema flatgeobuf.Schema) ([]byte, error) {
	var propsBuffer bytes.Buffer
	propsWriter := flatgeobuf.NewPropWriter(&propsBuffer)

	var columns []flatbuffers.Column
	if err := flatbuffers.SafeInteraction(func() error {
		cols, err := flatbuffers.VectorFromFlat(
			schema.ColumnsLength(),
			schema.Columns,
			flatbuffers.ColumnFromFlat,
		)
		columns = cols
		return err
	}); err != nil {
		return nil, fmt.Errorf("failed to get columns from schema: %w", err)
	}

	var err error
	for i, col := range columns {
		prop, ok := props[col.Name]
		if !ok {
			// TODO: Need to confirm desired behaviour when property in schema is not present .
			// Skip for now.
			continue
		}

		// The flatgeobuf spec requires props are prefixed by a ushort key,
		// using the property's index as a key.
		_, err = propsWriter.WriteUShort(uint16(i))
		if err != nil {
			break
		}

		switch prop := prop.(type) {
		case bool:
			_, err = propsWriter.WriteBool(prop)
		case int8:
			_, err = propsWriter.WriteByte(prop)
		case uint8:
			_, err = propsWriter.WriteUByte(prop)
		case int16:
			_, err = propsWriter.WriteShort(prop)
		case uint16:
			_, err = propsWriter.WriteUShort(prop)
		case int32:
			_, err = propsWriter.WriteInt(prop)
		case uint32:
			_, err = propsWriter.WriteUInt(prop)
		case int64:
			_, err = propsWriter.WriteLong(prop)
		case uint64:
			_, err = propsWriter.WriteULong(prop)
		case float32:
			_, err = propsWriter.WriteFloat(prop)
		case float64:
			_, err = propsWriter.WriteDouble(prop)
		case string:
			_, err = propsWriter.WriteString(prop)
		case []byte:
			_, err = propsWriter.WriteBinary(prop)
		default:
			err = fmt.Errorf("invalid flatgeobuf property type: %T", prop)
		}

		if err != nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	return propsBuffer.Bytes(), nil
}
