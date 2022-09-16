package model

import "sync"

// Model holds an abstracted data model representing the translation
// of various types of Kubernetes config to Cilium config.
type Model struct {
	HTTP []HTTPListener
	mu   sync.Mutex
}

// HTTPListener holds configuration for any listener that terminates and proxies HTTP
// including HTTP and HTTPS.
type HTTPListener struct {
	// Name of the HTTPListener
	Name string
	// Sources is a slice of fully qualified resources this HTTPListener is sourced
	// from.
	Sources []FullyQualifiedResource
	// IPAddress that the listener should listen on.
	// TODO(youngnick): Should this be a list, or should we have one listener per address?
	// The string must be parseable as an IP address.
	Address string
	// Port on which the service can be expected to be accessed by clients.
	Port uint32
	// VirtualHost details that the listener should match
	// Each VirtualHost contains a hostname, and optionally TLS secret details.
	// If no VirtualHosts contain TLS details, then this is an insecure HTTPListener
	VirtualHosts []VirtualHost
	// Routes associated with HTTP traffic to the service.
	// An empty list means that traffic will not be routed.
	Routes []HTTPRoute
}

type VirtualHost struct {
	Hostname string
	TLS      *TLSSecret
}

// FullyQualifiedResource stores the full details of a Kubernetes resource, including
// the Group, Version, and Kind.
// Namespace must be set to the empty string for cluster-scoped resources.
type FullyQualifiedResource struct {
	Name      string
	Namespace string
	Group     string
	Version   string
	Kind      string
}

// TLSSecret holds a reference to a secret containing a TLS keypair.
type TLSSecret struct {
	Name      string
	Namespace string
}

// HTTPRoute holds all the details needed to route HTTP traffic to a backend.
type HTTPRoute struct {
	Name string
	// PathMatch specifies that the HTTPRoute should match a path.
	PathMatch StringMatch
	Backends  []Backend
}

// StringMatch describes various types of string matching.
// Only one field may be set.
type StringMatch struct {
	Prefix string
	Exact  string
	Regex  string
}

// Backend holds a Kubernetes Service that points to a backend for traffic.
type Backend struct {
	// Name of the Service.
	Name string
	// Namespace of the Service.
	Namespace string
	// Port is the port on the Service to connect to.
	// If unset, the same port as the top-level Listener will be used.
	Port uint32
}
