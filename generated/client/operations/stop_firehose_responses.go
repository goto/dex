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

// StopFirehoseReader is a Reader for the StopFirehose structure.
type StopFirehoseReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *StopFirehoseReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewStopFirehoseOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewStopFirehoseBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewStopFirehoseNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewStopFirehoseInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewStopFirehoseOK creates a StopFirehoseOK with default headers values
func NewStopFirehoseOK() *StopFirehoseOK {
	return &StopFirehoseOK{}
}

/*
StopFirehoseOK describes a response with status code 200, with default header values.

Successfully applied update.
*/
type StopFirehoseOK struct {
	Payload *models.Firehose
}

// IsSuccess returns true when this stop firehose o k response has a 2xx status code
func (o *StopFirehoseOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this stop firehose o k response has a 3xx status code
func (o *StopFirehoseOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this stop firehose o k response has a 4xx status code
func (o *StopFirehoseOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this stop firehose o k response has a 5xx status code
func (o *StopFirehoseOK) IsServerError() bool {
	return false
}

// IsCode returns true when this stop firehose o k response a status code equal to that given
func (o *StopFirehoseOK) IsCode(code int) bool {
	return code == 200
}

func (o *StopFirehoseOK) Error() string {
	return fmt.Sprintf("[POST /api/firehoses/{firehoseUrn}/stop][%d] stopFirehoseOK  %+v", 200, o.Payload)
}

func (o *StopFirehoseOK) String() string {
	return fmt.Sprintf("[POST /api/firehoses/{firehoseUrn}/stop][%d] stopFirehoseOK  %+v", 200, o.Payload)
}

func (o *StopFirehoseOK) GetPayload() *models.Firehose {
	return o.Payload
}

func (o *StopFirehoseOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Firehose)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewStopFirehoseBadRequest creates a StopFirehoseBadRequest with default headers values
func NewStopFirehoseBadRequest() *StopFirehoseBadRequest {
	return &StopFirehoseBadRequest{}
}

/*
StopFirehoseBadRequest describes a response with status code 400, with default header values.

Update request is not valid.
*/
type StopFirehoseBadRequest struct {
	Payload *models.ErrorResponse
}

// IsSuccess returns true when this stop firehose bad request response has a 2xx status code
func (o *StopFirehoseBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this stop firehose bad request response has a 3xx status code
func (o *StopFirehoseBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this stop firehose bad request response has a 4xx status code
func (o *StopFirehoseBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this stop firehose bad request response has a 5xx status code
func (o *StopFirehoseBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this stop firehose bad request response a status code equal to that given
func (o *StopFirehoseBadRequest) IsCode(code int) bool {
	return code == 400
}

func (o *StopFirehoseBadRequest) Error() string {
	return fmt.Sprintf("[POST /api/firehoses/{firehoseUrn}/stop][%d] stopFirehoseBadRequest  %+v", 400, o.Payload)
}

func (o *StopFirehoseBadRequest) String() string {
	return fmt.Sprintf("[POST /api/firehoses/{firehoseUrn}/stop][%d] stopFirehoseBadRequest  %+v", 400, o.Payload)
}

func (o *StopFirehoseBadRequest) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *StopFirehoseBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewStopFirehoseNotFound creates a StopFirehoseNotFound with default headers values
func NewStopFirehoseNotFound() *StopFirehoseNotFound {
	return &StopFirehoseNotFound{}
}

/*
StopFirehoseNotFound describes a response with status code 404, with default header values.

Firehose with given URN was not found
*/
type StopFirehoseNotFound struct {
	Payload *models.ErrorResponse
}

// IsSuccess returns true when this stop firehose not found response has a 2xx status code
func (o *StopFirehoseNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this stop firehose not found response has a 3xx status code
func (o *StopFirehoseNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this stop firehose not found response has a 4xx status code
func (o *StopFirehoseNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this stop firehose not found response has a 5xx status code
func (o *StopFirehoseNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this stop firehose not found response a status code equal to that given
func (o *StopFirehoseNotFound) IsCode(code int) bool {
	return code == 404
}

func (o *StopFirehoseNotFound) Error() string {
	return fmt.Sprintf("[POST /api/firehoses/{firehoseUrn}/stop][%d] stopFirehoseNotFound  %+v", 404, o.Payload)
}

func (o *StopFirehoseNotFound) String() string {
	return fmt.Sprintf("[POST /api/firehoses/{firehoseUrn}/stop][%d] stopFirehoseNotFound  %+v", 404, o.Payload)
}

func (o *StopFirehoseNotFound) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *StopFirehoseNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewStopFirehoseInternalServerError creates a StopFirehoseInternalServerError with default headers values
func NewStopFirehoseInternalServerError() *StopFirehoseInternalServerError {
	return &StopFirehoseInternalServerError{}
}

/*
StopFirehoseInternalServerError describes a response with status code 500, with default header values.

internal error
*/
type StopFirehoseInternalServerError struct {
	Payload *models.ErrorResponse
}

// IsSuccess returns true when this stop firehose internal server error response has a 2xx status code
func (o *StopFirehoseInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this stop firehose internal server error response has a 3xx status code
func (o *StopFirehoseInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this stop firehose internal server error response has a 4xx status code
func (o *StopFirehoseInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this stop firehose internal server error response has a 5xx status code
func (o *StopFirehoseInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this stop firehose internal server error response a status code equal to that given
func (o *StopFirehoseInternalServerError) IsCode(code int) bool {
	return code == 500
}

func (o *StopFirehoseInternalServerError) Error() string {
	return fmt.Sprintf("[POST /api/firehoses/{firehoseUrn}/stop][%d] stopFirehoseInternalServerError  %+v", 500, o.Payload)
}

func (o *StopFirehoseInternalServerError) String() string {
	return fmt.Sprintf("[POST /api/firehoses/{firehoseUrn}/stop][%d] stopFirehoseInternalServerError  %+v", 500, o.Payload)
}

func (o *StopFirehoseInternalServerError) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *StopFirehoseInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
