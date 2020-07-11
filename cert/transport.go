package cert

import (
	"encoding/json"

	"github.com/cloudflare/cfssl/api"
	"github.com/cloudflare/cfssl/errors"
	"github.com/rs/zerolog/log"
)

// Utility functions for interacting with the transport package.

// authErrorHTTP builds the error code returned for an error.
var authErrorHTTP = int(errors.APIClientError) + int(errors.ClientHTTPError)

const tokenAuthErrorMessage = "invalid token"
const otherAuthErrorMessage = "not authorised"

// isAuthError returns true if the error is due to a CFSSL
// authentication error.
func isAuthError(err error) bool {
	cferr, ok := err.(*errors.Error)
	if !ok {
		return false
	}

	if cferr.ErrorCode == authErrorHTTP {
		var response api.Response
		innerErr := json.Unmarshal([]byte(cferr.Message), &response)
		if innerErr != nil {
			return false
		}
		log.Debug().Err(cferr).Msg("cfssl error received")
		for _, responseError := range response.Errors {
			if (responseError.Message == tokenAuthErrorMessage) || (responseError.Message == otherAuthErrorMessage) {
				return true
			}
		}
	}

	return false
}
