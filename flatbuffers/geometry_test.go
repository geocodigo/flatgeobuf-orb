package flatbuffers

import (
	"testing"

	"github.com/gogama/flatgeobuf/flatgeobuf/flat"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGeometry(t *testing.T) {
	type args struct {
		builder  *flatbuffers.Builder
		geometry Geometry
	}
	tests := []struct {
		name    string
		args    args
		want    flatbuffers.UOffsetT
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "empty fields",
			args: args{
				builder: flatbuffers.NewBuilder(0),
			},
			wantErr: require.NoError,
		},
		{
			name: "all fields",
			args: args{
				builder: flatbuffers.NewBuilder(0),
				geometry: Geometry{
					Ends:         []uint32{0, 1, 2, 3},
					XY:           []float64{0, 1, 2},
					Z:            []float64{0, 1, 2},
					M:            []float64{0, 1, 2},
					T:            []float64{0, 1, 2},
					TM:           []uint64{0, 1, 2},
					GeometryType: flat.GeometryTypeGeometryCollection,
					Parts: []Geometry{
						{GeometryType: flat.GeometryTypeMultiPoint},
						{GeometryType: flat.GeometryTypeLineString},
						{GeometryType: flat.GeometryTypeMultiPolygon},
					},
				},
			},
			wantErr: require.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Add geometry to builder
			offset, err := BuilderAddGeometry(tt.args.builder, tt.args.geometry)
			tt.wantErr(t, err)

			tt.args.builder.FinishSizePrefixed(offset)
			bytes := tt.args.builder.FinishedBytes()

			// Get geometry from flatbuffers
			geometry, err := GeometryFromFlat(flat.GetSizePrefixedRootAsGeometry(bytes, 0))
			tt.wantErr(t, err)

			assert.Equal(t, tt.args.geometry, *geometry)
		})
	}
}
