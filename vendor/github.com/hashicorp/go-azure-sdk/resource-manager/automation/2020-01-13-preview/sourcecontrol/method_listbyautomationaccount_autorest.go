package sourcecontrol

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByAutomationAccountOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]SourceControl

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByAutomationAccountOperationResponse, error)
}

type ListByAutomationAccountCompleteResult struct {
	Items []SourceControl
}

func (r ListByAutomationAccountOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByAutomationAccountOperationResponse) LoadMore(ctx context.Context) (resp ListByAutomationAccountOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListByAutomationAccountOperationOptions struct {
	Filter *string
}

func DefaultListByAutomationAccountOperationOptions() ListByAutomationAccountOperationOptions {
	return ListByAutomationAccountOperationOptions{}
}

func (o ListByAutomationAccountOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListByAutomationAccountOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	return out
}

// ListByAutomationAccount ...
func (c SourceControlClient) ListByAutomationAccount(ctx context.Context, id AutomationAccountId, options ListByAutomationAccountOperationOptions) (resp ListByAutomationAccountOperationResponse, err error) {
	req, err := c.preparerForListByAutomationAccount(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "sourcecontrol.SourceControlClient", "ListByAutomationAccount", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "sourcecontrol.SourceControlClient", "ListByAutomationAccount", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByAutomationAccount(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "sourcecontrol.SourceControlClient", "ListByAutomationAccount", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListByAutomationAccount prepares the ListByAutomationAccount request.
func (c SourceControlClient) preparerForListByAutomationAccount(ctx context.Context, id AutomationAccountId, options ListByAutomationAccountOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/sourceControls", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByAutomationAccountWithNextLink prepares the ListByAutomationAccount request with the given nextLink token.
func (c SourceControlClient) preparerForListByAutomationAccountWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
	uri, err := url.Parse(nextLink)
	if err != nil {
		return nil, fmt.Errorf("parsing nextLink %q: %+v", nextLink, err)
	}
	queryParameters := map[string]interface{}{}
	for k, v := range uri.Query() {
		if len(v) == 0 {
			continue
		}
		val := v[0]
		val = autorest.Encode("query", val)
		queryParameters[k] = val
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(uri.Path),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForListByAutomationAccount handles the response to the ListByAutomationAccount request. The method always
// closes the http.Response Body.
func (c SourceControlClient) responderForListByAutomationAccount(resp *http.Response) (result ListByAutomationAccountOperationResponse, err error) {
	type page struct {
		Values   []SourceControl `json:"value"`
		NextLink *string         `json:"nextLink"`
	}
	var respObj page
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&respObj),
		autorest.ByClosing())
	result.HttpResponse = resp
	result.Model = &respObj.Values
	result.nextLink = respObj.NextLink
	if respObj.NextLink != nil {
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByAutomationAccountOperationResponse, err error) {
			req, err := c.preparerForListByAutomationAccountWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "sourcecontrol.SourceControlClient", "ListByAutomationAccount", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "sourcecontrol.SourceControlClient", "ListByAutomationAccount", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByAutomationAccount(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "sourcecontrol.SourceControlClient", "ListByAutomationAccount", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListByAutomationAccountComplete retrieves all of the results into a single object
func (c SourceControlClient) ListByAutomationAccountComplete(ctx context.Context, id AutomationAccountId, options ListByAutomationAccountOperationOptions) (ListByAutomationAccountCompleteResult, error) {
	return c.ListByAutomationAccountCompleteMatchingPredicate(ctx, id, options, SourceControlOperationPredicate{})
}

// ListByAutomationAccountCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c SourceControlClient) ListByAutomationAccountCompleteMatchingPredicate(ctx context.Context, id AutomationAccountId, options ListByAutomationAccountOperationOptions, predicate SourceControlOperationPredicate) (resp ListByAutomationAccountCompleteResult, err error) {
	items := make([]SourceControl, 0)

	page, err := c.ListByAutomationAccount(ctx, id, options)
	if err != nil {
		err = fmt.Errorf("loading the initial page: %+v", err)
		return
	}
	if page.Model != nil {
		for _, v := range *page.Model {
			if predicate.Matches(v) {
				items = append(items, v)
			}
		}
	}

	for page.HasMore() {
		page, err = page.LoadMore(ctx)
		if err != nil {
			err = fmt.Errorf("loading the next page: %+v", err)
			return
		}

		if page.Model != nil {
			for _, v := range *page.Model {
				if predicate.Matches(v) {
					items = append(items, v)
				}
			}
		}
	}

	out := ListByAutomationAccountCompleteResult{
		Items: items,
	}
	return out, nil
}
