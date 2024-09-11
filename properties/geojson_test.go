package properties

import (
	"fmt"
	"testing"

	"github.com/geocodigo/flatgeobuf-orb/flatbuffers"
	"github.com/gogama/flatgeobuf/flatgeobuf"
	"github.com/gogama/flatgeobuf/flatgeobuf/flat"
	fb "github.com/google/flatbuffers/go"
	"github.com/paulmach/orb/geojson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromGeoJSON(t *testing.T) {
	type args struct {
		properties geojson.Properties
		schema     flatgeobuf.Schema
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "empty",
			args: args{
				properties: geojson.Properties{},
				schema:     simpleSchema{},
			},
			want:    nil,
			wantErr: require.NoError,
		},
		{
			name: "all types",
			args: args{
				properties: geojson.Properties{
					"bool":    true,
					"int8":    int8(1),
					"uint8":   uint8(2),
					"int16":   int16(3),
					"uint16":  uint16(4),
					"int32":   int32(5),
					"uint32":  uint32(6),
					"int64":   int64(7),
					"uint64":  uint64(8),
					"float32": float32(9),
					"float64": float64(10),
					"string":  "property",
					"binary":  []byte("property"),
				},
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
			want: []byte{
				0x0, 0x0, // Column 0
				0x1,      // true
				0x1, 0x0, // Column 1
				0x1,      // int8(1)
				0x2, 0x0, // Column 2
				0x2,      // uint8(2)
				0x3, 0x0, // Column 3
				0x3, 0x0, // int16(3)
				0x4, 0x0, // Column 4
				0x4, 0x0, // uint16(4)
				0x5, 0x0, // Column 5
				0x5, 0x0, 0x0, 0x0, // int32(5)
				0x6, 0x0, // Column 6
				0x6, 0x0, 0x0, 0x0, // uint32(6)
				0x7, 0x0, // Column 7
				0x7, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, // int64(7)
				0x8, 0x0, // Column 8
				0x8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, // uint64(8)
				0x9, 0x0, // Column 9
				0x0, 0x0, 0x10, 0x41, // float32(9)
				0xa, 0x0, // Column 10
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x24, 0x40, // float64(10)
				0xb, 0x0, // Column 11
				0x8, 0x0, 0x0, 0x0, 0x70, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x79, // "property"
				0xc, 0x0, // Column 12
				0x8, 0x0, 0x0, 0x0, 0x70, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x79, // []byte("property")
			},
			wantErr: require.NoError,
		},
		{
			name: "skip missing propertry",
			args: args{
				properties: geojson.Properties{
					"present": true,
				},
				schema: simpleSchema{
					flatbuffers.Column{Name: "present", ColumnType: flat.ColumnTypeBool},
					flatbuffers.Column{Name: "missing", ColumnType: flat.ColumnTypeBool},
				},
			},
			want: []byte{
				0x0, 0x0, // Column 0
				0x01, // true
			},
			wantErr: require.NoError,
		},
		{
			name: "invalid schema",
			args: args{
				properties: geojson.Properties{
					"bool": true,
				},
				schema: errSchema{
					flatbuffers.Column{Name: "bool", ColumnType: flat.ColumnTypeBool},
				},
			},
			want: nil,
			wantErr: func(tt require.TestingT, err error, i ...any) {
				require.ErrorContains(t, err, "failed to get columns from schema")
			},
		},
		{
			name: "invalid property type",
			args: args{
				properties: geojson.Properties{
					"invalid": struct{ invalid bool }{invalid: true},
				},
				schema: simpleSchema{
					flatbuffers.Column{Name: "invalid", ColumnType: flat.ColumnTypeJson},
				},
			},
			want: nil,
			wantErr: func(tt require.TestingT, err error, i ...any) {
				require.ErrorContains(t, err, "invalid flatgeobuf property type")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create new properties
			got, err := FromGeoJSON(tt.args.properties, tt.args.schema)

			// Assert written properties
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

var (
	_ flatgeobuf.Schema = simpleSchema{}
	_ flatgeobuf.Schema = errSchema{}
)

type simpleSchema []flatbuffers.Column

func (s simpleSchema) ColumnsLength() int {
	return len(s)
}

func (s simpleSchema) Columns(obj *flat.Column, i int) bool {
	if i < 0 || len(s) <= i {
		return false
	}

	builder := fb.NewBuilder(0)
	offset, _ := flatbuffers.BuilderAddColumn(builder, s[i])
	builder.FinishSizePrefixed(offset)
	*obj = *flat.GetSizePrefixedRootAsColumn(builder.FinishedBytes(), 0)

	return true
}

type errSchema []flatbuffers.Column

func (e errSchema) ColumnsLength() int {
	return len(e)
}

func (e errSchema) Columns(obj *flat.Column, i int) bool {
	panic(fmt.Errorf("failed to get column at index %d", i))
}
