package s3go

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	s3client "github.com/aws/aws-sdk-go/service/s3"
)

type s3 struct {
	client *s3client.S3
	Bucket string
}

func New(config *Config) (S3, error) {
	session, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, ""),
		Endpoint:    aws.String(config.Endpoint),
		Region:      aws.String(config.Region),
	})

	if err != nil {
		return nil, err
	}
	return &s3{
		client: s3client.New(session),
		Bucket: config.BucketName,
	}, nil

}

func (s *s3) SetBucket(name string) {
	s.Bucket = name
}

func (s *s3) MakeBucket(name string) error {
	_, err := s.client.CreateBucket(&s3client.CreateBucketInput{
		Bucket: aws.String(name),
	})

	return err
}

func (s *s3) RemoveBucket(name string) error {
	_, err := s.client.DeleteBucket(&s3client.DeleteBucketInput{
		Bucket: aws.String(name),
	})

	return err
}

func (s *s3) Upload(input string, output string, aclMode acl) error {

	// open file
	file, err := os.Open(input)
	if err != nil {
		return err
	}
	defer file.Close()

	// get file info and make buffer
	fileInfo, _ := os.Stat(input)
	size := fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// upload
	_, err = s.client.PutObject(&s3client.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(output),
		Body:   bytes.NewReader(buffer),
		ACL:    aws.String(aclMode.String()),
	})

	return err
}

func (s *s3) Download(input string, output string) error {

	// get object data from s3
	result, err := s.client.GetObject(&s3client.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(input),
	})
	if err != nil {
		log.Fatal(err)
	}

	// create file
	file, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// write file
	_, err = io.Copy(file, result.Body)

	return err
}

func (s *s3) Share(input string, duration int64) (string, error) {

	result, _ := s.client.GetObjectRequest(&s3client.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(input),
	})

	urlStr, err := result.Presign(time.Duration(duration) * time.Minute)
	if err != nil {
		fmt.Println(err.Error())
	}

	return urlStr, err
}

func (s *s3) Delete(input string) error {

	_, err := s.client.DeleteObject(&s3client.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(input),
	})

	return err
}

func (s *s3) ListFile(input string) (*[]FileObject, int, error) {

	result, err := s.client.ListObjects(&s3client.ListObjectsInput{
		Bucket: aws.String(s.Bucket),
		Prefix: aws.String(input),
	})
	if err != nil {
		return nil, 0, err
	}

	res := []FileObject{}
	space := 1

	for _, v := range result.Contents {

		if len(*v.Key) > space {
			space = len(*v.Key) + 2
		}

		res = append(res, FileObject{
			Name:    *v.Key,
			LastMod: v.LastModified.String(),
			Owner:   *v.Owner.DisplayName,
		})
	}

	return &res, space, nil

}

func (s *s3) ListBucket() (*[]BucketObject, int, error) {
	result, err := s.client.ListBuckets(nil)
	if err != nil {
		return nil, 0, err
	}

	data := []BucketObject{}
	space := 1

	for _, v := range result.Buckets {
		if len(*v.Name) > space {
			space = len(*v.Name) + 2
		}
		data = append(data, BucketObject{
			Name:      *v.Name,
			CreatedAt: v.CreationDate.String(),
		})
	}

	return &data, space, nil
}

func (s *s3) SetAcl(input string, aclMode acl) error {

	_, err := s.client.PutObjectAcl(&s3client.PutObjectAclInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(input),
		ACL:    aws.String(aclMode.String()),
	})

	return err
}
