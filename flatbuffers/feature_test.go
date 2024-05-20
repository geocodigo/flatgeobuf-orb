package flatbuffers

import (
	"testing"

	"github.com/gogama/flatgeobuf/flatgeobuf/flat"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFeature(t *testing.T) {
	type args struct {
		builder *flatbuffers.Builder
		feature Feature
	}
	tests := []struct {
		name    string
		args    args
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
				feature: Feature{
					Geometry:   &Geometry{GeometryType: flat.GeometryTypePolygon},
					Properties: []byte("some-properties"),
					Columns: []Column{
						{ColumnType: flat.ColumnTypeString},
						{ColumnType: flat.ColumnTypeBool},
						{ColumnType: flat.ColumnTypeDouble},
					},
				},
			},
			wantErr: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Add Feature to builder
			offset, err := BuilderAddFeature(tt.args.builder, tt.args.feature)
			tt.wantErr(t, err)

			tt.args.builder.FinishSizePrefixed(offset)
			bytes := tt.args.builder.FinishedBytes()

			// Get Feature from flatbuffers
			feature, err := FeatureFromFlat(flat.GetSizePrefixedRootAsFeature(bytes, 0))
			tt.wantErr(t, err)

			assert.Equal(t, tt.args.feature, *feature)
		})
	}
}
