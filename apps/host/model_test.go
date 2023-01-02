package host_test

import (
	"testing"

	"github.com/jacknotes/restful-api-demo/apps/host"
	"github.com/stretchr/testify/assert"
)

func TestHostUpdate(t *testing.T) {
	should := assert.New(t)
	h := host.NewDefaultHost()

	patch := host.NewDefaultHost()
	patch.Name = "patch01"

	err := h.Patch(patch.Resource, patch.Describe)

	if should.NoError(err) {
		should.Equal(h.Name, patch.Name)
	}
}
