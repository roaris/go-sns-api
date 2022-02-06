// Code generated by go-swagger; DO NOT EDIT.

package gen

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// CreateFolloweeRequest CreateFolloweeRequest
//
// swagger:model CreateFolloweeRequest
type CreateFolloweeRequest struct {

	// followee id
	FolloweeID int64 `json:"followee_id,omitempty"`
}

// Validate validates this create followee request
func (m *CreateFolloweeRequest) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this create followee request based on context it is used
func (m *CreateFolloweeRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *CreateFolloweeRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CreateFolloweeRequest) UnmarshalBinary(b []byte) error {
	var res CreateFolloweeRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}