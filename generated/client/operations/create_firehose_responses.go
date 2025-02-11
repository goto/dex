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

// CreateFirehoseReader is a Reader for the CreateFirehose structure.
type CreateFirehoseReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CreateFirehoseReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewCreateFirehoseCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewCreateFirehoseBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 409:
		result := NewCreateFirehoseConflict()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewCreateFirehoseInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewCreateFirehoseCreated creates a CreateFirehoseCreated with default headers values
func NewCreateFirehoseCreated() *CreateFirehoseCreated {
	return &CreateFirehoseCreated{}
}

/*
	CreateFirehoseCreated describes a response with status code 201, with default header values.

Successfully created
*/
type CreateFirehoseCreated struct {
	Payload *models.Firehose
}

func (o *CreateFirehoseCreated) Error() string {
	return fmt.Sprintf("[POST /dex/firehoses][%d] createFirehoseCreated  %+v", 201, o.Payload)
}
func (o *CreateFirehoseCreated) GetPayload() *models.Firehose {
	return o.Payload
}

func (o *CreateFirehoseCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Firehose)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateFirehoseBadRequest creates a CreateFirehoseBadRequest with default headers values
func NewCreateFirehoseBadRequest() *CreateFirehoseBadRequest {
	return &CreateFirehoseBadRequest{}
}

/*
	CreateFirehoseBadRequest describes a response with status code 400, with default header values.

Request was invalid.
*/
type CreateFirehoseBadRequest struct {
	Payload *models.ErrorResponse
}

func (o *CreateFirehoseBadRequest) Error() string {
	return fmt.Sprintf("[POST /dex/firehoses][%d] createFirehoseBadRequest  %+v", 400, o.Payload)
}
func (o *CreateFirehoseBadRequest) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *CreateFirehoseBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateFirehoseConflict creates a CreateFirehoseConflict with default headers values
func NewCreateFirehoseConflict() *CreateFirehoseConflict {
	return &CreateFirehoseConflict{}
}

/*
	CreateFirehoseConflict describes a response with status code 409, with default header values.

A firehose with same unique identifier already exists.
*/
type CreateFirehoseConflict struct {
	Payload *models.ErrorResponse
}

func (o *CreateFirehoseConflict) Error() string {
	return fmt.Sprintf("[POST /dex/firehoses][%d] createFirehoseConflict  %+v", 409, o.Payload)
}
func (o *CreateFirehoseConflict) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *CreateFirehoseConflict) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateFirehoseInternalServerError creates a CreateFirehoseInternalServerError with default headers values
func NewCreateFirehoseInternalServerError() *CreateFirehoseInternalServerError {
	return &CreateFirehoseInternalServerError{}
}

/*
	CreateFirehoseInternalServerError describes a response with status code 500, with default header values.

internal error
*/
type CreateFirehoseInternalServerError struct {
	Payload *models.ErrorResponse
}

func (o *CreateFirehoseInternalServerError) Error() string {
	return fmt.Sprintf("[POST /dex/firehoses][%d] createFirehoseInternalServerError  %+v", 500, o.Payload)
}
func (o *CreateFirehoseInternalServerError) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *CreateFirehoseInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
