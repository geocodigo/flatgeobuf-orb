package flatbuffers

import (
	"testing"

	"github.com/gogama/flatgeobuf/flatgeobuf/flat"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeader(t *testing.T) {
	type args struct {
		builder *flatbuffers.Builder
		header  Header
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
				header: Header{
					Name:         "some-name",
					Envelope:     []float64{1, 2, 3, 4},
					GeometryType: flat.GeometryTypePoint,
					HasZ:         true,
					HasM:         true,
					HasT:         true,
					HasTM:        true,
					Columns: []Column{
						{Name: "some-column-1"},
						{Name: "some-column-2"},
						{Name: "some-column-3"},
					},
					FeaturesCount: 123,
					IndexNodeSize: 456,
					CRS:           &CRS{Name: "some-crs"},
					Title:         "some-title",
					Description:   "some-description",
					Metadata:      "some-metadata",
				},
			},
			wantErr: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Add header to builder
			offset, err := BuilderAddHeader(tt.args.builder, tt.args.header)
			tt.wantErr(t, err)

			tt.args.builder.FinishSizePrefixed(offset)
			bytes := tt.args.builder.FinishedBytes()

			// Get header from flatbuffers
			header, err := HeaderFromFlat(flat.GetSizePrefixedRootAsHeader(bytes, 0))
			tt.wantErr(t, err)

			assert.Equal(t, tt.args.header, *header)
		})
	}
}
