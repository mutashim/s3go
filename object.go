package s3go

type FileObject struct {
	Name    string
	LastMod string
	Owner   string
}

type BucketObject struct {
	Name      string
	CreatedAt string
}
