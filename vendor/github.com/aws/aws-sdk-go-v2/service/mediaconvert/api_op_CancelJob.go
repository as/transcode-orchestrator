// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package mediaconvert

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
	"github.com/aws/aws-sdk-go-v2/private/protocol"
)

// Cancel a job by sending a request with the job ID
type CancelJobInput struct {
	_ struct{} `type:"structure"`

	// The Job ID of the job to be cancelled.
	//
	// Id is a required field
	Id *string `location:"uri" locationName:"id" type:"string" required:"true"`
}

// String returns the string representation
func (s CancelJobInput) String() string {
	return awsutil.Prettify(s)
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *CancelJobInput) Validate() error {
	invalidParams := aws.ErrInvalidParams{Context: "CancelJobInput"}

	if s.Id == nil {
		invalidParams.Add(aws.NewErrParamRequired("Id"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s CancelJobInput) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.Id != nil {
		v := *s.Id

		metadata := protocol.Metadata{}
		e.SetValue(protocol.PathTarget, "id", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	return nil
}

// A cancel job request will receive a response with an empty body.
type CancelJobOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s CancelJobOutput) String() string {
	return awsutil.Prettify(s)
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s CancelJobOutput) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opCancelJob = "CancelJob"

// CancelJobRequest returns a request value for making API operation for
// AWS Elemental MediaConvert.
//
// Permanently cancel a job. Once you have canceled a job, you can't start it
// again.
//
//    // Example sending a request using CancelJobRequest.
//    req := client.CancelJobRequest(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
//
// Please also see https://docs.aws.amazon.com/goto/WebAPI/mediaconvert-2017-08-29/CancelJob
func (c *Client) CancelJobRequest(input *CancelJobInput) CancelJobRequest {
	op := &aws.Operation{
		Name:       opCancelJob,
		HTTPMethod: "DELETE",
		HTTPPath:   "/2017-08-29/jobs/{id}",
	}

	if input == nil {
		input = &CancelJobInput{}
	}

	req := c.newRequest(op, input, &CancelJobOutput{})
	return CancelJobRequest{Request: req, Input: input, Copy: c.CancelJobRequest}
}

// CancelJobRequest is the request type for the
// CancelJob API operation.
type CancelJobRequest struct {
	*aws.Request
	Input *CancelJobInput
	Copy  func(*CancelJobInput) CancelJobRequest
}

// Send marshals and sends the CancelJob API request.
func (r CancelJobRequest) Send(ctx context.Context) (*CancelJobResponse, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &CancelJobResponse{
		CancelJobOutput: r.Request.Data.(*CancelJobOutput),
		response:        &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// CancelJobResponse is the response type for the
// CancelJob API operation.
type CancelJobResponse struct {
	*CancelJobOutput

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// CancelJob request.
func (r *CancelJobResponse) SDKResponseMetdata() *aws.Response {
	return r.response
}