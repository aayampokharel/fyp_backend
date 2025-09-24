package errorz

var (
	Status400BadRequest = Config{
		StatusCode: 400,
		Message:    "Bad Request: The server cannot process the request due to client error.",
	}

	Status401Unauthorized = Config{
		StatusCode: 401,
		Message:    "Unauthorized: Authentication is required to access this resource.",
	}

	Status403Forbidden = Config{
		StatusCode: 403,
		Message:    "Forbidden: You do not have permission to access this resource.",
	}

	Status404NotFound = Config{
		StatusCode: 404,
		Message:    "Not Found: The requested resource could not be found.",
	}

	Status405MethodNotAllowed = Config{
		StatusCode: 405,
		Message:    "Method Not Allowed: This HTTP method is not supported for this resource.",
	}

	Status406NotAcceptable = Config{
		StatusCode: 406,
		Message:    "Not Acceptable: Server cannot produce a response matching the criteria.",
	}

	Status407ProxyAuthRequired = Config{
		StatusCode: 407,
		Message:    "Proxy Authentication Required: Authenticate with the proxy first.",
	}

	Status408RequestTimeout = Config{
		StatusCode: 408,
		Message:    "Request Timeout: The server timed out waiting for the request.",
	}

	Status409Conflict = Config{
		StatusCode: 409,
		Message:    "Conflict: The request conflicts with the current state of the server.",
	}

	Status410Gone = Config{
		StatusCode: 410,
		Message:    "Gone: The requested resource has been permanently removed.",
	}

	Status411LengthRequired = Config{
		StatusCode: 411,
		Message:    "Length Required: The server needs a Content-Length header.",
	}

	Status412PreconditionFailed = Config{
		StatusCode: 412,
		Message:    "Precondition Failed: Server does not meet the request’s preconditions.",
	}

	Status413ContentTooLarge = Config{
		StatusCode: 413,
		Message:    "Content Too Large: The request body exceeds the server’s limit.",
	}
)
