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

// NewUpsertFirehoseAlertPolicyParams creates a new UpsertFirehoseAlertPolicyParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewUpsertFirehoseAlertPolicyParams() *UpsertFirehoseAlertPolicyParams {
	return &UpsertFirehoseAlertPolicyParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewUpsertFirehoseAlertPolicyParamsWithTimeout creates a new UpsertFirehoseAlertPolicyParams object
// with the ability to set a timeout on a request.
func NewUpsertFirehoseAlertPolicyParamsWithTimeout(timeout time.Duration) *UpsertFirehoseAlertPolicyParams {
	return &UpsertFirehoseAlertPolicyParams{
		timeout: timeout,
	}
}

// NewUpsertFirehoseAlertPolicyParamsWithContext creates a new UpsertFirehoseAlertPolicyParams object
// with the ability to set a context for a request.
func NewUpsertFirehoseAlertPolicyParamsWithContext(ctx context.Context) *UpsertFirehoseAlertPolicyParams {
	return &UpsertFirehoseAlertPolicyParams{
		Context: ctx,
	}
}

// NewUpsertFirehoseAlertPolicyParamsWithHTTPClient creates a new UpsertFirehoseAlertPolicyParams object
// with the ability to set a custom HTTPClient for a request.
func NewUpsertFirehoseAlertPolicyParamsWithHTTPClient(client *http.Client) *UpsertFirehoseAlertPolicyParams {
	return &UpsertFirehoseAlertPolicyParams{
		HTTPClient: client,
	}
}

/*
UpsertFirehoseAlertPolicyParams contains all the parameters to send to the API endpoint

	for the upsert firehose alert policy operation.

	Typically these are written to a http.Request.
*/
type UpsertFirehoseAlertPolicyParams struct {

	// Body.
	Body *models.AlertPolicy

	/* FirehoseUrn.

	   URN of the firehose.
	*/
	FirehoseUrn string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the upsert firehose alert policy params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpsertFirehoseAlertPolicyParams) WithDefaults() *UpsertFirehoseAlertPolicyParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the upsert firehose alert policy params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpsertFirehoseAlertPolicyParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the upsert firehose alert policy params
func (o *UpsertFirehoseAlertPolicyParams) WithTimeout(timeout time.Duration) *UpsertFirehoseAlertPolicyParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the upsert firehose alert policy params
func (o *UpsertFirehoseAlertPolicyParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the upsert firehose alert policy params
func (o *UpsertFirehoseAlertPolicyParams) WithContext(ctx context.Context) *UpsertFirehoseAlertPolicyParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the upsert firehose alert policy params
func (o *UpsertFirehoseAlertPolicyParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the upsert firehose alert policy params
func (o *UpsertFirehoseAlertPolicyParams) WithHTTPClient(client *http.Client) *UpsertFirehoseAlertPolicyParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the upsert firehose alert policy params
func (o *UpsertFirehoseAlertPolicyParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the upsert firehose alert policy params
func (o *UpsertFirehoseAlertPolicyParams) WithBody(body *models.AlertPolicy) *UpsertFirehoseAlertPolicyParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the upsert firehose alert policy params
func (o *UpsertFirehoseAlertPolicyParams) SetBody(body *models.AlertPolicy) {
	o.Body = body
}

// WithFirehoseUrn adds the firehoseUrn to the upsert firehose alert policy params
func (o *UpsertFirehoseAlertPolicyParams) WithFirehoseUrn(firehoseUrn string) *UpsertFirehoseAlertPolicyParams {
	o.SetFirehoseUrn(firehoseUrn)
	return o
}

// SetFirehoseUrn adds the firehoseUrn to the upsert firehose alert policy params
func (o *UpsertFirehoseAlertPolicyParams) SetFirehoseUrn(firehoseUrn string) {
	o.FirehoseUrn = firehoseUrn
}

// WriteToRequest writes these params to a swagger request
func (o *UpsertFirehoseAlertPolicyParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	// path param firehoseUrn
	if err := r.SetPathParam("firehoseUrn", o.FirehoseUrn); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
