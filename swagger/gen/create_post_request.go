// Code generated by go-swagger; DO NOT EDIT.

package gen

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// CreatePostRequest CreatePostRequest
//
// swagger:model CreatePostRequest
type CreatePostRequest struct {

	// content
	// Max Length: 140
	// Min Length: 1
	Content string `json:"content,omitempty"`
}

// Validate validates this create post request
func (m *CreatePostRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateContent(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CreatePostRequest) validateContent(formats strfmt.Registry) error {
	if swag.IsZero(m.Content) { // not required
		return nil
	}

	if err := validate.MinLength("content", "body", m.Content, 1); err != nil {
		return err
	}

	if err := validate.MaxLength("content", "body", m.Content, 140); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this create post request based on context it is used
func (m *CreatePostRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *CreatePostRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CreatePostRequest) UnmarshalBinary(b []byte) error {
	var res CreatePostRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}