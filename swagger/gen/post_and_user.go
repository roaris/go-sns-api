// Code generated by go-swagger; DO NOT EDIT.

package gen

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// PostAndUser PostAndUser
//
// swagger:model PostAndUser
type PostAndUser struct {

	// post
	Post *Post `json:"post,omitempty"`

	// user
	User *User `json:"user,omitempty"`
}

// Validate validates this post and user
func (m *PostAndUser) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePost(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateUser(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PostAndUser) validatePost(formats strfmt.Registry) error {
	if swag.IsZero(m.Post) { // not required
		return nil
	}

	if m.Post != nil {
		if err := m.Post.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("post")
			}
			return err
		}
	}

	return nil
}

func (m *PostAndUser) validateUser(formats strfmt.Registry) error {
	if swag.IsZero(m.User) { // not required
		return nil
	}

	if m.User != nil {
		if err := m.User.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("user")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this post and user based on the context it is used
func (m *PostAndUser) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidatePost(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateUser(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PostAndUser) contextValidatePost(ctx context.Context, formats strfmt.Registry) error {

	if m.Post != nil {
		if err := m.Post.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("post")
			}
			return err
		}
	}

	return nil
}

func (m *PostAndUser) contextValidateUser(ctx context.Context, formats strfmt.Registry) error {

	if m.User != nil {
		if err := m.User.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("user")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PostAndUser) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PostAndUser) UnmarshalBinary(b []byte) error {
	var res PostAndUser
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
