//go:build !go1.6
// +build !go1.6

package s3

import "github.com/IBM/ibm-cos-sdk-go/aws/request"

func platformRequestHandlers(r *request.Request) {
}
