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

// GetFirehoseLogsReader is a Reader for the GetFirehoseLogs structure.
type GetFirehoseLogsReader struct {
	formats strfmt.Registry
	writer  io.Writer
}

// ReadResponse reads a server response into the received o.
func (o *GetFirehoseLogsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetFirehoseLogsOK(o.writer)
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewGetFirehoseLogsBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewGetFirehoseLogsNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetFirehoseLogsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetFirehoseLogsOK creates a GetFirehoseLogsOK with default headers values
func NewGetFirehoseLogsOK(writer io.Writer) *GetFirehoseLogsOK {
	return &GetFirehoseLogsOK{

		Payload: writer,
	}
}

/*
GetFirehoseLogsOK describes a response with status code 200, with default header values.

Found logs for given firehose URN.
*/
type GetFirehoseLogsOK struct {
	Payload io.Writer
}

// IsSuccess returns true when this get firehose logs o k response has a 2xx status code
func (o *GetFirehoseLogsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get firehose logs o k response has a 3xx status code
func (o *GetFirehoseLogsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get firehose logs o k response has a 4xx status code
func (o *GetFirehoseLogsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get firehose logs o k response has a 5xx status code
func (o *GetFirehoseLogsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get firehose logs o k response a status code equal to that given
func (o *GetFirehoseLogsOK) IsCode(code int) bool {
	return code == 200
}

func (o *GetFirehoseLogsOK) Error() string {
	return fmt.Sprintf("[GET /dex/firehoses/{firehoseUrn}/logs][%d] getFirehoseLogsOK  %+v", 200, o.Payload)
}

func (o *GetFirehoseLogsOK) String() string {
	return fmt.Sprintf("[GET /dex/firehoses/{firehoseUrn}/logs][%d] getFirehoseLogsOK  %+v", 200, o.Payload)
}

func (o *GetFirehoseLogsOK) GetPayload() io.Writer {
	return o.Payload
}

func (o *GetFirehoseLogsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetFirehoseLogsBadRequest creates a GetFirehoseLogsBadRequest with default headers values
func NewGetFirehoseLogsBadRequest() *GetFirehoseLogsBadRequest {
	return &GetFirehoseLogsBadRequest{}
}

/*
GetFirehoseLogsBadRequest describes a response with status code 400, with default header values.

Get logs request is not valid.
*/
type GetFirehoseLogsBadRequest struct {
	Payload *models.ErrorResponse
}

// IsSuccess returns true when this get firehose logs bad request response has a 2xx status code
func (o *GetFirehoseLogsBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get firehose logs bad request response has a 3xx status code
func (o *GetFirehoseLogsBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get firehose logs bad request response has a 4xx status code
func (o *GetFirehoseLogsBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this get firehose logs bad request response has a 5xx status code
func (o *GetFirehoseLogsBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this get firehose logs bad request response a status code equal to that given
func (o *GetFirehoseLogsBadRequest) IsCode(code int) bool {
	return code == 400
}

func (o *GetFirehoseLogsBadRequest) Error() string {
	return fmt.Sprintf("[GET /dex/firehoses/{firehoseUrn}/logs][%d] getFirehoseLogsBadRequest  %+v", 400, o.Payload)
}

func (o *GetFirehoseLogsBadRequest) String() string {
	return fmt.Sprintf("[GET /dex/firehoses/{firehoseUrn}/logs][%d] getFirehoseLogsBadRequest  %+v", 400, o.Payload)
}

func (o *GetFirehoseLogsBadRequest) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *GetFirehoseLogsBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetFirehoseLogsNotFound creates a GetFirehoseLogsNotFound with default headers values
func NewGetFirehoseLogsNotFound() *GetFirehoseLogsNotFound {
	return &GetFirehoseLogsNotFound{}
}

/*
GetFirehoseLogsNotFound describes a response with status code 404, with default header values.

Firehose with given URN was not found
*/
type GetFirehoseLogsNotFound struct {
	Payload *models.ErrorResponse
}

// IsSuccess returns true when this get firehose logs not found response has a 2xx status code
func (o *GetFirehoseLogsNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get firehose logs not found response has a 3xx status code
func (o *GetFirehoseLogsNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get firehose logs not found response has a 4xx status code
func (o *GetFirehoseLogsNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this get firehose logs not found response has a 5xx status code
func (o *GetFirehoseLogsNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this get firehose logs not found response a status code equal to that given
func (o *GetFirehoseLogsNotFound) IsCode(code int) bool {
	return code == 404
}

func (o *GetFirehoseLogsNotFound) Error() string {
	return fmt.Sprintf("[GET /dex/firehoses/{firehoseUrn}/logs][%d] getFirehoseLogsNotFound  %+v", 404, o.Payload)
}

func (o *GetFirehoseLogsNotFound) String() string {
	return fmt.Sprintf("[GET /dex/firehoses/{firehoseUrn}/logs][%d] getFirehoseLogsNotFound  %+v", 404, o.Payload)
}

func (o *GetFirehoseLogsNotFound) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *GetFirehoseLogsNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetFirehoseLogsInternalServerError creates a GetFirehoseLogsInternalServerError with default headers values
func NewGetFirehoseLogsInternalServerError() *GetFirehoseLogsInternalServerError {
	return &GetFirehoseLogsInternalServerError{}
}

/*
GetFirehoseLogsInternalServerError describes a response with status code 500, with default header values.

internal error
*/
type GetFirehoseLogsInternalServerError struct {
	Payload *models.ErrorResponse
}

// IsSuccess returns true when this get firehose logs internal server error response has a 2xx status code
func (o *GetFirehoseLogsInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get firehose logs internal server error response has a 3xx status code
func (o *GetFirehoseLogsInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get firehose logs internal server error response has a 4xx status code
func (o *GetFirehoseLogsInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this get firehose logs internal server error response has a 5xx status code
func (o *GetFirehoseLogsInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this get firehose logs internal server error response a status code equal to that given
func (o *GetFirehoseLogsInternalServerError) IsCode(code int) bool {
	return code == 500
}

func (o *GetFirehoseLogsInternalServerError) Error() string {
	return fmt.Sprintf("[GET /dex/firehoses/{firehoseUrn}/logs][%d] getFirehoseLogsInternalServerError  %+v", 500, o.Payload)
}

func (o *GetFirehoseLogsInternalServerError) String() string {
	return fmt.Sprintf("[GET /dex/firehoses/{firehoseUrn}/logs][%d] getFirehoseLogsInternalServerError  %+v", 500, o.Payload)
}

func (o *GetFirehoseLogsInternalServerError) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *GetFirehoseLogsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
