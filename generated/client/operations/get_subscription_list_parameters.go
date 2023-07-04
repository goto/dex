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

// NewGetSubscriptionListParams creates a new GetSubscriptionListParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetSubscriptionListParams() *GetSubscriptionListParams {
	return &GetSubscriptionListParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetSubscriptionListParamsWithTimeout creates a new GetSubscriptionListParams object
// with the ability to set a timeout on a request.
func NewGetSubscriptionListParamsWithTimeout(timeout time.Duration) *GetSubscriptionListParams {
	return &GetSubscriptionListParams{
		timeout: timeout,
	}
}

// NewGetSubscriptionListParamsWithContext creates a new GetSubscriptionListParams object
// with the ability to set a context for a request.
func NewGetSubscriptionListParamsWithContext(ctx context.Context) *GetSubscriptionListParams {
	return &GetSubscriptionListParams{
		Context: ctx,
	}
}

// NewGetSubscriptionListParamsWithHTTPClient creates a new GetSubscriptionListParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetSubscriptionListParamsWithHTTPClient(client *http.Client) *GetSubscriptionListParams {
	return &GetSubscriptionListParams{
		HTTPClient: client,
	}
}

/*
GetSubscriptionListParams contains all the parameters to send to the API endpoint

	for the get subscription list operation.

	Typically these are written to a http.Request.
*/
type GetSubscriptionListParams struct {

	/* GroupID.

	   Shield's group id
	*/
	GroupID *string

	/* ResourceID.

	   Resource unique identifier, e.g. firehose's urn
	*/
	ResourceID *string

	/* ResourceType.

	   Resource type, e.g. firehose, dagger, optimus, etc. This is required if "resource_id" is passed.
	*/
	ResourceType *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get subscription list params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetSubscriptionListParams) WithDefaults() *GetSubscriptionListParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get subscription list params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetSubscriptionListParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get subscription list params
func (o *GetSubscriptionListParams) WithTimeout(timeout time.Duration) *GetSubscriptionListParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get subscription list params
func (o *GetSubscriptionListParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get subscription list params
func (o *GetSubscriptionListParams) WithContext(ctx context.Context) *GetSubscriptionListParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get subscription list params
func (o *GetSubscriptionListParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get subscription list params
func (o *GetSubscriptionListParams) WithHTTPClient(client *http.Client) *GetSubscriptionListParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get subscription list params
func (o *GetSubscriptionListParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithGroupID adds the groupID to the get subscription list params
func (o *GetSubscriptionListParams) WithGroupID(groupID *string) *GetSubscriptionListParams {
	o.SetGroupID(groupID)
	return o
}

// SetGroupID adds the groupId to the get subscription list params
func (o *GetSubscriptionListParams) SetGroupID(groupID *string) {
	o.GroupID = groupID
}

// WithResourceID adds the resourceID to the get subscription list params
func (o *GetSubscriptionListParams) WithResourceID(resourceID *string) *GetSubscriptionListParams {
	o.SetResourceID(resourceID)
	return o
}

// SetResourceID adds the resourceId to the get subscription list params
func (o *GetSubscriptionListParams) SetResourceID(resourceID *string) {
	o.ResourceID = resourceID
}

// WithResourceType adds the resourceType to the get subscription list params
func (o *GetSubscriptionListParams) WithResourceType(resourceType *string) *GetSubscriptionListParams {
	o.SetResourceType(resourceType)
	return o
}

// SetResourceType adds the resourceType to the get subscription list params
func (o *GetSubscriptionListParams) SetResourceType(resourceType *string) {
	o.ResourceType = resourceType
}

// WriteToRequest writes these params to a swagger request
func (o *GetSubscriptionListParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.GroupID != nil {

		// query param group_id
		var qrGroupID string

		if o.GroupID != nil {
			qrGroupID = *o.GroupID
		}
		qGroupID := qrGroupID
		if qGroupID != "" {

			if err := r.SetQueryParam("group_id", qGroupID); err != nil {
				return err
			}
		}
	}

	if o.ResourceID != nil {

		// query param resource_id
		var qrResourceID string

		if o.ResourceID != nil {
			qrResourceID = *o.ResourceID
		}
		qResourceID := qrResourceID
		if qResourceID != "" {

			if err := r.SetQueryParam("resource_id", qResourceID); err != nil {
				return err
			}
		}
	}

	if o.ResourceType != nil {

		// query param resource_type
		var qrResourceType string

		if o.ResourceType != nil {
			qrResourceType = *o.ResourceType
		}
		qResourceType := qrResourceType
		if qResourceType != "" {

			if err := r.SetQueryParam("resource_type", qResourceType); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
