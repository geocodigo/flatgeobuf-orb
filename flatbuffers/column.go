package flatbuffers

import (
	"github.com/gogama/flatgeobuf/flatgeobuf/flat"
	flatbuffers "github.com/google/flatbuffers/go"
)

type Column struct {
	Name        string
	ColumnType  flat.ColumnType
	Title       string
	Description string
	Width       int32
	Precision   int32
	Scale       int32
	Nullable    bool
	Unique      bool
	PrimaryKey  bool
	Metadata    string
}

func ColumnFromFlat(flatColumn *flat.Column) (*Column, error) {
	var column Column
	err := SafeInteraction(func() error {
		column.Name = string(flatColumn.Name())
		column.ColumnType = flatColumn.Type()
		column.Title = string(flatColumn.Title())
		column.Description = string(flatColumn.Description())
		column.Width = flatColumn.Width()
		column.Precision = flatColumn.Precision()
		column.Scale = flatColumn.Scale()
		column.Nullable = flatColumn.Nullable()
		column.Unique = flatColumn.Unique()
		column.PrimaryKey = flatColumn.PrimaryKey()
		column.Metadata = string(flatColumn.Metadata())

		return nil
	})

	return &column, err
}

func BuilderAddColumn(builder *flatbuffers.Builder, column Column) (flatbuffers.UOffsetT, error) {
	err := SafeInteraction(func() error {
		if column.Name != "" {
			offset := builder.CreateString(column.Name)
			defer flat.ColumnAddName(builder, offset)
		}

		defer flat.ColumnAddType(builder, column.ColumnType)

		if column.Title != "" {
			offset := builder.CreateString(column.Title)
			defer flat.ColumnAddTitle(builder, offset)
		}

		if column.Description != "" {
			offset := builder.CreateString(column.Description)
			defer flat.ColumnAddDescription(builder, offset)
		}

		defer flat.ColumnAddWidth(builder, column.Width)
		defer flat.ColumnAddPrecision(builder, column.Precision)
		defer flat.ColumnAddScale(builder, column.Scale)
		defer flat.ColumnAddNullable(builder, column.Nullable)
		defer flat.ColumnAddUnique(builder, column.Unique)
		defer flat.ColumnAddPrimaryKey(builder, column.PrimaryKey)

		if column.Metadata != "" {
			offset := builder.CreateString(column.Metadata)
			defer flat.ColumnAddMetadata(builder, offset)
		}

		flat.ColumnStart(builder)

		return nil
	})
	if err != nil {
		return 0, err
	}

	return flat.ColumnEnd(builder), nil
}

func BuilderAddColumns(builder *flatbuffers.Builder, columns []Column, startColumnsFunc StartVectorFunc) (flatbuffers.UOffsetT, error) {
	offsets := make([]flatbuffers.UOffsetT, len(columns))
	for i, column := range columns {
		offset, err := BuilderAddColumn(builder, column)
		if err != nil {
			return 0, err
		}
		offsets[i] = offset
	}

	return BuilderAddVector(builder, offsets, startColumnsFunc, builder.PrependUOffsetT), nil
}
