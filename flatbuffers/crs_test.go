package flatbuffers

import (
	"testing"

	"github.com/gogama/flatgeobuf/flatgeobuf/flat"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCRS(t *testing.T) {
	type args struct {
		builder *flatbuffers.Builder
		crs     CRS
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
				crs: CRS{
					Org:         "some-org",
					Code:        42,
					Name:        "some-name",
					Description: "some-description",
					WKT:         "some-wkt",
					CodeString:  "some-code-string",
				},
			},
			wantErr: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Add CRS to builder
			offset, err := BuilderAddCRS(tt.args.builder, tt.args.crs)
			tt.wantErr(t, err)

			tt.args.builder.FinishSizePrefixed(offset)
			bytes := tt.args.builder.FinishedBytes()

			// Get CRS from flatbuffers
			crs, err := CRSFromFlat(flat.GetSizePrefixedRootAsCrs(bytes, 0))
			tt.wantErr(t, err)

			assert.Equal(t, tt.args.crs, *crs)
		})
	}
}
