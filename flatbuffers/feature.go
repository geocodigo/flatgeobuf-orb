package flatbuffers

import (
	"github.com/gogama/flatgeobuf/flatgeobuf/flat"
	flatbuffers "github.com/google/flatbuffers/go"
)

type Feature struct {
	Geometry   *Geometry
	Properties []byte
	Columns    []Column
}

func FeatureFromFlat(flatFeature *flat.Feature) (*Feature, error) {
	var feature Feature

	err := SafeInteraction(func() error {
		var flatGeometry flat.Geometry
		if flatFeature.Geometry(&flatGeometry) != nil {
			geometry, err := GeometryFromFlat(&flatGeometry)
			if err != nil {
				return err
			}
			feature.Geometry = geometry
		}

		feature.Properties = flatFeature.PropertiesBytes()

		for _, col := range VectorFromFlat(flatFeature.ColumnsLength(), flatFeature.Columns) {
			column, err := ColumnFromFlat(&col)
			if err != nil {
				return err
			}
			feature.Columns = append(feature.Columns, *column)
		}

		return nil
	})

	return &feature, err
}

func BuilderAddFeature(builder *flatbuffers.Builder, feature Feature) (flatbuffers.UOffsetT, error) {
	err := SafeInteraction(func() error {
		if feature.Geometry != nil {
			geometryOffset, err := BuilderAddGeometry(builder, *feature.Geometry)
			if err != nil {
				return err
			}
			defer flat.FeatureAddGeometry(builder, geometryOffset)
		}

		if feature.Properties != nil {
			defer flat.FeatureAddProperties(builder, builder.CreateByteVector(feature.Properties))
		}

		if feature.Columns != nil {
			columnsOffset, err := BuilderAddColumns(builder, feature.Columns, flat.FeatureStartColumnsVector)
			if err != nil {
				return err
			}
			defer flat.FeatureAddColumns(builder, columnsOffset)
		}

		flat.FeatureStart(builder)

		return nil
	})
	if err != nil {
		return 0, err
	}

	return flat.FeatureEnd(builder), err
}
