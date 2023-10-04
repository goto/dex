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
		return nil, runtime.NewAPIError("[GET /dex/firehoses/{firehoseUrn}/alerts] getFirehoseAlerts", response, response.Code())
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

// IsSuccess returns true when this get firehose alerts o k response has a 2xx status code
func (o *GetFirehoseAlertsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get firehose alerts o k response has a 3xx status code
func (o *GetFirehoseAlertsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get firehose alerts o k response has a 4xx status code
func (o *GetFirehoseAlertsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get firehose alerts o k response has a 5xx status code
func (o *GetFirehoseAlertsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get firehose alerts o k response a status code equal to that given
func (o *GetFirehoseAlertsOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get firehose alerts o k response
func (o *GetFirehoseAlertsOK) Code() int {
	return 200
}

func (o *GetFirehoseAlertsOK) Error() string {
	return fmt.Sprintf("[GET /dex/firehoses/{firehoseUrn}/alerts][%d] getFirehoseAlertsOK  %+v", 200, o.Payload)
}

func (o *GetFirehoseAlertsOK) String() string {
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

// IsSuccess returns true when this get firehose alerts not found response has a 2xx status code
func (o *GetFirehoseAlertsNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get firehose alerts not found response has a 3xx status code
func (o *GetFirehoseAlertsNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get firehose alerts not found response has a 4xx status code
func (o *GetFirehoseAlertsNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this get firehose alerts not found response has a 5xx status code
func (o *GetFirehoseAlertsNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this get firehose alerts not found response a status code equal to that given
func (o *GetFirehoseAlertsNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the get firehose alerts not found response
func (o *GetFirehoseAlertsNotFound) Code() int {
	return 404
}

func (o *GetFirehoseAlertsNotFound) Error() string {
	return fmt.Sprintf("[GET /dex/firehoses/{firehoseUrn}/alerts][%d] getFirehoseAlertsNotFound  %+v", 404, o.Payload)
}

func (o *GetFirehoseAlertsNotFound) String() string {
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

// IsSuccess returns true when this get firehose alerts internal server error response has a 2xx status code
func (o *GetFirehoseAlertsInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get firehose alerts internal server error response has a 3xx status code
func (o *GetFirehoseAlertsInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get firehose alerts internal server error response has a 4xx status code
func (o *GetFirehoseAlertsInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this get firehose alerts internal server error response has a 5xx status code
func (o *GetFirehoseAlertsInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this get firehose alerts internal server error response a status code equal to that given
func (o *GetFirehoseAlertsInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the get firehose alerts internal server error response
func (o *GetFirehoseAlertsInternalServerError) Code() int {
	return 500
}

func (o *GetFirehoseAlertsInternalServerError) Error() string {
	return fmt.Sprintf("[GET /dex/firehoses/{firehoseUrn}/alerts][%d] getFirehoseAlertsInternalServerError  %+v", 500, o.Payload)
}

func (o *GetFirehoseAlertsInternalServerError) String() string {
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
