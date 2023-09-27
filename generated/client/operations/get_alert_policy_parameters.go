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

// NewGetAlertPolicyParams creates a new GetAlertPolicyParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetAlertPolicyParams() *GetAlertPolicyParams {
	return &GetAlertPolicyParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetAlertPolicyParamsWithTimeout creates a new GetAlertPolicyParams object
// with the ability to set a timeout on a request.
func NewGetAlertPolicyParamsWithTimeout(timeout time.Duration) *GetAlertPolicyParams {
	return &GetAlertPolicyParams{
		timeout: timeout,
	}
}

// NewGetAlertPolicyParamsWithContext creates a new GetAlertPolicyParams object
// with the ability to set a context for a request.
func NewGetAlertPolicyParamsWithContext(ctx context.Context) *GetAlertPolicyParams {
	return &GetAlertPolicyParams{
		Context: ctx,
	}
}

// NewGetAlertPolicyParamsWithHTTPClient creates a new GetAlertPolicyParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetAlertPolicyParamsWithHTTPClient(client *http.Client) *GetAlertPolicyParams {
	return &GetAlertPolicyParams{
		HTTPClient: client,
	}
}

/*
GetAlertPolicyParams contains all the parameters to send to the API endpoint

	for the get alert policy operation.

	Typically these are written to a http.Request.
*/
type GetAlertPolicyParams struct {

	/* ProjectSlug.

	   Shield project slug.
	*/
	ProjectSlug string

	/* ResourceUrn.

	   Siren resource identifier.
	*/
	ResourceUrn string

	/* Template.

	   Siren template tag.
	*/
	Template string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get alert policy params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetAlertPolicyParams) WithDefaults() *GetAlertPolicyParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get alert policy params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetAlertPolicyParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get alert policy params
func (o *GetAlertPolicyParams) WithTimeout(timeout time.Duration) *GetAlertPolicyParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get alert policy params
func (o *GetAlertPolicyParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get alert policy params
func (o *GetAlertPolicyParams) WithContext(ctx context.Context) *GetAlertPolicyParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get alert policy params
func (o *GetAlertPolicyParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get alert policy params
func (o *GetAlertPolicyParams) WithHTTPClient(client *http.Client) *GetAlertPolicyParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get alert policy params
func (o *GetAlertPolicyParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithProjectSlug adds the projectSlug to the get alert policy params
func (o *GetAlertPolicyParams) WithProjectSlug(projectSlug string) *GetAlertPolicyParams {
	o.SetProjectSlug(projectSlug)
	return o
}

// SetProjectSlug adds the projectSlug to the get alert policy params
func (o *GetAlertPolicyParams) SetProjectSlug(projectSlug string) {
	o.ProjectSlug = projectSlug
}

// WithResourceUrn adds the resourceUrn to the get alert policy params
func (o *GetAlertPolicyParams) WithResourceUrn(resourceUrn string) *GetAlertPolicyParams {
	o.SetResourceUrn(resourceUrn)
	return o
}

// SetResourceUrn adds the resourceUrn to the get alert policy params
func (o *GetAlertPolicyParams) SetResourceUrn(resourceUrn string) {
	o.ResourceUrn = resourceUrn
}

// WithTemplate adds the template to the get alert policy params
func (o *GetAlertPolicyParams) WithTemplate(template string) *GetAlertPolicyParams {
	o.SetTemplate(template)
	return o
}

// SetTemplate adds the template to the get alert policy params
func (o *GetAlertPolicyParams) SetTemplate(template string) {
	o.Template = template
}

// WriteToRequest writes these params to a swagger request
func (o *GetAlertPolicyParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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

	// query param template
	qrTemplate := o.Template
	qTemplate := qrTemplate
	if qTemplate != "" {

		if err := r.SetQueryParam("template", qTemplate); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
