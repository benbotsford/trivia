// Package billing defines the entitlement contract for the trivia platform.
// The real Stripe-backed implementation lives here once billing goes live.
// For now, NoopChecker grants access to everything.
package billing

import "context"

// EntitlementChecker answers gate questions about what a host is allowed to do.
// Handlers call this; swapping the implementation is the only change needed
// when Stripe billing is introduced.
type EntitlementChecker interface {
	// CanCreateGame returns true if the host is allowed to create a new game.
	CanCreateGame(ctx context.Context, userID string) (bool, error)
	// CanUseFeature returns true if the host is allowed to use a named feature.
	CanUseFeature(ctx context.Context, userID, feature string) (bool, error)
}

// NoopChecker is a pass-through implementation that always grants access.
// Used until Stripe billing is implemented.
type NoopChecker struct{}

var _ EntitlementChecker = (*NoopChecker)(nil)

func (NoopChecker) CanCreateGame(_ context.Context, _ string) (bool, error) {
	return true, nil
}

func (NoopChecker) CanUseFeature(_ context.Context, _, _ string) (bool, error) {
	return true, nil
}
