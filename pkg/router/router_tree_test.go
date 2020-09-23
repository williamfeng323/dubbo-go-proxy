package router

import (
	"testing"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	"github.com/emirpasic/gods/trees/avltree"
	"github.com/stretchr/testify/assert"
)

func TestPut(t *testing.T) {
	rt := &RouterTree{
		tree:         avltree.NewWithStringComparator(),
		wildcardTree: avltree.NewWithStringComparator(),
	}
	n0 := getMockMethod(config.MethodGet)
	rt.Put("/", n0)
	_, ok := rt.tree.Get("/")
	assert.True(t, ok)

	err := rt.Put("/", n0)
	assert.Error(t, err, "Method GET already exists in path /")

	n1 := getMockMethod(config.MethodPost)
	err = rt.Put("/mock", n0)
	assert.Nil(t, err)
	err = rt.Put("/mock", n1)
	assert.Nil(t, err)
	mNode, ok := rt.tree.Get("/mock")
	assert.True(t, ok)
	assert.Equal(t, len(mNode.(*RouterNode).methods), 2)

	err = rt.Put("/mock/test", n0)
	assert.Nil(t, err)
	_, ok = rt.tree.Get("/mock")
	assert.True(t, ok)

	rt.Put("/test/:id", n0)
	tNode, ok := rt.tree.Get("/test/:id")
	assert.True(t, ok)
	assert.True(t, tNode.(*RouterNode).wildcard)

	err = rt.Put("/test/:id", n1)
	assert.Nil(t, err)
	err = rt.Put("/test/js", n0)
	assert.Error(t, err, "/test/:id wildcard already exist so that cannot add path /test/js")

	err = rt.Put("/test/:id/mock", n0)
	tNode, ok = rt.tree.Get("/test/:id/mock")
	assert.True(t, ok)
	assert.True(t, tNode.(*RouterNode).wildcard)
	assert.Nil(t, err)
}

func TestSearchWildcard(t *testing.T) {
	rt := &RouterTree{
		tree:         avltree.NewWithStringComparator(),
		wildcardTree: avltree.NewWithStringComparator(),
	}
	n0 := getMockMethod(config.MethodGet)
	e := rt.Put("/theboys", n0)
	assert.Nil(t, e)
	e = rt.Put("/theboys/:id", n0)
	assert.Nil(t, e)
	e = rt.Put("/vought/:id/supe/:name", n0)
	assert.Nil(t, e)

	_, ok := rt.searchWildcard("/marvel")
	assert.False(t, ok)
	_, ok = rt.searchWildcard("/theboys/:id/age")
	assert.False(t, ok)
	_, ok = rt.searchWildcard("/theboys/butcher")
	assert.True(t, ok)
	_, ok = rt.searchWildcard("/vought/:id/supe/homelander")
	assert.True(t, ok)
}

func TestContainParam(t *testing.T) {
	assert.True(t, containParam("/test/:id"))
	assert.False(t, containParam("/test"))
	assert.True(t, containParam("/test/:id/mock"))
}

func TestWildcardMatch(t *testing.T) {
	assert.True(t, wildcardMatch("/vought/:id", "/vought/12345"))
	assert.True(t, wildcardMatch("/vought/:id", "/vought/125abc"))
	assert.False(t, wildcardMatch("/vought/:id", "/vought/1234abcd/status"))
	assert.True(t, wildcardMatch("/voughT/:id/:action", "/Vought/1234abcd/attack"))
}
