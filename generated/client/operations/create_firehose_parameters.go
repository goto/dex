// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/goto/dex/generated/models"
)

// NewCreateFirehoseParams creates a new CreateFirehoseParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewCreateFirehoseParams() *CreateFirehoseParams {
	return &CreateFirehoseParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewCreateFirehoseParamsWithTimeout creates a new CreateFirehoseParams object
// with the ability to set a timeout on a request.
func NewCreateFirehoseParamsWithTimeout(timeout time.Duration) *CreateFirehoseParams {
	return &CreateFirehoseParams{
		timeout: timeout,
	}
}

// NewCreateFirehoseParamsWithContext creates a new CreateFirehoseParams object
// with the ability to set a context for a request.
func NewCreateFirehoseParamsWithContext(ctx context.Context) *CreateFirehoseParams {
	return &CreateFirehoseParams{
		Context: ctx,
	}
}

// NewCreateFirehoseParamsWithHTTPClient creates a new CreateFirehoseParams object
// with the ability to set a custom HTTPClient for a request.
func NewCreateFirehoseParamsWithHTTPClient(client *http.Client) *CreateFirehoseParams {
	return &CreateFirehoseParams{
		HTTPClient: client,
	}
}

/*
CreateFirehoseParams contains all the parameters to send to the API endpoint

	for the create firehose operation.

	Typically these are written to a http.Request.
*/
type CreateFirehoseParams struct {

	// Body.
	Body *models.Firehose

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the create firehose params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateFirehoseParams) WithDefaults() *CreateFirehoseParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the create firehose params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *CreateFirehoseParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the create firehose params
func (o *CreateFirehoseParams) WithTimeout(timeout time.Duration) *CreateFirehoseParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create firehose params
func (o *CreateFirehoseParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create firehose params
func (o *CreateFirehoseParams) WithContext(ctx context.Context) *CreateFirehoseParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create firehose params
func (o *CreateFirehoseParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create firehose params
func (o *CreateFirehoseParams) WithHTTPClient(client *http.Client) *CreateFirehoseParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create firehose params
func (o *CreateFirehoseParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the create firehose params
func (o *CreateFirehoseParams) WithBody(body *models.Firehose) *CreateFirehoseParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the create firehose params
func (o *CreateFirehoseParams) SetBody(body *models.Firehose) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *CreateFirehoseParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
