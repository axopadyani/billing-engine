package businesserror

// Kind represents different types of business errors.
type Kind int

const (
	// KindInternal represents an internal server error.
	KindInternal Kind = iota

	// KindBadRequest indicates that the client sent an invalid request.
	KindBadRequest

	// KindUnprocessableEntity signifies that the server understands the content type of the request entity,
	// and the request entity is correct, but it was unable to process the contained instructions.
	KindUnprocessableEntity

	// KindNotFound denotes that the requested resource could not be found.
	KindNotFound

	// KindAlreadyExists indicates that an attempt to create an entity failed because it already exists.
	KindAlreadyExists
)
