package generator

import (
	"sort"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMediaWellKnownMime(t *testing.T) {
	w, ok := wellKnownMime(runtime.JSONMime)
	assert.True(t, ok)
	assert.Equal(t, jsonSerializer, w)

	w, ok = wellKnownMime(runtime.YAMLMime)
	assert.True(t, ok)
	assert.Equal(t, "yaml", w)

	w, ok = wellKnownMime(runtime.JSONMime + "+version=1;param=1")
	assert.True(t, ok)
	assert.Equal(t, jsonSerializer, w)

	w, ok = wellKnownMime("unknown")
	assert.False(t, ok)
	assert.Equal(t, "", w)
}

func TestMediaMime(t *testing.T) {
	params := "param=1;param=2"
	withParams := runtime.JSONMime + ";" + params
	assert.Equal(t, runtime.JSONMime, mediaMime(runtime.JSONMime))
	assert.Equal(t, runtime.JSONMime, mediaMime(withParams))

	assert.Equal(t, params, mediaParameters(withParams))
	assert.Equal(t, "", mediaParameters(runtime.JSONMime))
}

func TestMediaMakeSerializers(t *testing.T) {
	app := appGenerator{
		Name:     "myapp",
		Receiver: "myReceiver",
	}

	res, supportsJSON := app.makeSerializers([]string{
		runtime.JSONMime,
		"application/json;param=1",
		"application/json;param=2",
		"application/json",
		"application/json+subtype;param=6",
		"application/yaml;param=1",
		runtime.YAMLMime,
		runtime.YAMLMime + "; param=xy",
		runtime.YAMLMime + ";  param=xy", // duplicate
		"application/funny;param=x",
	}, func(media string) (string, bool) {
		w, ok := knownConsumers[media]
		if !ok {
			w = "custom.FunnyConsume()"
		}
		return w, true
	})
	assert.True(t, supportsJSON)
	assert.True(t, sort.IsSorted(res))
	assert.Len(t, res, 3)

	for _, ser := range res {
		assert.NotEmpty(t, ser.AppName)
		assert.NotEmpty(t, ser.ReceiverName)
		assert.NotEmpty(t, ser.Implementation)

		switch ser.Name {
		case jsonSerializer:
			assert.Len(t, ser.AllSerializers, 2)
			for _, media := range ser.AllSerializers {
				assert.Equal(t, ser.AppName, media.AppName)
				assert.Equal(t, ser.ReceiverName, media.ReceiverName)
				assert.Equal(t, ser.Implementation, media.Implementation)
				switch media.MediaType {
				case "application/json":
					assert.Len(t, media.Parameters, 2)
				case "application/json+subtype":
					assert.Len(t, media.Parameters, 1)
				default:
					t.Logf("unexpected media type: %s in %v", media.MediaType, ser.AllSerializers)
					t.Fail()
				}
			}

		case "yaml":
			assert.Len(t, ser.AllSerializers, 2)
			for _, media := range ser.AllSerializers {
				assert.Equal(t, ser.AppName, media.AppName)
				assert.Equal(t, ser.ReceiverName, media.ReceiverName)
				assert.Equal(t, ser.Implementation, media.Implementation)
				switch media.MediaType {
				case runtime.YAMLMime:
					assert.Len(t, media.Parameters, 1)
				case "application/yaml":
					assert.Len(t, media.Parameters, 1)
				default:
					t.Logf("unexpected media type: %s in %v", media.MediaType, ser.AllSerializers)
					t.Fail()
				}
			}

		case "applicationFunny":
			assert.Len(t, ser.AllSerializers, 1)
			for _, media := range ser.AllSerializers {
				assert.Equal(t, ser.AppName, media.AppName)
				assert.Equal(t, ser.ReceiverName, media.ReceiverName)
				assert.Equal(t, ser.Implementation, media.Implementation)
				switch media.MediaType {
				case "application/funny":
					assert.Len(t, media.Parameters, 1)
				default:
					t.Logf("unexpected media type: %s in %v", media.MediaType, ser.AllSerializers)
					t.Fail()
				}
			}
		default:
			t.Logf("unexpected serializer name: %s", ser.Name)
			t.Fail()
		}
	}

	// no json, one non default serializer
	res, supportsJSON = app.makeSerializers([]string{
		"application/yaml",
		runtime.TextMime,
		"application/funny",
	}, func(media string) (string, bool) {
		w, ok := knownConsumers[media]
		return w, ok
	})
	assert.False(t, supportsJSON)
	assert.True(t, sort.IsSorted(res))
	assert.Len(t, res, 3)
	for _, ser := range res {
		assert.NotEmpty(t, ser.AppName)
		assert.NotEmpty(t, ser.ReceiverName)
		switch ser.Name {
		case "yaml":
			assert.NotEmpty(t, ser.Implementation)
			assert.Len(t, ser.AllSerializers, 1)
		case "txt":
			assert.NotEmpty(t, ser.Implementation)
			assert.Len(t, ser.AllSerializers, 1)
		case "applicationFunny":
			assert.Empty(t, ser.Implementation)
			assert.Len(t, ser.AllSerializers, 1)
		default:
			t.Logf("unexpected mime type: %s", ser.MediaType)
			t.Fail()
		}
	}

	// empty: defaults as json
	res, supportsJSON = app.makeSerializers([]string{}, func(_ string) (string, bool) { return "fake", true })
	assert.True(t, supportsJSON)
	assert.True(t, sort.IsSorted(res))
	require.Len(t, res, 1)
	assert.Equal(t, jsonSerializer, res[0].Name)
}
