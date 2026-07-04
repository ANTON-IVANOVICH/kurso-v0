// Package domain is the core of the hexagon: pure business entities and the
// port interfaces (repositories, gateways) that the application depends on.
//
// It has NO dependencies on frameworks, transports, or databases. Adapters in
// internal/adapter and internal/platform implement these ports; use cases in
// internal/service orchestrate them.
//
// Subpackages (added as features land, from Stage 1 onward), e.g.:
//
//	internal/domain/catalog   — currencies, directions, exchangers
//	internal/domain/rates     — rate snapshots and history
//	internal/domain/alert     — user alerts and matching rules
package domain
