package protecteditems

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/polling"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreateOrUpdateOperationResponse struct {
	Poller       polling.LongRunningPoller
	HttpResponse *http.Response
	Model        *ProtectedItemResource
}

// CreateOrUpdate ...
func (c ProtectedItemsClient) CreateOrUpdate(ctx context.Context, id ProtectedItemId, input ProtectedItemResource) (result CreateOrUpdateOperationResponse, err error) {
	req, err := c.preparerForCreateOrUpdate(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "protecteditems.ProtectedItemsClient", "CreateOrUpdate", nil, "Failure preparing request")
		return
	}

	result, err = c.senderForCreateOrUpdate(ctx, req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "protecteditems.ProtectedItemsClient", "CreateOrUpdate", result.HttpResponse, "Failure sending request")
		return
	}

	return
}

// CreateOrUpdateThenPoll performs CreateOrUpdate then polls until it's completed
func (c ProtectedItemsClient) CreateOrUpdateThenPoll(ctx context.Context, id ProtectedItemId, input ProtectedItemResource) error {
	result, err := c.CreateOrUpdate(ctx, id, input)
	if err != nil {
		return fmt.Errorf("performing CreateOrUpdate: %+v", err)
	}

	if err := result.Poller.PollUntilDone(); err != nil {
		return fmt.Errorf("polling after CreateOrUpdate: %+v", err)
	}

	return nil
}

// preparerForCreateOrUpdate prepares the CreateOrUpdate request.
func (c ProtectedItemsClient) preparerForCreateOrUpdate(ctx context.Context, id ProtectedItemId, input ProtectedItemResource) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// senderForCreateOrUpdate sends the CreateOrUpdate request. The method will close the
// http.Response Body if it receives an error.
func (c ProtectedItemsClient) senderForCreateOrUpdate(ctx context.Context, req *http.Request) (future CreateOrUpdateOperationResponse, err error) {
	var resp *http.Response
	resp, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		return
	}

	future.Poller, err = polling.NewPollerFromResponse(ctx, resp, c.Client, req.Method)
	return
}
