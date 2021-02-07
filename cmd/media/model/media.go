package model

import (
	"errors"
	"fmt"
	"strings"

	"github.com/reecerussell/open-social/cmd/media/dao"
)

var validContentTypes = [...]string{
	"image/jpeg",
	"image/png",
}

// Media is a domain model for the media domain.
type Media struct {
	id          int
	referenceID string
	contentType string
}

// NewMedia returns a new instance of the Media domain model.
func NewMedia(contentType string) (*Media, error) {
	m := &Media{}

	err := m.setContentType(contentType)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// ID returns the media's id.
func (m *Media) ID() int {
	return m.id
}

// ReferenceID returns the media's reference id.
func (m *Media) ReferenceID() string {
	return m.referenceID
}

func (m *Media) setContentType(contentType string) error {
	if contentType == "" {
		return errors.New("contentType is a required field")
	}

	contentType = strings.ToLower(contentType)

	for _, ct := range validContentTypes {
		if ct == contentType {
			m.contentType = contentType
			return nil
		}
	}

	return fmt.Errorf("the content type '%s' is not allowed", contentType)
}

// SetID sets the id of the Media.
func (m *Media) SetID(id int) {
	m.id = id
}

// SetReferenceID sets the reference id of the Media.
func (m *Media) SetReferenceID(referenceID string) {
	m.referenceID = referenceID
}

// Dao returns a data access object for the media instance.
func (m *Media) Dao() *dao.Media {
	return &dao.Media{
		ID:          m.id,
		ReferenceID: m.referenceID,
		ContentType: m.contentType,
	}
}
