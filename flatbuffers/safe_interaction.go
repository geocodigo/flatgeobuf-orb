package flatbuffers

import (
	"fmt"
)

func SafeInteraction(f func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: flatbuffers: %v", r)
		}
	}()

	err = f()

	return
}
