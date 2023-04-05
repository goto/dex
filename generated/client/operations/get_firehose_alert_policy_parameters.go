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
)

// NewGetFirehoseAlertPolicyParams creates a new GetFirehoseAlertPolicyParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetFirehoseAlertPolicyParams() *GetFirehoseAlertPolicyParams {
	return &GetFirehoseAlertPolicyParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetFirehoseAlertPolicyParamsWithTimeout creates a new GetFirehoseAlertPolicyParams object
// with the ability to set a timeout on a request.
func NewGetFirehoseAlertPolicyParamsWithTimeout(timeout time.Duration) *GetFirehoseAlertPolicyParams {
	return &GetFirehoseAlertPolicyParams{
		timeout: timeout,
	}
}

// NewGetFirehoseAlertPolicyParamsWithContext creates a new GetFirehoseAlertPolicyParams object
// with the ability to set a context for a request.
func NewGetFirehoseAlertPolicyParamsWithContext(ctx context.Context) *GetFirehoseAlertPolicyParams {
	return &GetFirehoseAlertPolicyParams{
		Context: ctx,
	}
}

// NewGetFirehoseAlertPolicyParamsWithHTTPClient creates a new GetFirehoseAlertPolicyParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetFirehoseAlertPolicyParamsWithHTTPClient(client *http.Client) *GetFirehoseAlertPolicyParams {
	return &GetFirehoseAlertPolicyParams{
		HTTPClient: client,
	}
}

/*
GetFirehoseAlertPolicyParams contains all the parameters to send to the API endpoint

	for the get firehose alert policy operation.

	Typically these are written to a http.Request.
*/
type GetFirehoseAlertPolicyParams struct {

	/* FirehoseUrn.

	   URN of the firehose.
	*/
	FirehoseUrn string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get firehose alert policy params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetFirehoseAlertPolicyParams) WithDefaults() *GetFirehoseAlertPolicyParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get firehose alert policy params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetFirehoseAlertPolicyParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get firehose alert policy params
func (o *GetFirehoseAlertPolicyParams) WithTimeout(timeout time.Duration) *GetFirehoseAlertPolicyParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get firehose alert policy params
func (o *GetFirehoseAlertPolicyParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get firehose alert policy params
func (o *GetFirehoseAlertPolicyParams) WithContext(ctx context.Context) *GetFirehoseAlertPolicyParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get firehose alert policy params
func (o *GetFirehoseAlertPolicyParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get firehose alert policy params
func (o *GetFirehoseAlertPolicyParams) WithHTTPClient(client *http.Client) *GetFirehoseAlertPolicyParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get firehose alert policy params
func (o *GetFirehoseAlertPolicyParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithFirehoseUrn adds the firehoseUrn to the get firehose alert policy params
func (o *GetFirehoseAlertPolicyParams) WithFirehoseUrn(firehoseUrn string) *GetFirehoseAlertPolicyParams {
	o.SetFirehoseUrn(firehoseUrn)
	return o
}

// SetFirehoseUrn adds the firehoseUrn to the get firehose alert policy params
func (o *GetFirehoseAlertPolicyParams) SetFirehoseUrn(firehoseUrn string) {
	o.FirehoseUrn = firehoseUrn
}

// WriteToRequest writes these params to a swagger request
func (o *GetFirehoseAlertPolicyParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param firehoseUrn
	if err := r.SetPathParam("firehoseUrn", o.FirehoseUrn); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
