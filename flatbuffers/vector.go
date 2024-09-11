package flatbuffers

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type FromFlatFunc[T any, FlatT any] func(*FlatT) (*T, error)

type FlatGetElementAtFunc[T any] func(*T, int) bool

func VectorFromFlat[T any, FlatT any](numElems int, flatElemFunc FlatGetElementAtFunc[FlatT], fromFlatFunc FromFlatFunc[T, FlatT]) ([]T, error) {
	var elems []T

	for i := 0; i < numElems; i++ {
		var flatElem FlatT
		if flatElemFunc(&flatElem, i) {
			elem, err := fromFlatFunc(&flatElem)
			if err != nil {
				return nil, err
			}
			elems = append(elems, *elem)
		}
	}

	return elems, nil
}

type StartVectorFunc func(builder *flatbuffers.Builder, numElems int) (offset flatbuffers.UOffsetT)

func BuilderAddVector[T any](builder *flatbuffers.Builder, elems []T, starter StartVectorFunc, prepender func(T)) flatbuffers.UOffsetT {
	starter(builder, len(elems))

	for i := len(elems) - 1; i >= 0; i-- {
		prepender(elems[i])
	}

	return builder.EndVector(len(elems))
}
