package flatbuffers

import (
	"github.com/gogama/flatgeobuf/flatgeobuf/flat"
	flatbuffers "github.com/google/flatbuffers/go"
)

type Geometry struct {
	Ends         []uint32
	XY           []float64
	Z            []float64
	M            []float64
	T            []float64
	TM           []uint64
	GeometryType flat.GeometryType
	Parts        []Geometry
}

func GeometryFromFlat(flatGeometry *flat.Geometry) (*Geometry, error) {
	var geometry Geometry
	err := SafeInteraction(func() error {
		for i := 0; i < flatGeometry.EndsLength(); i++ {
			geometry.Ends = append(geometry.Ends, flatGeometry.Ends(i))
		}

		for i := 0; i < flatGeometry.XyLength(); i++ {
			geometry.XY = append(geometry.XY, flatGeometry.Xy(i))
		}

		for i := 0; i < flatGeometry.ZLength(); i++ {
			geometry.Z = append(geometry.Z, flatGeometry.Z(i))
		}

		for i := 0; i < flatGeometry.MLength(); i++ {
			geometry.M = append(geometry.M, flatGeometry.M(i))
		}

		for i := 0; i < flatGeometry.TLength(); i++ {
			geometry.T = append(geometry.T, flatGeometry.T(i))
		}

		for i := 0; i < flatGeometry.TmLength(); i++ {
			geometry.TM = append(geometry.TM, flatGeometry.Tm(i))
		}

		geometry.GeometryType = flatGeometry.Type()

		if parts, err := VectorFromFlat(flatGeometry.PartsLength(), flatGeometry.Parts, GeometryFromFlat); err != nil {
			return err
		} else {
			geometry.Parts = parts
		}

		return nil
	})

	return &geometry, err
}

func BuilderAddGeometry(builder *flatbuffers.Builder, geometry Geometry) (flatbuffers.UOffsetT, error) {
	err := SafeInteraction(func() error {
		if geometry.Ends != nil {
			offset := BuilderAddVector(builder, geometry.Ends, flat.GeometryStartEndsVector, builder.PrependUint32)
			defer flat.GeometryAddEnds(builder, offset)
		}

		if geometry.XY != nil {
			offset := BuilderAddVector(builder, geometry.XY, flat.GeometryStartXyVector, builder.PrependFloat64)
			defer flat.GeometryAddXy(builder, offset)
		}

		if geometry.Z != nil {
			offset := BuilderAddVector(builder, geometry.Z, flat.GeometryStartZVector, builder.PrependFloat64)
			defer flat.GeometryAddZ(builder, offset)
		}

		if geometry.M != nil {
			offset := BuilderAddVector(builder, geometry.M, flat.GeometryStartMVector, builder.PrependFloat64)
			defer flat.GeometryAddM(builder, offset)
		}

		if geometry.T != nil {
			offset := BuilderAddVector(builder, geometry.T, flat.GeometryStartTVector, builder.PrependFloat64)
			defer flat.GeometryAddT(builder, offset)
		}

		if geometry.TM != nil {
			offset := BuilderAddVector(builder, geometry.TM, flat.GeometryStartTmVector, builder.PrependUint64)
			defer flat.GeometryAddTm(builder, offset)
		}

		defer flat.GeometryAddType(builder, geometry.GeometryType)

		if geometry.Parts != nil {
			offset, err := BuilderAddGeometries(builder, geometry.Parts, flat.GeometryStartPartsVector)
			if err != nil {
				return err
			}
			defer flat.GeometryAddParts(builder, offset)
		}

		flat.GeometryStart(builder)

		return nil
	})
	if err != nil {
		return 0, err
	}

	return flat.GeometryEnd(builder), nil
}

func BuilderAddGeometries(builder *flatbuffers.Builder, geometries []Geometry, startGeometriesFunc StartVectorFunc) (flatbuffers.UOffsetT, error) {
	offsets := make([]flatbuffers.UOffsetT, len(geometries))
	for i, geometry := range geometries {
		offset, err := BuilderAddGeometry(builder, geometry)
		if err != nil {
			return 0, err
		}
		offsets[i] = offset
	}

	return BuilderAddVector(builder, offsets, startGeometriesFunc, builder.PrependUOffsetT), nil
}
