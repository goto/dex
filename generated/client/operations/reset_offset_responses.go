// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"github.com/goto/dex/generated/models"
)

// ResetOffsetReader is a Reader for the ResetOffset structure.
type ResetOffsetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ResetOffsetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewResetOffsetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewResetOffsetBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewResetOffsetNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewResetOffsetInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewResetOffsetOK creates a ResetOffsetOK with default headers values
func NewResetOffsetOK() *ResetOffsetOK {
	return &ResetOffsetOK{}
}

/*
ResetOffsetOK describes a response with status code 200, with default header values.

Found firehose with given URN
*/
type ResetOffsetOK struct {
	Payload *models.Firehose
}

// IsSuccess returns true when this reset offset o k response has a 2xx status code
func (o *ResetOffsetOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this reset offset o k response has a 3xx status code
func (o *ResetOffsetOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this reset offset o k response has a 4xx status code
func (o *ResetOffsetOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this reset offset o k response has a 5xx status code
func (o *ResetOffsetOK) IsServerError() bool {
	return false
}

// IsCode returns true when this reset offset o k response a status code equal to that given
func (o *ResetOffsetOK) IsCode(code int) bool {
	return code == 200
}

func (o *ResetOffsetOK) Error() string {
	return fmt.Sprintf("[POST /dex/firehoses/{firehoseUrn}/reset][%d] resetOffsetOK  %+v", 200, o.Payload)
}

func (o *ResetOffsetOK) String() string {
	return fmt.Sprintf("[POST /dex/firehoses/{firehoseUrn}/reset][%d] resetOffsetOK  %+v", 200, o.Payload)
}

func (o *ResetOffsetOK) GetPayload() *models.Firehose {
	return o.Payload
}

func (o *ResetOffsetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Firehose)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewResetOffsetBadRequest creates a ResetOffsetBadRequest with default headers values
func NewResetOffsetBadRequest() *ResetOffsetBadRequest {
	return &ResetOffsetBadRequest{}
}

/*
ResetOffsetBadRequest describes a response with status code 400, with default header values.

Update request is not valid.
*/
type ResetOffsetBadRequest struct {
	Payload *models.ErrorResponse
}

// IsSuccess returns true when this reset offset bad request response has a 2xx status code
func (o *ResetOffsetBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this reset offset bad request response has a 3xx status code
func (o *ResetOffsetBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this reset offset bad request response has a 4xx status code
func (o *ResetOffsetBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this reset offset bad request response has a 5xx status code
func (o *ResetOffsetBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this reset offset bad request response a status code equal to that given
func (o *ResetOffsetBadRequest) IsCode(code int) bool {
	return code == 400
}

func (o *ResetOffsetBadRequest) Error() string {
	return fmt.Sprintf("[POST /dex/firehoses/{firehoseUrn}/reset][%d] resetOffsetBadRequest  %+v", 400, o.Payload)
}

func (o *ResetOffsetBadRequest) String() string {
	return fmt.Sprintf("[POST /dex/firehoses/{firehoseUrn}/reset][%d] resetOffsetBadRequest  %+v", 400, o.Payload)
}

func (o *ResetOffsetBadRequest) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *ResetOffsetBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewResetOffsetNotFound creates a ResetOffsetNotFound with default headers values
func NewResetOffsetNotFound() *ResetOffsetNotFound {
	return &ResetOffsetNotFound{}
}

/*
ResetOffsetNotFound describes a response with status code 404, with default header values.

Firehose with given URN was not found
*/
type ResetOffsetNotFound struct {
	Payload *models.ErrorResponse
}

// IsSuccess returns true when this reset offset not found response has a 2xx status code
func (o *ResetOffsetNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this reset offset not found response has a 3xx status code
func (o *ResetOffsetNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this reset offset not found response has a 4xx status code
func (o *ResetOffsetNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this reset offset not found response has a 5xx status code
func (o *ResetOffsetNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this reset offset not found response a status code equal to that given
func (o *ResetOffsetNotFound) IsCode(code int) bool {
	return code == 404
}

func (o *ResetOffsetNotFound) Error() string {
	return fmt.Sprintf("[POST /dex/firehoses/{firehoseUrn}/reset][%d] resetOffsetNotFound  %+v", 404, o.Payload)
}

func (o *ResetOffsetNotFound) String() string {
	return fmt.Sprintf("[POST /dex/firehoses/{firehoseUrn}/reset][%d] resetOffsetNotFound  %+v", 404, o.Payload)
}

func (o *ResetOffsetNotFound) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *ResetOffsetNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewResetOffsetInternalServerError creates a ResetOffsetInternalServerError with default headers values
func NewResetOffsetInternalServerError() *ResetOffsetInternalServerError {
	return &ResetOffsetInternalServerError{}
}

/*
ResetOffsetInternalServerError describes a response with status code 500, with default header values.

internal error
*/
type ResetOffsetInternalServerError struct {
	Payload *models.ErrorResponse
}

// IsSuccess returns true when this reset offset internal server error response has a 2xx status code
func (o *ResetOffsetInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this reset offset internal server error response has a 3xx status code
func (o *ResetOffsetInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this reset offset internal server error response has a 4xx status code
func (o *ResetOffsetInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this reset offset internal server error response has a 5xx status code
func (o *ResetOffsetInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this reset offset internal server error response a status code equal to that given
func (o *ResetOffsetInternalServerError) IsCode(code int) bool {
	return code == 500
}

func (o *ResetOffsetInternalServerError) Error() string {
	return fmt.Sprintf("[POST /dex/firehoses/{firehoseUrn}/reset][%d] resetOffsetInternalServerError  %+v", 500, o.Payload)
}

func (o *ResetOffsetInternalServerError) String() string {
	return fmt.Sprintf("[POST /dex/firehoses/{firehoseUrn}/reset][%d] resetOffsetInternalServerError  %+v", 500, o.Payload)
}

func (o *ResetOffsetInternalServerError) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *ResetOffsetInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
ResetOffsetBody reset offset body
swagger:model ResetOffsetBody
*/
type ResetOffsetBody struct {

	// datetime
	// Example: 2022-10-10T10:10:10.100Z
	// Format: date-time
	Datetime strfmt.DateTime `json:"datetime,omitempty"`

	// to
	// Required: true
	// Enum: [DATETIME EARLIEST LATEST]
	To *string `json:"to"`
}

// Validate validates this reset offset body
func (o *ResetOffsetBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateDatetime(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateTo(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *ResetOffsetBody) validateDatetime(formats strfmt.Registry) error {
	if swag.IsZero(o.Datetime) { // not required
		return nil
	}

	if err := validate.FormatOf("body"+"."+"datetime", "body", "date-time", o.Datetime.String(), formats); err != nil {
		return err
	}

	return nil
}

var resetOffsetBodyTypeToPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["DATETIME","EARLIEST","LATEST"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		resetOffsetBodyTypeToPropEnum = append(resetOffsetBodyTypeToPropEnum, v)
	}
}

const (

	// ResetOffsetBodyToDATETIME captures enum value "DATETIME"
	ResetOffsetBodyToDATETIME string = "DATETIME"

	// ResetOffsetBodyToEARLIEST captures enum value "EARLIEST"
	ResetOffsetBodyToEARLIEST string = "EARLIEST"

	// ResetOffsetBodyToLATEST captures enum value "LATEST"
	ResetOffsetBodyToLATEST string = "LATEST"
)

// prop value enum
func (o *ResetOffsetBody) validateToEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, resetOffsetBodyTypeToPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (o *ResetOffsetBody) validateTo(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"to", "body", o.To); err != nil {
		return err
	}

	// value enum
	if err := o.validateToEnum("body"+"."+"to", "body", *o.To); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this reset offset body based on context it is used
func (o *ResetOffsetBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *ResetOffsetBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ResetOffsetBody) UnmarshalBinary(b []byte) error {
	var res ResetOffsetBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
