package share

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ShareProperties struct {
	CreatedAt         *string            `json:"createdAt,omitempty"`
	Description       *string            `json:"description,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	ShareKind         *ShareKind         `json:"shareKind,omitempty"`
	Terms             *string            `json:"terms,omitempty"`
	UserEmail         *string            `json:"userEmail,omitempty"`
	UserName          *string            `json:"userName,omitempty"`
}

func (o *ShareProperties) GetCreatedAtAsTime() (*time.Time, error) {
	if o.CreatedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *ShareProperties) SetCreatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedAt = &formatted
}
