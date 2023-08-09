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

// NewGetOptimusJobParams creates a new GetOptimusJobParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetOptimusJobParams() *GetOptimusJobParams {
	return &GetOptimusJobParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetOptimusJobParamsWithTimeout creates a new GetOptimusJobParams object
// with the ability to set a timeout on a request.
func NewGetOptimusJobParamsWithTimeout(timeout time.Duration) *GetOptimusJobParams {
	return &GetOptimusJobParams{
		timeout: timeout,
	}
}

// NewGetOptimusJobParamsWithContext creates a new GetOptimusJobParams object
// with the ability to set a context for a request.
func NewGetOptimusJobParamsWithContext(ctx context.Context) *GetOptimusJobParams {
	return &GetOptimusJobParams{
		Context: ctx,
	}
}

// NewGetOptimusJobParamsWithHTTPClient creates a new GetOptimusJobParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetOptimusJobParamsWithHTTPClient(client *http.Client) *GetOptimusJobParams {
	return &GetOptimusJobParams{
		HTTPClient: client,
	}
}

/*
GetOptimusJobParams contains all the parameters to send to the API endpoint

	for the get optimus job operation.

	Typically these are written to a http.Request.
*/
type GetOptimusJobParams struct {

	/* Job.

	   Optimus job's name.
	*/
	Job string

	/* Project.

	   Unique identifier of the project.
	*/
	Project string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get optimus job params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetOptimusJobParams) WithDefaults() *GetOptimusJobParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get optimus job params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetOptimusJobParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get optimus job params
func (o *GetOptimusJobParams) WithTimeout(timeout time.Duration) *GetOptimusJobParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get optimus job params
func (o *GetOptimusJobParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get optimus job params
func (o *GetOptimusJobParams) WithContext(ctx context.Context) *GetOptimusJobParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get optimus job params
func (o *GetOptimusJobParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get optimus job params
func (o *GetOptimusJobParams) WithHTTPClient(client *http.Client) *GetOptimusJobParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get optimus job params
func (o *GetOptimusJobParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithJob adds the job to the get optimus job params
func (o *GetOptimusJobParams) WithJob(job string) *GetOptimusJobParams {
	o.SetJob(job)
	return o
}

// SetJob adds the job to the get optimus job params
func (o *GetOptimusJobParams) SetJob(job string) {
	o.Job = job
}

// WithProject adds the project to the get optimus job params
func (o *GetOptimusJobParams) WithProject(project string) *GetOptimusJobParams {
	o.SetProject(project)
	return o
}

// SetProject adds the project to the get optimus job params
func (o *GetOptimusJobParams) SetProject(project string) {
	o.Project = project
}

// WriteToRequest writes these params to a swagger request
func (o *GetOptimusJobParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param job
	if err := r.SetPathParam("job", o.Job); err != nil {
		return err
	}

	// path param project
	if err := r.SetPathParam("project", o.Project); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
