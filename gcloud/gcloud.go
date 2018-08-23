/*
DESC: Uploads file to a Google Cloud stroage bucket
Author: Joshua Haupt
Last Modified: 08-23-2018
Modified from source: https://github.com/GoogleCloudPlatform/golang-samples/blob/master/storage/gcsupload/gcsupload.go
*/

package gcloud

import (
  "fmt"
  "log"
  "io"
  "os"
  // Imports the Google Cloud Storage client package.
  "cloud.google.com/go/storage"
  "golang.org/x/net/context"
)


/*
DESC: Driver to upload file to GCloud storage
IN: GCloud project ID as projectID; GCloud bucket name as bucket; file path as source; whether the file is public as public
OUT: Error
*/
func GCUpload(projectID, bucket, source, name string, public bool) error {

  // file reader
  var reader io.Reader
  file, err := os.Open(source)
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()
  reader = file

  ctx := context.Background()
  _, objAttrs, err := upload(ctx, reader, projectID, bucket, name, true)
  if err != nil {
    switch err {
    case storage.ErrBucketNotExist:
      log.Fatal("Bucket Does Not Exist")
    default:
      log.Fatal(err)
    }
  }

  fmt.Printf("File uploaded to https://%s/%s\n", objAttrs.Bucket, objAttrs.Name)

  return err

}


/*
DESC: Uploads a file to GCloud storage
IN: GCloud context as ctx; file reader as reader; GCloud project ID as projectID; GCloud bucket name as bucket; file path as source; whether the file is public as public
OUT: Error
*/
func upload(ctx context.Context, reader io.Reader, projectID, bucket, name string, public bool) (*storage.ObjectHandle, *storage.ObjectAttrs, error) {

  client, err := storage.NewClient(ctx)
  if err != nil {
    return nil, nil, err
  }

  bh := client.Bucket(bucket)
  // Check if bucket exists
  if _, err = bh.Attrs(ctx); err != nil {
    return nil, nil, err
  }

  obj := bh.Object(name)
  w := obj.NewWriter(ctx)
  if _, err := io.Copy(w, reader); err != nil {
    return nil, nil, err
  }

  if err := w.Close(); err != nil {
    return nil, nil, err
  }

  if public {
    if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
      return nil, nil, err
    }
  }

  attrs, err := obj.Attrs(ctx)
  return obj, attrs, err

}
