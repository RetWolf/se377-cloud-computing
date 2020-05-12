package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func listS3Buckets(svc *s3.S3) {
	result, err := svc.ListBuckets(nil)
	if err != nil {
		fmt.Printf("Unable to list buckets, %v", err)
		os.Exit(1)
	}

	fmt.Println("Buckets:")
	for _, bucket := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(bucket.Name), aws.TimeValue(bucket.CreationDate))
	}
}

func createS3Bucket(svc *s3.S3, bucketName string) string {
	bucket := s3.CreateBucketInput{
		Bucket: &bucketName,
	}
	result, err := svc.CreateBucket(&bucket)
	if err != nil {
		fmt.Printf("Unable to create bucket, %v", err)
		os.Exit(1)
	}

	fmt.Printf("Created bucket at location: %s", *result.Location)

	return *result.Location
}

func main() {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	svc := s3.New(sess)

	bucket := "s3programmatic-access-se377-conway"
	obj := "example.txt"

	// Setup - create bucket, list buckets
	createS3Bucket(svc, bucket)
	listS3Buckets(svc)
	fmt.Println("")

	file, err := os.Open(obj)
	if err != nil {
		fmt.Printf("Unable to open file for upload, %v", err)
		os.Exit(1)
	}

	defer file.Close()

	uploader := s3manager.NewUploader(sess)

	// Upload example file
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(obj),
		Body:   file,
	})
	if err != nil {
		fmt.Printf("Unable to upload file, %v", err)
		os.Exit(1)
	}

	// List files in bucket
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		fmt.Printf("Unable to list items in bucket %q, %v", bucket, err)
		os.Exit(1)
	}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}

	// Delete file
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(obj)})
	if err != nil {
		fmt.Printf("Unable to delete object %q from bucket %q, %v", obj, bucket, err)
		os.Exit(1)
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(obj),
	})

	if err != nil {
		fmt.Printf("Error occurred while waiting for bucket to be deleted, %v", bucket)
	}

	fmt.Printf("Object %q successfully deleted\n", obj)
	fmt.Println("")

	// List files in bucket to demo deletion
	resp, err = svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(bucket)})
	if err != nil {
		fmt.Printf("Unable to list items in bucket %q, %v", bucket, err)
		os.Exit(1)
	}

	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}

	_, err = svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		fmt.Printf("Unable to delete bucket %q, %v", bucket, err)
	}

	// Wait until bucket is deleted before finishing
	fmt.Printf("Waiting for bucket %q to be deleted...\n", bucket)
	fmt.Println("")

	err = svc.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		fmt.Printf("Error occurred while waiting for bucket to be deleted, %v", bucket)
	}

	fmt.Printf("Bucket %q successfully deleted\n", bucket)

}
