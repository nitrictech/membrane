// Copyright 2021 Nitric Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/durationpb"

	mock_provider "github.com/nitrictech/nitric/cloud/aws/mocks/provider"
	mock_s3iface "github.com/nitrictech/nitric/cloud/aws/mocks/s3"
	"github.com/nitrictech/nitric/cloud/aws/runtime/resource"
	s3_service "github.com/nitrictech/nitric/cloud/aws/runtime/storage"
	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
)

var _ = Describe("S3", func() {
	When("Write", func() {
		When("Given the S3 backend is available", func() {
			When("Creating an object in an existing bucket", func() {
				testPayload := []byte("Test")
				ctrl := gomock.NewController(GinkgoT())

				mockStorageClient := mock_s3iface.NewMockS3API(ctrl)
				mockPSStorageClient := mock_s3iface.NewMockPreSignAPI(ctrl)
				mockProvider := mock_provider.NewMockAwsResourceResolver(ctrl)
				storagePlugin, _ := s3_service.NewWithClient(mockProvider, mockStorageClient, mockPSStorageClient)

				It("Should successfully store the object", func() {
					By("the bucket existing")
					mockProvider.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Bucket).Return(map[string]resource.ResolvedResource{
						"my-bucket": {ARN: "arn:aws:s3:::my-bucket"},
					}, nil)

					By("writing the item")
					mockStorageClient.EXPECT().PutObject(gomock.Any(), gomock.Any()).Return(&s3.PutObjectOutput{}, nil)

					_, err := storagePlugin.Write(context.TODO(), &storagepb.StorageWriteRequest{
						BucketName: "my-bucket",
						Key:        "test-item",
						Body:       testPayload,
					})
					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("Creating an object in a non-existent bucket", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockStorageClient := mock_s3iface.NewMockS3API(ctrl)
				mockPSStorageClient := mock_s3iface.NewMockPreSignAPI(ctrl)
				mockProvider := mock_provider.NewMockAwsResourceResolver(ctrl)
				storagePlugin, _ := s3_service.NewWithClient(mockProvider, mockStorageClient, mockPSStorageClient)
				It("Should fail to store the item", func() {
					By("the bucket not existing")
					mockProvider.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Bucket).Return(map[string]resource.ResolvedResource{}, nil)

					_, err := storagePlugin.Write(context.TODO(), &storagepb.StorageWriteRequest{
						BucketName: "my-bucket",
						Key:        "test-item",
						Body:       []byte("Test"),
					})
					By("Returning an error")
					Expect(err).Should(HaveOccurred())
				})
			})

			When("Creating in a bucket without permissions", func() {
				testPayload := []byte("Test")
				ctrl := gomock.NewController(GinkgoT())

				mockStorageClient := mock_s3iface.NewMockS3API(ctrl)
				mockPSStorageClient := mock_s3iface.NewMockPreSignAPI(ctrl)
				mockProvider := mock_provider.NewMockAwsResourceResolver(ctrl)
				storagePlugin, _ := s3_service.NewWithClient(mockProvider, mockStorageClient, mockPSStorageClient)

				It("Should fail to store the object", func() {
					By("the bucket existing")
					mockProvider.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Bucket).Return(map[string]resource.ResolvedResource{
						"my-bucket": {ARN: "arn:aws:s3:::my-bucket"},
					}, nil)

					opErr := &smithy.OperationError{
						ServiceID: "S3",
						Err:       errors.New("AccessDenied"),
					}

					By("failing to write the item")
					mockStorageClient.EXPECT().PutObject(gomock.Any(), gomock.Any()).Return(nil, opErr)

					_, err := storagePlugin.Write(context.TODO(), &storagepb.StorageWriteRequest{
						BucketName: "my-bucket",
						Key:        "test-item",
						Body:       testPayload,
					})

					By("Returning an error")
					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("unable to write file"))
				})
			})
		})
	})
	When("Read", func() {
		When("The S3 backend is available", func() {
			When("The bucket exists", func() {
				When("The item exists", func() {
					ctrl := gomock.NewController(GinkgoT())
					mockStorageClient := mock_s3iface.NewMockS3API(ctrl)
					mockPSStorageClient := mock_s3iface.NewMockPreSignAPI(ctrl)
					mockProvider := mock_provider.NewMockAwsResourceResolver(ctrl)
					storagePlugin, _ := s3_service.NewWithClient(mockProvider, mockStorageClient, mockPSStorageClient)

					It("Should successfully retrieve the object", func() {
						By("the bucket existing")
						mockProvider.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Bucket).Return(map[string]resource.ResolvedResource{
							"test-bucket": {ARN: "arn:aws:s3:::test-bucket"},
						}, nil)

						By("the object existing")
						mockStorageClient.EXPECT().GetObject(gomock.Any(), &s3.GetObjectInput{
							Bucket: aws.String("test-bucket"),
							Key:    aws.String("test-key"),
						}).Return(&s3.GetObjectOutput{
							Body: io.NopCloser(bytes.NewReader([]byte("Test"))),
						}, nil)

						resp, err := storagePlugin.Read(context.TODO(), &storagepb.StorageReadRequest{
							BucketName: "test-bucket",
							Key:        "test-key",
						})

						By("Not returning an error")
						Expect(err).ShouldNot(HaveOccurred())

						By("Returning the item")
						Expect(resp.Body).To(Equal([]byte("Test")))
					})
				})
				When("The item doesn't exist", func() {
				})

				When("Accessing a file without permissions", func() {
					ctrl := gomock.NewController(GinkgoT())
					mockStorageClient := mock_s3iface.NewMockS3API(ctrl)
					mockPSStorageClient := mock_s3iface.NewMockPreSignAPI(ctrl)
					mockProvider := mock_provider.NewMockAwsResourceResolver(ctrl)
					storagePlugin, _ := s3_service.NewWithClient(mockProvider, mockStorageClient, mockPSStorageClient)

					It("Should fail to retrieve the object", func() {
						By("the bucket existing")
						mockProvider.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Bucket).Return(map[string]resource.ResolvedResource{
							"test-bucket": {ARN: "arn:aws:s3:::test-bucket"},
						}, nil)

						opErr := &smithy.OperationError{
							ServiceID: "S3",
							Err:       errors.New("AccessDenied"),
						}

						By("the object existing")
						mockStorageClient.EXPECT().GetObject(gomock.Any(), &s3.GetObjectInput{
							Bucket: aws.String("test-bucket"),
							Key:    aws.String("test-key"),
						}).Return(nil, opErr)

						resp, err := storagePlugin.Read(context.TODO(), &storagepb.StorageReadRequest{
							BucketName: "test-bucket",
							Key:        "test-key",
						})

						By("Returning an error")
						Expect(err).Should(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("unable to read file"))
						By("returning a nil response")
						Expect(resp).Should(BeNil())
					})
				})
			})
			When("The bucket doesn't exist", func() {
			})
		})
	})
	When("Delete", func() {
		When("The S3 backend is available", func() {
			When("The bucket exists", func() {
				When("The item exists", func() {
					ctrl := gomock.NewController(GinkgoT())
					mockStorageClient := mock_s3iface.NewMockS3API(ctrl)
					mockPSStorageClient := mock_s3iface.NewMockPreSignAPI(ctrl)
					mockProvider := mock_provider.NewMockAwsResourceResolver(ctrl)
					storagePlugin, _ := s3_service.NewWithClient(mockProvider, mockStorageClient, mockPSStorageClient)

					It("Should successfully delete the object", func() {
						By("the bucket existing")
						mockProvider.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Bucket).Return(map[string]resource.ResolvedResource{
							"test-bucket": {ARN: "arn:aws:s3:::test-bucket"},
						}, nil)

						opErr := &smithy.OperationError{
							ServiceID: "S3",
							Err:       errors.New("AccessDenied"),
						}

						By("not deleting the object")
						mockStorageClient.EXPECT().DeleteObject(gomock.Any(), gomock.Any()).Return(nil, opErr)

						_, err := storagePlugin.Delete(context.TODO(), &storagepb.StorageDeleteRequest{
							BucketName: "test-bucket",
							Key:        "test-key",
						})

						By("Returning an error")
						Expect(err).Should(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("unable to delete file"))
					})
				})
				When("The item doesn't exist", func() {
				})

				When("Accessing without permissions", func() {
					ctrl := gomock.NewController(GinkgoT())
					mockStorageClient := mock_s3iface.NewMockS3API(ctrl)
					mockPSStorageClient := mock_s3iface.NewMockPreSignAPI(ctrl)
					mockProvider := mock_provider.NewMockAwsResourceResolver(ctrl)
					storagePlugin, _ := s3_service.NewWithClient(mockProvider, mockStorageClient, mockPSStorageClient)

					It("Should fail to delete the object", func() {
						By("the bucket existing")
						mockProvider.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Bucket).Return(map[string]resource.ResolvedResource{
							"test-bucket": {ARN: "arn:aws:s3:::test-bucket-aaa111"},
						}, nil)

						By("successfully deleting the object")
						mockStorageClient.EXPECT().DeleteObject(gomock.Any(), gomock.Any()).Return(&s3.DeleteObjectOutput{}, nil)

						_, err := storagePlugin.Delete(context.TODO(), &storagepb.StorageDeleteRequest{
							BucketName: "test-bucket",
							Key:        "test-key",
						})
						By("Not returning an error")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
			})
			When("The bucket doesn't exist", func() {
			})
		})
	})
	When("PreSignUrl", func() {
		When("The bucket exists", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockProvider := mock_provider.NewMockAwsResourceResolver(ctrl)
			mockStorageClient := mock_s3iface.NewMockS3API(ctrl)
			mockPSStorageClient := mock_s3iface.NewMockPreSignAPI(ctrl)
			storagePlugin, _ := s3_service.NewWithClient(mockProvider, mockStorageClient, mockPSStorageClient)

			When("A URL is requested for a known operation", func() {
				It("Should successfully generate the URL", func() {
					By("the bucket existing")
					mockProvider.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Bucket).Return(map[string]resource.ResolvedResource{
						"test-bucket": {ARN: "arn:aws:s3:::test-bucket-aaa111"},
					}, nil)

					mockPSStorageClient.EXPECT().PresignPutObject(gomock.Any(), &s3.PutObjectInput{
						Bucket: aws.String("test-bucket-aaa111"), // the real bucket name should be provided here, not the nitric name
						Key:    aws.String("test-key"),
					}, gomock.Any()).Times(1).Return(&v4.PresignedHTTPRequest{
						URL: "aws.example.com",
					}, nil)

					resp, err := storagePlugin.PreSignUrl(context.TODO(), &storagepb.StoragePreSignUrlRequest{
						BucketName: "test-bucket",
						Key:        "test-key",
						Expiry:     durationpb.New(time.Second * 60),
						Operation:  storagepb.StoragePreSignUrlRequest_WRITE,
					})
					Expect(err).ShouldNot(HaveOccurred())

					By("Return the correct url")
					// always blank - it's the best we can do without a real mock.
					Expect(resp.Url).To(Equal("aws.example.com"))
				})
			})
		})
	})

	When("ListFiles", func() {
		When("The bucket exists", func() {
			When("The s3 backend is available", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockStorageClient := mock_s3iface.NewMockS3API(ctrl)
				mockPSStorageClient := mock_s3iface.NewMockPreSignAPI(ctrl)
				mockProvider := mock_provider.NewMockAwsResourceResolver(ctrl)
				storagePlugin, _ := s3_service.NewWithClient(mockProvider, mockStorageClient, mockPSStorageClient)

				It("should list the files contained in the bucket", func() {
					By("the bucket existing")
					mockProvider.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Bucket).Return(map[string]resource.ResolvedResource{
						"test-bucket": {ARN: "arn:aws:s3:::test-bucket-aaa111"},
					}, nil)

					By("s3 returning files")
					mockStorageClient.EXPECT().ListObjectsV2(gomock.Any(), &s3.ListObjectsV2Input{
						Bucket: aws.String("test-bucket-aaa111"),
						Prefix: aws.String("test/"),
					}).Return(&s3.ListObjectsV2Output{
						Contents: []types.Object{{
							Key: aws.String("test/test"),
						}},
					}, nil)

					resp, err := storagePlugin.ListBlobs(context.TODO(), &storagepb.StorageListBlobsRequest{
						BucketName: "test-bucket",
						Prefix:     "test/",
					})

					By("not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("returning the file listing from s3")
					Expect(resp.Blobs).To(HaveLen(1))

					By("having the returned keys")
					Expect(resp.Blobs[0].Key).To(Equal("test/test"))
				})
			})
		})
	})

	When("Exists", func() {
		When("The bucket exists", func() {
			When("The s3 backend is available", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockStorageClient := mock_s3iface.NewMockS3API(ctrl)
				mockPSStorageClient := mock_s3iface.NewMockPreSignAPI(ctrl)
				mockProvider := mock_provider.NewMockAwsResourceResolver(ctrl)
				storagePlugin, _ := s3_service.NewWithClient(mockProvider, mockStorageClient, mockPSStorageClient)

				When("the file exists", func() {
					It("should return true", func() {
						By("the bucket existing")
						mockProvider.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Bucket).Return(map[string]resource.ResolvedResource{
							"test-bucket": {ARN: "arn:aws:s3:::test-bucket-aaa111"},
						}, nil)

						By("the file existing")
						mockStorageClient.EXPECT().HeadObject(gomock.Any(), &s3.HeadObjectInput{
							Bucket: aws.String("test-bucket-aaa111"),
							Key:    aws.String("test-file"),
						}).Return(&s3.HeadObjectOutput{}, nil)

						exists, err := storagePlugin.Exists(context.TODO(), &storagepb.StorageExistsRequest{
							BucketName: "test-bucket",
							Key:        "test-file",
						})

						By("not returning an error")
						Expect(err).ShouldNot(HaveOccurred())

						By("returning false")
						Expect(exists.Exists).To(Equal(true))
					})
				})

				When("the file does not exist", func() {
					It("should return false", func() {
						By("the bucket existing")
						mockProvider.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Bucket).Return(map[string]resource.ResolvedResource{
							"test-bucket": {ARN: "arn:aws:s3:::test-bucket-aaa111"},
						}, nil)

						By("the file not existing")
						mockStorageClient.EXPECT().HeadObject(gomock.Any(), &s3.HeadObjectInput{
							Bucket: aws.String("test-bucket-aaa111"),
							Key:    aws.String("test-file"),
						}).Return(nil, fmt.Errorf("mock-error"))

						exists, err := storagePlugin.Exists(context.TODO(), &storagepb.StorageExistsRequest{
							BucketName: "test-bucket",
							Key:        "test-file",
						})

						By("not returning an error")
						Expect(err).ShouldNot(HaveOccurred())

						By("returning false")
						Expect(exists.Exists).To(Equal(false))
					})
				})
			})
		})
	})
})
