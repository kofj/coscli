package util

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
)

// GetBucketVersioning retrieves versioning information for a COS bucket.
// Parameters:
// - c: *cos.Client, the COS client to use.
// Returns:
// - res: *cos.BucketGetVersionResult, the versioning information of the bucket.
// - resp: *cos.Response, the response object containing the result.
// - err: error, if an error occurs during the request.
func GetBucketVersioning(c *cos.Client) (res *cos.BucketGetVersionResult, resp *cos.Response, err error) {
	res, resp, err = c.Bucket.GetVersioning(context.Background())
	if err != nil {
		return nil, nil, err
	}
	return res, resp, err
}

// PutBucketVersioning sets the versioning configuration for a bucket.
//
// Parameters:
// - c: *cos.Client - The COS client.
// - status: string - The status to set for versioning.
//
// Returns:
// - resp: *cos.Response - The response from the server.
// - err: error - An error occurred if the operation failed.
func PutBucketVersioning(c *cos.Client, status string) (resp *cos.Response, err error) {
	opt := &cos.BucketPutVersionOptions{
		Status: status,
	}
	resp, err = c.Bucket.PutVersioning(context.Background(), opt)
	return resp, err
}
