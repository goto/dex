// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/goto/dex/generated/models"
)

// GetAlertPolicyReader is a Reader for the GetAlertPolicy structure.
type GetAlertPolicyReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetAlertPolicyReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetAlertPolicyOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewGetAlertPolicyNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetAlertPolicyInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetAlertPolicyOK creates a GetAlertPolicyOK with default headers values
func NewGetAlertPolicyOK() *GetAlertPolicyOK {
	return &GetAlertPolicyOK{}
}

/*
	GetAlertPolicyOK describes a response with status code 200, with default header values.

Found alert policy for given URN.
*/
type GetAlertPolicyOK struct {
	Payload *models.AlertPolicy
}

func (o *GetAlertPolicyOK) Error() string {
	return fmt.Sprintf("[GET /dex/alerts/{projectSlug}/{resourceUrn}/policies][%d] getAlertPolicyOK  %+v", 200, o.Payload)
}
func (o *GetAlertPolicyOK) GetPayload() *models.AlertPolicy {
	return o.Payload
}

func (o *GetAlertPolicyOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.AlertPolicy)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetAlertPolicyNotFound creates a GetAlertPolicyNotFound with default headers values
func NewGetAlertPolicyNotFound() *GetAlertPolicyNotFound {
	return &GetAlertPolicyNotFound{}
}

/*
	GetAlertPolicyNotFound describes a response with status code 404, with default header values.

Could not find policies for given URN
*/
type GetAlertPolicyNotFound struct {
	Payload *models.ErrorResponse
}

func (o *GetAlertPolicyNotFound) Error() string {
	return fmt.Sprintf("[GET /dex/alerts/{projectSlug}/{resourceUrn}/policies][%d] getAlertPolicyNotFound  %+v", 404, o.Payload)
}
func (o *GetAlertPolicyNotFound) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *GetAlertPolicyNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetAlertPolicyInternalServerError creates a GetAlertPolicyInternalServerError with default headers values
func NewGetAlertPolicyInternalServerError() *GetAlertPolicyInternalServerError {
	return &GetAlertPolicyInternalServerError{}
}

/*
	GetAlertPolicyInternalServerError describes a response with status code 500, with default header values.

internal error
*/
type GetAlertPolicyInternalServerError struct {
	Payload *models.ErrorResponse
}

func (o *GetAlertPolicyInternalServerError) Error() string {
	return fmt.Sprintf("[GET /dex/alerts/{projectSlug}/{resourceUrn}/policies][%d] getAlertPolicyInternalServerError  %+v", 500, o.Payload)
}
func (o *GetAlertPolicyInternalServerError) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *GetAlertPolicyInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
