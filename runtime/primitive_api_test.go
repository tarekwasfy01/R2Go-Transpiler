package runtime

import "testing"

func TestEveryMatrixPrimitiveIsPubliclyAddressable(t *testing.T) {
	for _, descriptor := range PrimitiveTable {
		if !PrimitiveKnown(descriptor.Name) {
			t.Errorf("matrix primitive %q is not addressable", descriptor.Name)
		}
	}
}
