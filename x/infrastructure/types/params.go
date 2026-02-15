package types

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		DefaultCurrency:            "PHP",
		RequireProjectVerification: false,
		MaxProjectSize:             2097152, // 2MB
		AllowedProjectPrefixes:     []string{"ph-dpwh-"},
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}
