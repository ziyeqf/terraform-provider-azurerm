package galleryapplicationversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TargetRegion struct {
	Encryption           *EncryptionImages   `json:"encryption,omitempty"`
	Name                 string              `json:"name"`
	RegionalReplicaCount *int64              `json:"regionalReplicaCount,omitempty"`
	StorageAccountType   *StorageAccountType `json:"storageAccountType,omitempty"`
}
