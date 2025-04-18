// Code generated by smithy-go-codegen DO NOT EDIT.

package s3

import (
	"context"
	"fmt"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	s3cust "github.com/aws/aws-sdk-go-v2/service/s3/internal/customizations"
	"github.com/aws/smithy-go/middleware"
	"github.com/aws/smithy-go/ptr"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// This operation is not supported for directory buckets.
//
// Deletes the S3 Intelligent-Tiering configuration from the specified bucket.
//
// The S3 Intelligent-Tiering storage class is designed to optimize storage costs
// by automatically moving data to the most cost-effective storage access tier,
// without performance impact or operational overhead. S3 Intelligent-Tiering
// delivers automatic cost savings in three low latency and high throughput access
// tiers. To get the lowest storage cost on data that can be accessed in minutes to
// hours, you can choose to activate additional archiving capabilities.
//
// The S3 Intelligent-Tiering storage class is the ideal storage class for data
// with unknown, changing, or unpredictable access patterns, independent of object
// size or retention period. If the size of an object is less than 128 KB, it is
// not monitored and not eligible for auto-tiering. Smaller objects can be stored,
// but they are always charged at the Frequent Access tier rates in the S3
// Intelligent-Tiering storage class.
//
// For more information, see [Storage class for automatically optimizing frequently and infrequently accessed objects].
//
// Operations related to DeleteBucketIntelligentTieringConfiguration include:
//
// [GetBucketIntelligentTieringConfiguration]
//
// [PutBucketIntelligentTieringConfiguration]
//
// [ListBucketIntelligentTieringConfigurations]
//
// [ListBucketIntelligentTieringConfigurations]: https://docs.aws.amazon.com/AmazonS3/latest/API/API_ListBucketIntelligentTieringConfigurations.html
// [GetBucketIntelligentTieringConfiguration]: https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketIntelligentTieringConfiguration.html
// [PutBucketIntelligentTieringConfiguration]: https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketIntelligentTieringConfiguration.html
// [Storage class for automatically optimizing frequently and infrequently accessed objects]: https://docs.aws.amazon.com/AmazonS3/latest/dev/storage-class-intro.html#sc-dynamic-data-access
func (c *Client) DeleteBucketIntelligentTieringConfiguration(ctx context.Context, params *DeleteBucketIntelligentTieringConfigurationInput, optFns ...func(*Options)) (*DeleteBucketIntelligentTieringConfigurationOutput, error) {
	if params == nil {
		params = &DeleteBucketIntelligentTieringConfigurationInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "DeleteBucketIntelligentTieringConfiguration", params, optFns, c.addOperationDeleteBucketIntelligentTieringConfigurationMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*DeleteBucketIntelligentTieringConfigurationOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type DeleteBucketIntelligentTieringConfigurationInput struct {

	// The name of the Amazon S3 bucket whose configuration you want to modify or
	// retrieve.
	//
	// This member is required.
	Bucket *string

	// The ID used to identify the S3 Intelligent-Tiering configuration.
	//
	// This member is required.
	Id *string

	noSmithyDocumentSerde
}

func (in *DeleteBucketIntelligentTieringConfigurationInput) bindEndpointParams(p *EndpointParameters) {

	p.Bucket = in.Bucket
	p.UseS3ExpressControlEndpoint = ptr.Bool(true)
}

type DeleteBucketIntelligentTieringConfigurationOutput struct {
	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationDeleteBucketIntelligentTieringConfigurationMiddlewares(stack *middleware.Stack, options Options) (err error) {
	if err := stack.Serialize.Add(&setOperationInputMiddleware{}, middleware.After); err != nil {
		return err
	}
	err = stack.Serialize.Add(&awsRestxml_serializeOpDeleteBucketIntelligentTieringConfiguration{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestxml_deserializeOpDeleteBucketIntelligentTieringConfiguration{}, middleware.After)
	if err != nil {
		return err
	}
	if err := addProtocolFinalizerMiddlewares(stack, options, "DeleteBucketIntelligentTieringConfiguration"); err != nil {
		return fmt.Errorf("add protocol finalizers: %v", err)
	}

	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addClientRequestID(stack); err != nil {
		return err
	}
	if err = addComputeContentLength(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addComputePayloadSHA256(stack); err != nil {
		return err
	}
	if err = addRetry(stack, options); err != nil {
		return err
	}
	if err = addRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = addRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addSpanRetryLoop(stack, options); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addSetLegacyContextSigningOptionsMiddleware(stack); err != nil {
		return err
	}
	if err = addPutBucketContextMiddleware(stack); err != nil {
		return err
	}
	if err = addTimeOffsetBuild(stack, c); err != nil {
		return err
	}
	if err = addUserAgentRetryMode(stack, options); err != nil {
		return err
	}
	if err = addIsExpressUserAgent(stack); err != nil {
		return err
	}
	if err = addCredentialSource(stack, options); err != nil {
		return err
	}
	if err = addOpDeleteBucketIntelligentTieringConfigurationValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opDeleteBucketIntelligentTieringConfiguration(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addMetadataRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addRecursionDetection(stack); err != nil {
		return err
	}
	if err = addDeleteBucketIntelligentTieringConfigurationUpdateEndpoint(stack, options); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = v4.AddContentSHA256HeaderMiddleware(stack); err != nil {
		return err
	}
	if err = disableAcceptEncodingGzip(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	if err = addSerializeImmutableHostnameBucketMiddleware(stack, options); err != nil {
		return err
	}
	if err = addSpanInitializeStart(stack); err != nil {
		return err
	}
	if err = addSpanInitializeEnd(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestStart(stack); err != nil {
		return err
	}
	if err = addSpanBuildRequestEnd(stack); err != nil {
		return err
	}
	return nil
}

func (v *DeleteBucketIntelligentTieringConfigurationInput) bucket() (string, bool) {
	if v.Bucket == nil {
		return "", false
	}
	return *v.Bucket, true
}

func newServiceMetadataMiddleware_opDeleteBucketIntelligentTieringConfiguration(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		OperationName: "DeleteBucketIntelligentTieringConfiguration",
	}
}

// getDeleteBucketIntelligentTieringConfigurationBucketMember returns a pointer to
// string denoting a provided bucket member valueand a boolean indicating if the
// input has a modeled bucket name,
func getDeleteBucketIntelligentTieringConfigurationBucketMember(input interface{}) (*string, bool) {
	in := input.(*DeleteBucketIntelligentTieringConfigurationInput)
	if in.Bucket == nil {
		return nil, false
	}
	return in.Bucket, true
}
func addDeleteBucketIntelligentTieringConfigurationUpdateEndpoint(stack *middleware.Stack, options Options) error {
	return s3cust.UpdateEndpoint(stack, s3cust.UpdateEndpointOptions{
		Accessor: s3cust.UpdateEndpointParameterAccessor{
			GetBucketFromInput: getDeleteBucketIntelligentTieringConfigurationBucketMember,
		},
		UsePathStyle:                   options.UsePathStyle,
		UseAccelerate:                  options.UseAccelerate,
		SupportsAccelerate:             true,
		TargetS3ObjectLambda:           false,
		EndpointResolver:               options.EndpointResolver,
		EndpointResolverOptions:        options.EndpointOptions,
		UseARNRegion:                   options.UseARNRegion,
		DisableMultiRegionAccessPoints: options.DisableMultiRegionAccessPoints,
	})
}
