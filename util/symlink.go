package util

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
)

// CreateSymlink creates a symbolic link in the specified COS bucket.
// c: *cos.Client - the COS client to use for the operation.
// cosUrl: *StorageUrl - the COS URL of the bucket where the link should be created.
// linkKey: string - the name of the symbolic link to be created.
func CreateSymlink(c *cos.Client, cosUrl StorageUrl, linkKey string) error {
	opt := &cos.ObjectPutSymlinkOptions{
		SymlinkTarget: cosUrl.(*CosUrl).Object,
	}
	_, err := c.Object.PutSymlink(context.Background(), linkKey, opt)
	return err
}

// GetSymlink retrieves a symlink from the COS client.
// Parameters:
// - c: *cos.Client, the COS client instance.
// - linkKey: string, the key of the symlink to retrieve.
// Returns:
// - res: string, the content of the symlink.
// - err: error, if any error occurs during retrieval.
func GetSymlink(c *cos.Client, linkKey string) (res string, err error) {
	res, _, err = c.Object.GetSymlink(context.Background(), linkKey, nil)
	return res, err
}
