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

// NewGetAlertsParams creates a new GetAlertsParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetAlertsParams() *GetAlertsParams {
	return &GetAlertsParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetAlertsParamsWithTimeout creates a new GetAlertsParams object
// with the ability to set a timeout on a request.
func NewGetAlertsParamsWithTimeout(timeout time.Duration) *GetAlertsParams {
	return &GetAlertsParams{
		timeout: timeout,
	}
}

// NewGetAlertsParamsWithContext creates a new GetAlertsParams object
// with the ability to set a context for a request.
func NewGetAlertsParamsWithContext(ctx context.Context) *GetAlertsParams {
	return &GetAlertsParams{
		Context: ctx,
	}
}

// NewGetAlertsParamsWithHTTPClient creates a new GetAlertsParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetAlertsParamsWithHTTPClient(client *http.Client) *GetAlertsParams {
	return &GetAlertsParams{
		HTTPClient: client,
	}
}

/*
GetAlertsParams contains all the parameters to send to the API endpoint

	for the get alerts operation.

	Typically these are written to a http.Request.
*/
type GetAlertsParams struct {

	/* ProjectSlug.

	   Shield project slug.
	*/
	ProjectSlug string

	/* ResourceUrn.

	   Siren resource identifier.
	*/
	ResourceUrn string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get alerts params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetAlertsParams) WithDefaults() *GetAlertsParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get alerts params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetAlertsParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get alerts params
func (o *GetAlertsParams) WithTimeout(timeout time.Duration) *GetAlertsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get alerts params
func (o *GetAlertsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get alerts params
func (o *GetAlertsParams) WithContext(ctx context.Context) *GetAlertsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get alerts params
func (o *GetAlertsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get alerts params
func (o *GetAlertsParams) WithHTTPClient(client *http.Client) *GetAlertsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get alerts params
func (o *GetAlertsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithProjectSlug adds the projectSlug to the get alerts params
func (o *GetAlertsParams) WithProjectSlug(projectSlug string) *GetAlertsParams {
	o.SetProjectSlug(projectSlug)
	return o
}

// SetProjectSlug adds the projectSlug to the get alerts params
func (o *GetAlertsParams) SetProjectSlug(projectSlug string) {
	o.ProjectSlug = projectSlug
}

// WithResourceUrn adds the resourceUrn to the get alerts params
func (o *GetAlertsParams) WithResourceUrn(resourceUrn string) *GetAlertsParams {
	o.SetResourceUrn(resourceUrn)
	return o
}

// SetResourceUrn adds the resourceUrn to the get alerts params
func (o *GetAlertsParams) SetResourceUrn(resourceUrn string) {
	o.ResourceUrn = resourceUrn
}

// WriteToRequest writes these params to a swagger request
func (o *GetAlertsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param projectSlug
	if err := r.SetPathParam("projectSlug", o.ProjectSlug); err != nil {
		return err
	}

	// path param resourceUrn
	if err := r.SetPathParam("resourceUrn", o.ResourceUrn); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
