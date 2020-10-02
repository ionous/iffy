// Package composer runs a local web server to host the iffy story editor.
//
// Http endpoints:
//  - /compose: the composer web app
//  - /stories/<path>/<to>/<file>.if: get or put individual story files
//  - /stories/<path>/<to>/<file>.if/check: post to test story tests
package composer
