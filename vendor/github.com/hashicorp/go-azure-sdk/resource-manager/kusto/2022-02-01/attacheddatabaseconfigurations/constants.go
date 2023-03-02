package attacheddatabaseconfigurations

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AttachedDatabaseType string

const (
	AttachedDatabaseTypeMicrosoftPointKustoClustersAttachedDatabaseConfigurations AttachedDatabaseType = "Microsoft.Kusto/clusters/attachedDatabaseConfigurations"
)

func PossibleValuesForAttachedDatabaseType() []string {
	return []string{
		string(AttachedDatabaseTypeMicrosoftPointKustoClustersAttachedDatabaseConfigurations),
	}
}

func parseAttachedDatabaseType(input string) (*AttachedDatabaseType, error) {
	vals := map[string]AttachedDatabaseType{
		"microsoft.kusto/clusters/attacheddatabaseconfigurations": AttachedDatabaseTypeMicrosoftPointKustoClustersAttachedDatabaseConfigurations,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AttachedDatabaseType(input)
	return &out, nil
}

type DefaultPrincipalsModificationKind string

const (
	DefaultPrincipalsModificationKindNone    DefaultPrincipalsModificationKind = "None"
	DefaultPrincipalsModificationKindReplace DefaultPrincipalsModificationKind = "Replace"
	DefaultPrincipalsModificationKindUnion   DefaultPrincipalsModificationKind = "Union"
)

func PossibleValuesForDefaultPrincipalsModificationKind() []string {
	return []string{
		string(DefaultPrincipalsModificationKindNone),
		string(DefaultPrincipalsModificationKindReplace),
		string(DefaultPrincipalsModificationKindUnion),
	}
}

func parseDefaultPrincipalsModificationKind(input string) (*DefaultPrincipalsModificationKind, error) {
	vals := map[string]DefaultPrincipalsModificationKind{
		"none":    DefaultPrincipalsModificationKindNone,
		"replace": DefaultPrincipalsModificationKindReplace,
		"union":   DefaultPrincipalsModificationKindUnion,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DefaultPrincipalsModificationKind(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateMoving    ProvisioningState = "Moving"
	ProvisioningStateRunning   ProvisioningState = "Running"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateMoving),
		string(ProvisioningStateRunning),
		string(ProvisioningStateSucceeded),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"creating":  ProvisioningStateCreating,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"moving":    ProvisioningStateMoving,
		"running":   ProvisioningStateRunning,
		"succeeded": ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type Reason string

const (
	ReasonAlreadyExists Reason = "AlreadyExists"
	ReasonInvalid       Reason = "Invalid"
)

func PossibleValuesForReason() []string {
	return []string{
		string(ReasonAlreadyExists),
		string(ReasonInvalid),
	}
}

func parseReason(input string) (*Reason, error) {
	vals := map[string]Reason{
		"alreadyexists": ReasonAlreadyExists,
		"invalid":       ReasonInvalid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Reason(input)
	return &out, nil
}
