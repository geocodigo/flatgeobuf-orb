package flatbuffers

import (
	"testing"

	"github.com/gogama/flatgeobuf/flatgeobuf/flat"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestColumn(t *testing.T) {
	type args struct {
		builder *flatbuffers.Builder
		column  Column
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
				column: Column{
					Name:        "some-name",
					ColumnType:  flat.ColumnTypeString,
					Title:       "some-title",
					Description: "some-description",
					Width:       1,
					Precision:   2,
					Scale:       3,
					Nullable:    true,
					Unique:      true,
					PrimaryKey:  true,
					Metadata:    "some-metadata",
				},
			},
			wantErr: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Add column to builder
			offset, err := BuilderAddColumn(tt.args.builder, tt.args.column)
			tt.wantErr(t, err)

			tt.args.builder.FinishSizePrefixed(offset)
			bytes := tt.args.builder.FinishedBytes()

			// Get column from flatbuffers
			column, err := ColumnFromFlat(flat.GetSizePrefixedRootAsColumn(bytes, 0))
			tt.wantErr(t, err)

			assert.Equal(t, tt.args.column, *column)
		})
	}
}
