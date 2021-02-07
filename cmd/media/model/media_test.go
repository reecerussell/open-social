package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMedia(t *testing.T) {
	const testContentType = "image/png"

	media, err := NewMedia(testContentType)
	assert.NoError(t, err)
	assert.Equal(t, testContentType, media.contentType)
}

func TestNewMedia_GivenInvalidData_ReturnsError(t *testing.T) {
	t.Run("Invalid ContentType", func(t *testing.T) {
		media, err := NewMedia("")
		assert.Nil(t, media)
		assert.Equal(t, "contentType is a required field", err.Error())
	})
}

func TestMedia_SetContentType_ReturnsNoError(t *testing.T) {
	for _, contentType := range validContentTypes {
		var media Media
		err := media.setContentType(contentType)
		assert.NoError(t, err)
	}
}

func TestMedia_SetContentType_ReturnsError(t *testing.T) {
	t.Run("Given Empty Value", func(t *testing.T) {
		var media Media
		err := media.setContentType("")
		assert.Equal(t, "contentType is a required field", err.Error())
	})

	t.Run("Given Invalid Value", func(t *testing.T) {
		invalidValues := []string{
			"application/json",
			"text/plain",
		}

		for _, value := range invalidValues {
			var media Media
			err := media.setContentType(value)
			exp := fmt.Sprintf("the content type '%s' is not allowed", value)
			assert.Equal(t, exp, err.Error())
		}
	})
}

func TestMedia_SetID(t *testing.T) {
	const testID = 123

	var media Media
	media.SetID(testID)
	assert.Equal(t, testID, media.id)
}

func TestMedia_SetReferenceID(t *testing.T) {
	const testReferenceID = "1234"

	var media Media
	media.SetReferenceID(testReferenceID)
	assert.Equal(t, testReferenceID, media.referenceID)
}

func TestMedia_Dao(t *testing.T) {
	m := &Media{
		id:          123,
		referenceID: "2913",
		contentType: "image/jpeg",
	}

	d := m.Dao()

	assert.Equal(t, 123, d.ID)
	assert.Equal(t, "2913", d.ReferenceID)
	assert.Equal(t, "image/jpeg", d.ContentType)
}

func TestMedia_ID(t *testing.T) {
	const testID = 123

	m := &Media{id: testID}

	assert.Equal(t, testID, m.ID())
}

func TestMedia_ReferenceID(t *testing.T) {
	const testReferenceID = "123"

	m := &Media{referenceID: testReferenceID}

	assert.Equal(t, testReferenceID, m.ReferenceID())
}
