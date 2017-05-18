package ops

import (
	"bytes"
	"io"

	"github.com/sirupsen/logrus"

	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Engine defines an AWS S3 backed object storage engine
type S3Engine struct {
	sess       *session.Session
	downloader *s3manager.Downloader
	bucket     *string
}

// NewS3 creates an
func NewS3(region string, bucket string) *S3Engine {
	config := aws.NewConfig().WithRegion(region).WithS3UseAccelerate(false)
	e := &S3Engine{}
	e.sess = session.New(config)
	e.downloader = s3manager.NewDownloader(e.sess)
	e.bucket = aws.String(bucket)

	return e
}

// WriteTo reads key from S3 and writes the bytes to w
func (e *S3Engine) WriteTo(key string, w io.Writer) error {
	logrus.Debug("excuting S3Engine WriteTo")
	if writerAt, ok := w.(io.WriterAt); ok {
		return e.download(key, writerAt)
	}
	data := make([]byte, bytes.MinRead)
	wab := aws.NewWriteAtBuffer(data)
	err := e.download(key, wab)
	if err != nil {
		return err
	}
	_, err = w.Write(wab.Bytes())
	return err
}

func (e *S3Engine) download(key string, w io.WriterAt) error {
	obj := &s3.GetObjectInput{
		Bucket: e.bucket,
		Key:    aws.String(key),
	}
	numbytes, err := e.downloader.Download(w, obj)
	if err != nil {
		if rf, ok := err.(awserr.RequestFailure); ok {
			if rf.StatusCode() == 404 {
				logrus.WithField("key", key).Info("key does not exist")
				return errors.New(rf.Message())
			}
		}
		logrus.WithField("key", key).Debug("failed to read data from key")
		return err
	}
	logrus.WithFields(logrus.Fields{"key": key, "bytes": numbytes}).Info("read bytes from S3")
	return nil
}

// ReadFrom reads data from r and stores it under key
func (e *S3Engine) ReadFrom(key string, r io.Reader) error {
	return s3upload(e, key, r)
}

// Delete remove the
func (e *S3Engine) Delete(key string) error {
	return nil
}

func s3upload(e *S3Engine, key string, reader io.Reader) error {
	logrus.WithField("bucket", e.bucket).Info("engine configuration")
	uploader := s3manager.NewUploader(e.sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   reader,
		Bucket: e.bucket,
		Key:    aws.String(key),
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{"awserr": err, "key": key}).Error("failed to upload")
		return err
	}
	logrus.WithFields(logrus.Fields{"key": key, "location": result.Location}).Info("uploaded key")
	return nil
}
