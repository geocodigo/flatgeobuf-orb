package flatbuffers

import (
	"github.com/gogama/flatgeobuf/flatgeobuf/flat"
	flatbuffers "github.com/google/flatbuffers/go"
)

type Header struct {
	Name          string
	Envelope      []float64
	GeometryType  flat.GeometryType
	HasZ          bool
	HasM          bool
	HasT          bool
	HasTM         bool
	Columns       []Column
	FeaturesCount uint64
	IndexNodeSize uint16
	CRS           *CRS
	Title         string
	Description   string
	Metadata      string
}

func HeaderFromFlat(flatHeader *flat.Header) (*Header, error) {
	var header Header

	err := SafeInteraction(func() error {
		header.Name = string(flatHeader.Name())

		for i := 0; i < flatHeader.EnvelopeLength(); i++ {
			header.Envelope = append(header.Envelope, flatHeader.Envelope(i))
		}

		header.GeometryType = flatHeader.GeometryType()
		header.HasZ = flatHeader.HasZ()
		header.HasM = flatHeader.HasM()
		header.HasT = flatHeader.HasT()
		header.HasTM = flatHeader.HasTm()

		if cols, err := VectorFromFlat(flatHeader.ColumnsLength(), flatHeader.Columns, ColumnFromFlat); err != nil {
			return err
		} else {
			header.Columns = cols
		}

		header.FeaturesCount = flatHeader.FeaturesCount()
		header.IndexNodeSize = flatHeader.IndexNodeSize()

		var flatCRS flat.Crs
		if flatHeader.Crs(&flatCRS) != nil {
			crs, err := CRSFromFlat(&flatCRS)
			if err != nil {
				return err
			}
			header.CRS = crs
		}

		header.Title = string(flatHeader.Title())
		header.Description = string(flatHeader.Description())
		header.Metadata = string(flatHeader.Metadata())

		return nil
	})

	return &header, err
}

func BuilderAddHeader(builder *flatbuffers.Builder, header Header) (flatbuffers.UOffsetT, error) {
	err := SafeInteraction(func() error {
		if header.Name != "" {
			offset := builder.CreateString(header.Name)
			defer flat.HeaderAddName(builder, offset)
		}

		if header.Envelope != nil {
			offset := BuilderAddVector(builder, header.Envelope, flat.HeaderStartEnvelopeVector, builder.PrependFloat64)
			defer flat.HeaderAddEnvelope(builder, offset)
		}

		defer flat.HeaderAddGeometryType(builder, header.GeometryType)
		defer flat.HeaderAddHasZ(builder, header.HasZ)
		defer flat.HeaderAddHasM(builder, header.HasM)
		defer flat.HeaderAddHasT(builder, header.HasT)
		defer flat.HeaderAddHasTm(builder, header.HasTM)

		if header.Columns != nil {
			offset, err := BuilderAddColumns(builder, header.Columns, flat.HeaderStartColumnsVector)
			if err != nil {
				return err
			}
			defer flat.HeaderAddColumns(builder, offset)
		}

		defer flat.HeaderAddFeaturesCount(builder, header.FeaturesCount)
		defer flat.HeaderAddIndexNodeSize(builder, header.IndexNodeSize)

		if header.CRS != nil {
			offset, err := BuilderAddCRS(builder, *header.CRS)
			if err != nil {
				return err
			}
			defer flat.HeaderAddCrs(builder, offset)
		}

		if header.Title != "" {
			offset := builder.CreateString(header.Title)
			defer flat.HeaderAddTitle(builder, offset)
		}

		if header.Description != "" {
			offset := builder.CreateString(header.Description)
			defer flat.HeaderAddDescription(builder, offset)
		}

		if header.Metadata != "" {
			offset := builder.CreateString(header.Metadata)
			defer flat.HeaderAddMetadata(builder, offset)
		}

		flat.HeaderStart(builder)

		return nil
	})
	if err != nil {
		return 0, err
	}

	return flat.HeaderEnd(builder), nil
}
