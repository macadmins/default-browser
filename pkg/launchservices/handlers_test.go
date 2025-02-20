package launchservices

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddLSHandlers(t *testing.T) {
	plist := Plist{}

	handlersToAdd := []LSHandler{
		{LSHandlerContentType: "public.text"},
		{LSHandlerURLScheme: "ftp"},
	}

	plist.AddLSHandlers(handlersToAdd)

	assert.Equal(t, handlersToAdd, plist.LSHandlers, "AddLSHandlers should add the provided handlers to the LSHandlers list")
}

func TestCleanHandlers(t *testing.T) {
	plist := Plist{
		LSHandlers: []LSHandler{
			{LSHandlerContentType: "public.url"},
			{LSHandlerContentType: "public.html"},
			{LSHandlerContentType: "public.xhtml"},
			{LSHandlerURLScheme: "https"},
			{LSHandlerURLScheme: "http"},
			{LSHandlerContentType: "public.text"},
			{LSHandlerURLScheme: "ftp"},
		},
	}

	plist.CleanHandlers()

	expectedHandlers := []LSHandler{
		{LSHandlerContentType: "public.text"},
		{LSHandlerURLScheme: "ftp"},
	}

	assert.Equal(t, expectedHandlers, plist.LSHandlers, "CleanHandlers should remove handlers with specific content types and URL schemes")
}
