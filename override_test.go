package keyfunc_test

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang-jwt/jwt/v4"

	"github.com/MicahParks/keyfunc"
)

const (
	// TODO
	givenKID  = "givenKID"
	remoteKID = "remoteKID"
)

// TODO
type pseudoJWKs struct {
	Keys []pseudoJSONKey `json:"keys"`
}

// TODO
type pseudoJSONKey struct {
	KID string `json:"kid"`
	E   string `json:"e"`
	N   string `json:"n"`
}

// TestNewGiven tests that given keys will be added to a JWKs with a remote resource.
func TestNewGiven(t *testing.T) {

	// Create a temporary directory to serve the JWKs from.
	tempDir, err := ioutil.TempDir("", "*")
	if err != nil {
		t.Errorf("Failed to create a temporary directory.\nError: %s", err.Error())
		t.FailNow()
	}
	defer func() {
		if err = os.RemoveAll(tempDir); err != nil {
			t.Errorf("Failed to remove temporary directory.\nError: %s", err.Error())
			t.FailNow()
		}
	}()

	// Create the JWKs file path.
	jwksFile := filepath.Join(tempDir, jwksFilePath)

	// Create the keys used for this test.
	givenKeys, givenPrivateKeys, jwksBytes, remotePrivateKeys, err := keysAndJWKs()
	if err != nil {
		t.Errorf("Failed to create cryptographic keys for the test: %s.", err.Error())
		t.FailNow()
	}

	// Write the empty JWKs.
	if err = ioutil.WriteFile(jwksFile, jwksBytes, 0600); err != nil {
		t.Errorf("Failed to write JWKs file to temporary directory.\nError: %s", err.Error())
		t.FailNow()
	}

	// Create the HTTP test server.
	server := httptest.NewServer(http.FileServer(http.FS(os.DirFS(tempDir))))
	defer server.Close()

	// Create testing options.
	testingRefreshErrorHandler := func(err error) {
		panic(fmt.Sprintf("Unhandled JWKs error: %s", err.Error()))
	}

	// Set the JWKs URL.
	jwksURL := server.URL + jwksFilePath

	// Create the test options.
	options := keyfunc.Options{
		GivenKeys:           givenKeys,
		RefreshErrorHandler: testingRefreshErrorHandler,
	}

	// Get the remote JWKs.
	var jwks *keyfunc.JWKs
	if jwks, err = keyfunc.Get(jwksURL, options); err != nil {
		t.Errorf("Failed to get the JWKs the testing URL.\nError: %s", err.Error())
		t.FailNow()
	}

	// Test the given key with a unique key ID.
	createSignParseValidate(t, givenPrivateKeys, jwks, givenKID, true)

	// Test the given key with a non-unique key ID that should be overwritten.
	createSignParseValidate(t, givenPrivateKeys, jwks, remoteKID, false)

	// Test the remote key that should not have been overwritten.
	createSignParseValidate(t, remotePrivateKeys, jwks, remoteKID, true)

	// Change the JWKs options to overwrite remote keys.
	givenKidOverride := true
	options.GivenKIDOverride = &givenKidOverride
	if jwks, err = keyfunc.Get(jwksURL, options); err != nil {
		t.Errorf("Failed to recreate JWKs: %s.", err.Error())
		t.FailNow()
	}

	// Test the given key with a unique key ID.
	createSignParseValidate(t, givenPrivateKeys, jwks, givenKID, true)

	// Test the given key with a non-unique key ID that should overwrite the remote key.
	createSignParseValidate(t, givenPrivateKeys, jwks, remoteKID, true)

	// Test the remote key that should have been overwritten.
	createSignParseValidate(t, remotePrivateKeys, jwks, remoteKID, false)

}

// TODO
func createSignParseValidate(t *testing.T, keys map[string]*rsa.PrivateKey, jwks *keyfunc.JWKs, kid string, shouldValidate bool) {

	// Create the JWT.
	unsignedToken := jwt.New(jwt.SigningMethodRS256)
	unsignedToken.Header[kidAttribute] = kid

	// Sign the JWT.
	tokenString, err := unsignedToken.SignedString(keys[kid])
	if err != nil {
		t.Errorf("Failed to sign the JWT: %s.", err.Error())
		t.FailNow()
	}

	// Parse the JWT.
	var token *jwt.Token
	token, err = jwt.Parse(tokenString, jwks.Keyfunc)
	if err != nil {
		if !shouldValidate && !errors.Is(err, rsa.ErrVerification) {
			return
		}
		t.Errorf("Failed to parse the JWT: %s", err.Error())
		t.FailNow()

	}
	if !shouldValidate {
		t.Errorf("The token should not have parsed properly.")
		t.FailNow()
	}

	// Validate the JWT.
	if !token.Valid {
		t.Errorf("The JWT is not valid.")
		t.FailNow()
	}
}

// TODO
func keysAndJWKs() (givenKeys map[string]keyfunc.GivenKey, givenPrivateKeys map[string]*rsa.PrivateKey, jwksBytes []byte, remotePrivateKeys map[string]*rsa.PrivateKey, err error) {

	// Initialize the function's assets.
	const rsaErrMessage = "failed to create RSA key: %w"
	givenKeys = make(map[string]keyfunc.GivenKey)
	givenPrivateKeys = make(map[string]*rsa.PrivateKey)
	remotePrivateKeys = make(map[string]*rsa.PrivateKey)

	// Create a key not in the remote JWKs.
	var key1 *rsa.PrivateKey
	if key1, err = addRSA(givenKeys, givenKID); err != nil {
		return nil, nil, nil, nil, fmt.Errorf(rsaErrMessage, err)
	}
	givenPrivateKeys[givenKID] = key1

	// Create a key to be overwritten by or override the one with the same key ID in the remote JWKs.
	var key2 *rsa.PrivateKey
	if key2, err = addRSA(givenKeys, remoteKID); err != nil {
		return nil, nil, nil, nil, fmt.Errorf(rsaErrMessage, err)
	}
	givenPrivateKeys[remoteKID] = key2

	// Create a key that exists in the remote JWKs.
	var key3 *rsa.PrivateKey
	if key3, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
		return nil, nil, nil, nil, fmt.Errorf(rsaErrMessage, err)
	}
	remotePrivateKeys[remoteKID] = key3

	// Create a pseudo-JWKs.
	jwks := pseudoJWKs{Keys: []pseudoJSONKey{{
		KID: remoteKID,
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key3.PublicKey.E)).Bytes()),
		N:   base64.RawURLEncoding.EncodeToString(key3.PublicKey.N.Bytes()),
	}}}

	// Marshal the JWKs to JSON.
	if jwksBytes, err = json.Marshal(jwks); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to marshal the JWKs to JSON: %w", err)
	}

	return givenKeys, givenPrivateKeys, jwksBytes, remotePrivateKeys, nil
}