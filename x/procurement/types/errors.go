package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrInvalidOCID          = errorsmod.Register(ModuleName, 1100, "invalid OCID format")
	ErrProcessNotFound      = errorsmod.Register(ModuleName, 1101, "contracting process not found")
	ErrProcessAlreadyExists = errorsmod.Register(ModuleName, 1102, "contracting process already exists")
	ErrEntityNotFound       = errorsmod.Register(ModuleName, 1103, "entity not found")
	ErrEntityAlreadyExists  = errorsmod.Register(ModuleName, 1104, "entity already exists")
	ErrUnauthorizedPublisher = errorsmod.Register(ModuleName, 1105, "unauthorized publisher")
	ErrReleaseNotFound      = errorsmod.Register(ModuleName, 1106, "release not found")
	ErrInvalidReleaseTag    = errorsmod.Register(ModuleName, 1107, "invalid release tag")
	ErrReleaseTooLarge      = errorsmod.Register(ModuleName, 1108, "release exceeds maximum size")
)
