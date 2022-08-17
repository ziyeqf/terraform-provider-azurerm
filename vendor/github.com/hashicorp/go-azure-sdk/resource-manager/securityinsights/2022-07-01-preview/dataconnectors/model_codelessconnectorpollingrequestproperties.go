package dataconnectors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CodelessConnectorPollingRequestProperties struct {
	ApiEndpoint             string       `json:"apiEndpoint"`
	EndTimeAttributeName    *string      `json:"endTimeAttributeName,omitempty"`
	Headers                 *interface{} `json:"headers,omitempty"`
	HttpMethod              string       `json:"httpMethod"`
	QueryParameters         *interface{} `json:"queryParameters,omitempty"`
	QueryParametersTemplate *string      `json:"queryParametersTemplate,omitempty"`
	QueryTimeFormat         string       `json:"queryTimeFormat"`
	QueryWindowInMin        int64        `json:"queryWindowInMin"`
	RateLimitQps            *int64       `json:"rateLimitQps,omitempty"`
	RetryCount              *int64       `json:"retryCount,omitempty"`
	StartTimeAttributeName  *string      `json:"startTimeAttributeName,omitempty"`
	TimeoutInSeconds        *int64       `json:"timeoutInSeconds,omitempty"`
}
