package client_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/jacknotes/restful-api-demo/apps/host"
	"github.com/jacknotes/restful-api-demo/client"
	"github.com/stretchr/testify/assert"
)

func TestHost(t *testing.T) {
	should := assert.New(t)

	c, err := client.NewClient(client.NewDefaultConfig())
	should.NoError(err)

	set, err := c.Host().QueryHost(context.Background(), host.NewQueryHostRequest())
	should.NoError(err)

	fmt.Println(set)
}
