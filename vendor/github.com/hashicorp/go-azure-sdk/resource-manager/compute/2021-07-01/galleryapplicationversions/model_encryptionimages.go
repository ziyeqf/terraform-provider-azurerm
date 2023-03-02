package galleryapplicationversions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EncryptionImages struct {
	DataDiskImages *[]DataDiskImageEncryption `json:"dataDiskImages,omitempty"`
	OsDiskImage    *DiskImageEncryption       `json:"osDiskImage,omitempty"`
}
