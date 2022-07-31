package s3go

type S3 interface {
	Upload(input string, output string, aclMode acl) error
	Download(input string, output string) error
	Share(input string, duration int64) (string, error)
	Delete(input string) error
	ListFile(input string) (*[]FileObject, int, error)
	ListBucket() (*[]BucketObject, int, error)
	MakeBucket(name string) error
	RemoveBucket(name string) error
	SetAcl(input string, aclMode acl) error
	SetBucket(name string)
}
