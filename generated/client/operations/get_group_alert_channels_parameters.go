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

// NewGetGroupAlertChannelsParams creates a new GetGroupAlertChannelsParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetGroupAlertChannelsParams() *GetGroupAlertChannelsParams {
	return &GetGroupAlertChannelsParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetGroupAlertChannelsParamsWithTimeout creates a new GetGroupAlertChannelsParams object
// with the ability to set a timeout on a request.
func NewGetGroupAlertChannelsParamsWithTimeout(timeout time.Duration) *GetGroupAlertChannelsParams {
	return &GetGroupAlertChannelsParams{
		timeout: timeout,
	}
}

// NewGetGroupAlertChannelsParamsWithContext creates a new GetGroupAlertChannelsParams object
// with the ability to set a context for a request.
func NewGetGroupAlertChannelsParamsWithContext(ctx context.Context) *GetGroupAlertChannelsParams {
	return &GetGroupAlertChannelsParams{
		Context: ctx,
	}
}

// NewGetGroupAlertChannelsParamsWithHTTPClient creates a new GetGroupAlertChannelsParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetGroupAlertChannelsParamsWithHTTPClient(client *http.Client) *GetGroupAlertChannelsParams {
	return &GetGroupAlertChannelsParams{
		HTTPClient: client,
	}
}

/*
GetGroupAlertChannelsParams contains all the parameters to send to the API endpoint

	for the get group alert channels operation.

	Typically these are written to a http.Request.
*/
type GetGroupAlertChannelsParams struct {

	/* ID.

	   Shield Group's ID
	*/
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get group alert channels params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetGroupAlertChannelsParams) WithDefaults() *GetGroupAlertChannelsParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get group alert channels params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetGroupAlertChannelsParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get group alert channels params
func (o *GetGroupAlertChannelsParams) WithTimeout(timeout time.Duration) *GetGroupAlertChannelsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get group alert channels params
func (o *GetGroupAlertChannelsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get group alert channels params
func (o *GetGroupAlertChannelsParams) WithContext(ctx context.Context) *GetGroupAlertChannelsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get group alert channels params
func (o *GetGroupAlertChannelsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get group alert channels params
func (o *GetGroupAlertChannelsParams) WithHTTPClient(client *http.Client) *GetGroupAlertChannelsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get group alert channels params
func (o *GetGroupAlertChannelsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the get group alert channels params
func (o *GetGroupAlertChannelsParams) WithID(id string) *GetGroupAlertChannelsParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get group alert channels params
func (o *GetGroupAlertChannelsParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *GetGroupAlertChannelsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
