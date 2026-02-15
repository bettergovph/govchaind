package types

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		AllowedOcidPrefixes:        []string{"ocds-ph-"},
		MaxReleaseSize:             1048576, // 1MB
		RequireEntityVerification:  false,
		DefaultCurrency:            "PHP",
		MinTenderersOpen:           3,
		AlternativeMethodThreshold: "1000000",  // 1M PHP
		PublicBiddingThreshold:     "50000000", // 50M PHP
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}
