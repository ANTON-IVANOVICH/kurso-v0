// Package service holds the application layer: use cases that orchestrate
// domain entities and ports. Inbound adapters (HTTP) call into these services;
// services call out through domain ports implemented by driven adapters.
//
// Kept as a thin, explicit boundary so transports stay free of business logic.
package service
