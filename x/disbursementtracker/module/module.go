package disbursementtracker

import "govchain/x/disbursementtracker/types"

// AppModule will coordinate the disbursement tracker features once wiring is complete.
type AppModule struct{}

// NewAppModule yields a placeholder AppModule instance.
func NewAppModule() AppModule {
	return AppModule{}
}

// Name returns the module name.
func (AppModule) Name() string {
	return types.ModuleName
}
