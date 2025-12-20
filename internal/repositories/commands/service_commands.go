package commands

// UpdateServiceCommand represents a partial update request for a service. This
// struct should be used internally only by the service and repository layers
// to apply selective updates to a service record.
type UpdateServiceCommand struct {
	// ContainerID is the new container ID for the service.
	ContainerID *string
}
