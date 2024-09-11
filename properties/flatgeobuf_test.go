package properties

import (
	"bytes"
	"testing"

	"github.com/geocodigo/flatgeobuf-orb/flatbuffers"
	"github.com/gogama/flatgeobuf/flatgeobuf"
	"github.com/gogama/flatgeobuf/flatgeobuf/flat"
	"github.com/paulmach/orb/geojson"
	"github.com/stretchr/testify/require"
)

func TestFromFlatGeobuf(t *testing.T) {
	type args struct {
		propertiesBytes []byte
		schema          flatgeobuf.Schema
	}
	tests := []struct {
		name    string
		args    args
		want    geojson.Properties
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "empty",
			args: args{
				propertiesBytes: []byte{},
				schema:          simpleSchema{},
			},
			want:    make(geojson.Properties),
			wantErr: require.NoError,
		},
		{
			name: "all types",
			args: args{
				propertiesBytes: func() []byte {
					var propsBuffer bytes.Buffer
					propsWriter := flatgeobuf.NewPropWriter(&propsBuffer)

					propsWriter.WriteUShort(0) // Column 0
					propsWriter.WriteBool(true)
					propsWriter.WriteUShort(1) // Column 1
					propsWriter.WriteByte(1)
					propsWriter.WriteUShort(2) // Column 2
					propsWriter.WriteUByte(2)
					propsWriter.WriteUShort(3) // Column 3
					propsWriter.WriteShort(3)
					propsWriter.WriteUShort(4) // Column 4
					propsWriter.WriteUShort(4)
					propsWriter.WriteUShort(5) // Column 5
					propsWriter.WriteInt(5)
					propsWriter.WriteUShort(6) // Column 6
					propsWriter.WriteUInt(6)
					propsWriter.WriteUShort(7) // Column 7
					propsWriter.WriteLong(7)
					propsWriter.WriteUShort(8) // Column 8
					propsWriter.WriteULong(8)
					propsWriter.WriteUShort(9) // Column 9
					propsWriter.WriteFloat(9.0)
					propsWriter.WriteUShort(10) // Column 10
					propsWriter.WriteDouble(10.0)
					propsWriter.WriteUShort(11) // Column 11
					propsWriter.WriteString("property")
					propsWriter.WriteUShort(12) // Column 12
					propsWriter.WriteBinary([]byte("property"))

					return propsBuffer.Bytes()
				}(),
				schema: simpleSchema{
					flatbuffers.Column{Name: "bool", ColumnType: flat.ColumnTypeBool},
					flatbuffers.Column{Name: "int8", ColumnType: flat.ColumnTypeByte},
					flatbuffers.Column{Name: "uint8", ColumnType: flat.ColumnTypeUByte},
					flatbuffers.Column{Name: "int16", ColumnType: flat.ColumnTypeShort},
					flatbuffers.Column{Name: "uint16", ColumnType: flat.ColumnTypeUShort},
					flatbuffers.Column{Name: "int32", ColumnType: flat.ColumnTypeInt},
					flatbuffers.Column{Name: "uint32", ColumnType: flat.ColumnTypeUInt},
					flatbuffers.Column{Name: "int64", ColumnType: flat.ColumnTypeLong},
					flatbuffers.Column{Name: "uint64", ColumnType: flat.ColumnTypeULong},
					flatbuffers.Column{Name: "float32", ColumnType: flat.ColumnTypeFloat},
					flatbuffers.Column{Name: "float64", ColumnType: flat.ColumnTypeDouble},
					flatbuffers.Column{Name: "string", ColumnType: flat.ColumnTypeString},
					flatbuffers.Column{Name: "binary", ColumnType: flat.ColumnTypeBinary},
				},
			},
			want: geojson.Properties{
				"bool":    true,
				"int8":    int8(1),
				"uint8":   uint8(2),
				"int16":   int16(3),
				"uint16":  uint16(4),
				"int32":   int32(5),
				"uint32":  uint32(6),
				"int64":   int64(7),
				"uint64":  uint64(8),
				"float32": float32(9.0),
				"float64": float64(10.0),
				"string":  "property",
				"binary":  []byte("property"),
			},
			wantErr: require.NoError,
		},
		{
			name: "invalid schema",
			args: args{
				propertiesBytes: []byte{},
				schema: errSchema{
					flatbuffers.Column{Name: "bool", ColumnType: flat.ColumnTypeBool},
				},
			},
			want: make(geojson.Properties),
			wantErr: func(tt require.TestingT, err error, i ...any) {
				require.ErrorContains(t, err, "failed to read propetries")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create new properties
			got, err := FromFlatGeobuf(tt.args.propertiesBytes, tt.args.schema)

			// Assert written properties
			tt.wantErr(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
