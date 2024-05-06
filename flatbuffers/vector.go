package flatbuffers

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type VectorElementsFunc[T any] func(*T, int) bool

func VectorFromFlat[T any](numElems int, elemsFunc VectorElementsFunc[T]) []T {
	var elems []T
	for i := 0; i < numElems; i++ {
		var elem T
		if elemsFunc(&elem, i) {
			elems = append(elems, elem)
		}
	}

	return elems
}

type StartVectorFunc func(builder *flatbuffers.Builder, numElems int) (offset flatbuffers.UOffsetT)

func BuilderAddVector[T any](builder *flatbuffers.Builder, elems []T, starter StartVectorFunc, prepender func(T)) flatbuffers.UOffsetT {
	starter(builder, len(elems))

	for i := len(elems) - 1; i >= 0; i-- {
		prepender(elems[i])
	}

	return builder.EndVector(len(elems))
}
