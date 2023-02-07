package packetcorecontrolplane

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthenticationType string

const (
	AuthenticationTypeAAD      AuthenticationType = "AAD"
	AuthenticationTypePassword AuthenticationType = "Password"
)

func PossibleValuesForAuthenticationType() []string {
	return []string{
		string(AuthenticationTypeAAD),
		string(AuthenticationTypePassword),
	}
}

func parseAuthenticationType(input string) (*AuthenticationType, error) {
	vals := map[string]AuthenticationType{
		"aad":      AuthenticationTypeAAD,
		"password": AuthenticationTypePassword,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AuthenticationType(input)
	return &out, nil
}

type BillingSku string

const (
	BillingSkuGFive    BillingSku = "G5"
	BillingSkuGFour    BillingSku = "G4"
	BillingSkuGOne     BillingSku = "G1"
	BillingSkuGOneZero BillingSku = "G10"
	BillingSkuGThree   BillingSku = "G3"
	BillingSkuGTwo     BillingSku = "G2"
	BillingSkuGZero    BillingSku = "G0"
)

func PossibleValuesForBillingSku() []string {
	return []string{
		string(BillingSkuGFive),
		string(BillingSkuGFour),
		string(BillingSkuGOne),
		string(BillingSkuGOneZero),
		string(BillingSkuGThree),
		string(BillingSkuGTwo),
		string(BillingSkuGZero),
	}
}

func parseBillingSku(input string) (*BillingSku, error) {
	vals := map[string]BillingSku{
		"g5":  BillingSkuGFive,
		"g4":  BillingSkuGFour,
		"g1":  BillingSkuGOne,
		"g10": BillingSkuGOneZero,
		"g3":  BillingSkuGThree,
		"g2":  BillingSkuGTwo,
		"g0":  BillingSkuGZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BillingSku(input)
	return &out, nil
}

type CertificateProvisioningState string

const (
	CertificateProvisioningStateFailed         CertificateProvisioningState = "Failed"
	CertificateProvisioningStateNotProvisioned CertificateProvisioningState = "NotProvisioned"
	CertificateProvisioningStateProvisioned    CertificateProvisioningState = "Provisioned"
)

func PossibleValuesForCertificateProvisioningState() []string {
	return []string{
		string(CertificateProvisioningStateFailed),
		string(CertificateProvisioningStateNotProvisioned),
		string(CertificateProvisioningStateProvisioned),
	}
}

func parseCertificateProvisioningState(input string) (*CertificateProvisioningState, error) {
	vals := map[string]CertificateProvisioningState{
		"failed":         CertificateProvisioningStateFailed,
		"notprovisioned": CertificateProvisioningStateNotProvisioned,
		"provisioned":    CertificateProvisioningStateProvisioned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CertificateProvisioningState(input)
	return &out, nil
}

type CoreNetworkType string

const (
	CoreNetworkTypeEPC    CoreNetworkType = "EPC"
	CoreNetworkTypeFiveGC CoreNetworkType = "5GC"
)

func PossibleValuesForCoreNetworkType() []string {
	return []string{
		string(CoreNetworkTypeEPC),
		string(CoreNetworkTypeFiveGC),
	}
}

func parseCoreNetworkType(input string) (*CoreNetworkType, error) {
	vals := map[string]CoreNetworkType{
		"epc": CoreNetworkTypeEPC,
		"5gc": CoreNetworkTypeFiveGC,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CoreNetworkType(input)
	return &out, nil
}

type InstallationState string

const (
	InstallationStateFailed       InstallationState = "Failed"
	InstallationStateInstalled    InstallationState = "Installed"
	InstallationStateInstalling   InstallationState = "Installing"
	InstallationStateReinstalling InstallationState = "Reinstalling"
	InstallationStateRollingBack  InstallationState = "RollingBack"
	InstallationStateUninstalled  InstallationState = "Uninstalled"
	InstallationStateUninstalling InstallationState = "Uninstalling"
	InstallationStateUpdating     InstallationState = "Updating"
	InstallationStateUpgrading    InstallationState = "Upgrading"
)

func PossibleValuesForInstallationState() []string {
	return []string{
		string(InstallationStateFailed),
		string(InstallationStateInstalled),
		string(InstallationStateInstalling),
		string(InstallationStateReinstalling),
		string(InstallationStateRollingBack),
		string(InstallationStateUninstalled),
		string(InstallationStateUninstalling),
		string(InstallationStateUpdating),
		string(InstallationStateUpgrading),
	}
}

func parseInstallationState(input string) (*InstallationState, error) {
	vals := map[string]InstallationState{
		"failed":       InstallationStateFailed,
		"installed":    InstallationStateInstalled,
		"installing":   InstallationStateInstalling,
		"reinstalling": InstallationStateReinstalling,
		"rollingback":  InstallationStateRollingBack,
		"uninstalled":  InstallationStateUninstalled,
		"uninstalling": InstallationStateUninstalling,
		"updating":     InstallationStateUpdating,
		"upgrading":    InstallationStateUpgrading,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := InstallationState(input)
	return &out, nil
}

type PlatformType string

const (
	PlatformTypeAKSNegativeHCI                              PlatformType = "AKS-HCI"
	PlatformTypeThreePNegativeAZURENegativeSTACKNegativeHCI PlatformType = "3P-AZURE-STACK-HCI"
)

func PossibleValuesForPlatformType() []string {
	return []string{
		string(PlatformTypeAKSNegativeHCI),
		string(PlatformTypeThreePNegativeAZURENegativeSTACKNegativeHCI),
	}
}

func parsePlatformType(input string) (*PlatformType, error) {
	vals := map[string]PlatformType{
		"aks-hci":            PlatformTypeAKSNegativeHCI,
		"3p-azure-stack-hci": PlatformTypeThreePNegativeAZURENegativeSTACKNegativeHCI,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PlatformType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateAccepted  ProvisioningState = "Accepted"
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateDeleted   ProvisioningState = "Deleted"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUnknown   ProvisioningState = "Unknown"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateAccepted),
		string(ProvisioningStateCanceled),
		string(ProvisioningStateDeleted),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUnknown),
	}
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"accepted":  ProvisioningStateAccepted,
		"canceled":  ProvisioningStateCanceled,
		"deleted":   ProvisioningStateDeleted,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
		"unknown":   ProvisioningStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}
