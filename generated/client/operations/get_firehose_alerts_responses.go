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

// GetFirehoseAlertsReader is a Reader for the GetFirehoseAlerts structure.
type GetFirehoseAlertsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetFirehoseAlertsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetFirehoseAlertsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewGetFirehoseAlertsNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetFirehoseAlertsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetFirehoseAlertsOK creates a GetFirehoseAlertsOK with default headers values
func NewGetFirehoseAlertsOK() *GetFirehoseAlertsOK {
	return &GetFirehoseAlertsOK{}
}

/*
	GetFirehoseAlertsOK describes a response with status code 200, with default header values.

alerts for given firehose URN.
*/
type GetFirehoseAlertsOK struct {
	Payload *models.AlertArray
}

func (o *GetFirehoseAlertsOK) Error() string {
	return fmt.Sprintf("[GET /dex/firehoses/{firehoseUrn}/alerts][%d] getFirehoseAlertsOK  %+v", 200, o.Payload)
}
func (o *GetFirehoseAlertsOK) GetPayload() *models.AlertArray {
	return o.Payload
}

func (o *GetFirehoseAlertsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.AlertArray)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetFirehoseAlertsNotFound creates a GetFirehoseAlertsNotFound with default headers values
func NewGetFirehoseAlertsNotFound() *GetFirehoseAlertsNotFound {
	return &GetFirehoseAlertsNotFound{}
}

/*
	GetFirehoseAlertsNotFound describes a response with status code 404, with default header values.

Firehose with given URN was not found
*/
type GetFirehoseAlertsNotFound struct {
	Payload *models.ErrorResponse
}

func (o *GetFirehoseAlertsNotFound) Error() string {
	return fmt.Sprintf("[GET /dex/firehoses/{firehoseUrn}/alerts][%d] getFirehoseAlertsNotFound  %+v", 404, o.Payload)
}
func (o *GetFirehoseAlertsNotFound) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *GetFirehoseAlertsNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetFirehoseAlertsInternalServerError creates a GetFirehoseAlertsInternalServerError with default headers values
func NewGetFirehoseAlertsInternalServerError() *GetFirehoseAlertsInternalServerError {
	return &GetFirehoseAlertsInternalServerError{}
}

/*
	GetFirehoseAlertsInternalServerError describes a response with status code 500, with default header values.

internal error
*/
type GetFirehoseAlertsInternalServerError struct {
	Payload *models.ErrorResponse
}

func (o *GetFirehoseAlertsInternalServerError) Error() string {
	return fmt.Sprintf("[GET /dex/firehoses/{firehoseUrn}/alerts][%d] getFirehoseAlertsInternalServerError  %+v", 500, o.Payload)
}
func (o *GetFirehoseAlertsInternalServerError) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *GetFirehoseAlertsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
