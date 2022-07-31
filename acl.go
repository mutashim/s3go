/* More info about AWS ACL:
   https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html#canned-acl
*/

package s3go

type acl string

const (
	PRIVATE acl = "private"
	PUBLIC  acl = "public-read"
)

func (a acl) String() string {
	return string(a)
}
