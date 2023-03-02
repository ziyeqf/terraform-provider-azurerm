package galleryapplicationversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GalleryApplicationVersionProperties struct {
	ProvisioningState *ProvisioningState                         `json:"provisioningState,omitempty"`
	PublishingProfile GalleryApplicationVersionPublishingProfile `json:"publishingProfile"`
	ReplicationStatus *ReplicationStatus                         `json:"replicationStatus,omitempty"`
}
