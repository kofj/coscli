package util

// FindBucket searches for a bucket by name or alias in the given config and returns it along with its index.
// If not found, it returns a temporary bucket with the specified name and -1 as the index.
func FindBucket(config *Config, bucketName string) (Bucket, int, error) {
	for i, b := range config.Buckets {
		if b.Alias == bucketName {
			return b, i, nil
		}
	}
	for i, b := range config.Buckets {
		if b.Name == bucketName {
			return b, i, nil
		}
	}
	var tmpBucket Bucket
	tmpBucket.Name = bucketName
	return tmpBucket, -1, nil
	// return Bucket{}, -1, errors.New("Bucket not exist! Use \"./coscli config show\" to check config file please! ")
}
