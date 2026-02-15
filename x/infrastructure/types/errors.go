package types

import errorsmod "cosmossdk.io/errors"

var (
	ErrInvalidProjectID      = errorsmod.Register(ModuleName, 1200, "invalid project ID format")
	ErrProjectNotFound       = errorsmod.Register(ModuleName, 1201, "infrastructure project not found")
	ErrProjectAlreadyExists  = errorsmod.Register(ModuleName, 1202, "infrastructure project already exists")
	ErrUnauthorizedPublisher = errorsmod.Register(ModuleName, 1203, "unauthorized publisher")
	ErrProjectTooLarge       = errorsmod.Register(ModuleName, 1204, "project data exceeds maximum size")
)
