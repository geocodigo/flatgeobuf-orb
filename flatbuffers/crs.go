package flatbuffers

import (
	"github.com/gogama/flatgeobuf/flatgeobuf/flat"
	flatbuffers "github.com/google/flatbuffers/go"
)

type CRS struct {
	Org         string
	Code        int32
	Name        string
	Description string
	WKT         string
	CodeString  string
}

func CRSFromFlat(flatCRS *flat.Crs) (*CRS, error) {
	var crs CRS
	err := SafeInteraction(func() error {
		crs.Org = string(flatCRS.Org())
		crs.Code = flatCRS.Code()
		crs.Name = string(flatCRS.Name())
		crs.Description = string(flatCRS.Description())
		crs.WKT = string(flatCRS.Wkt())
		crs.CodeString = string(flatCRS.CodeString())

		return nil
	})

	return &crs, err
}

func BuilderAddCRS(builder *flatbuffers.Builder, crs CRS) (flatbuffers.UOffsetT, error) {
	err := SafeInteraction(func() error {
		if crs.Org != "" {
			offset := builder.CreateString(crs.Org)
			defer flat.CrsAddOrg(builder, offset)
		}

		defer flat.CrsAddCode(builder, crs.Code)

		if crs.Name != "" {
			offset := builder.CreateString(crs.Name)
			defer flat.CrsAddName(builder, offset)
		}

		if crs.Description != "" {
			offset := builder.CreateString(crs.Description)
			defer flat.CrsAddDescription(builder, offset)
		}

		if crs.WKT != "" {
			offset := builder.CreateString(crs.WKT)
			defer flat.CrsAddWkt(builder, offset)
		}

		if crs.CodeString != "" {
			offset := builder.CreateString(crs.CodeString)
			defer flat.CrsAddCodeString(builder, offset)
		}

		flat.CrsStart(builder)

		return nil
	})
	if err != nil {
		return 0, err
	}

	return flat.CrsEnd(builder), nil
}
