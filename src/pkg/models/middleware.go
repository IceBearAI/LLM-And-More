package models

type contextKey string

const (
	// contextKeyModelName is the context key that holds the model name.
	contextKeyModelName contextKey = "ctx-model-name"
	// contextKeyModelId is the context key that holds the model id.
	contextKeyModelId contextKey = "ctx-model-id"
	// contextKeyModelVersion is the context key that holds the model version.
	contextKeyModelVersion contextKey = "ctx-model-version"
	// contextKeyModelContainerName is the context key that holds the model container name.
	contextKeyModelContainerName contextKey = "ctx-model-container-name"
)
