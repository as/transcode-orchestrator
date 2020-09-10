// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package mediaconvert

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
	"github.com/aws/aws-sdk-go-v2/private/protocol"
)

// Removes an association between the Amazon Resource Name (ARN) of an AWS Certificate
// Manager (ACM) certificate and an AWS Elemental MediaConvert resource.
type DisassociateCertificateInput struct {
	_ struct{} `type:"structure"`

	// The ARN of the ACM certificate that you want to disassociate from your MediaConvert
	// resource.
	//
	// Arn is a required field
	Arn *string `location:"uri" locationName:"arn" type:"string" required:"true"`
}

// String returns the string representation
func (s DisassociateCertificateInput) String() string {
	return awsutil.Prettify(s)
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DisassociateCertificateInput) Validate() error {
	invalidParams := aws.ErrInvalidParams{Context: "DisassociateCertificateInput"}

	if s.Arn == nil {
		invalidParams.Add(aws.NewErrParamRequired("Arn"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s DisassociateCertificateInput) MarshalFields(e protocol.FieldEncoder) error {
	e.SetValue(protocol.HeaderTarget, "Content-Type", protocol.StringValue("application/json"), protocol.Metadata{})

	if s.Arn != nil {
		v := *s.Arn

		metadata := protocol.Metadata{}
		e.SetValue(protocol.PathTarget, "arn", protocol.QuotedValue{ValueMarshaler: protocol.StringValue(v)}, metadata)
	}
	return nil
}

// Successful disassociation of Certificate Manager Amazon Resource Name (ARN)
// with Mediaconvert returns an OK message.
type DisassociateCertificateOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DisassociateCertificateOutput) String() string {
	return awsutil.Prettify(s)
}

// MarshalFields encodes the AWS API shape using the passed in protocol encoder.
func (s DisassociateCertificateOutput) MarshalFields(e protocol.FieldEncoder) error {
	return nil
}

const opDisassociateCertificate = "DisassociateCertificate"

// DisassociateCertificateRequest returns a request value for making API operation for
// AWS Elemental MediaConvert.
//
// Removes an association between the Amazon Resource Name (ARN) of an AWS Certificate
// Manager (ACM) certificate and an AWS Elemental MediaConvert resource.
//
//    // Example sending a request using DisassociateCertificateRequest.
//    req := client.DisassociateCertificateRequest(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
//
// Please also see https://docs.aws.amazon.com/goto/WebAPI/mediaconvert-2017-08-29/DisassociateCertificate
func (c *Client) DisassociateCertificateRequest(input *DisassociateCertificateInput) DisassociateCertificateRequest {
	op := &aws.Operation{
		Name:       opDisassociateCertificate,
		HTTPMethod: "DELETE",
		HTTPPath:   "/2017-08-29/certificates/{arn}",
	}

	if input == nil {
		input = &DisassociateCertificateInput{}
	}

	req := c.newRequest(op, input, &DisassociateCertificateOutput{})

	return DisassociateCertificateRequest{Request: req, Input: input, Copy: c.DisassociateCertificateRequest}
}

// DisassociateCertificateRequest is the request type for the
// DisassociateCertificate API operation.
type DisassociateCertificateRequest struct {
	*aws.Request
	Input *DisassociateCertificateInput
	Copy  func(*DisassociateCertificateInput) DisassociateCertificateRequest
}

// Send marshals and sends the DisassociateCertificate API request.
func (r DisassociateCertificateRequest) Send(ctx context.Context) (*DisassociateCertificateResponse, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &DisassociateCertificateResponse{
		DisassociateCertificateOutput: r.Request.Data.(*DisassociateCertificateOutput),
		response:                      &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// DisassociateCertificateResponse is the response type for the
// DisassociateCertificate API operation.
type DisassociateCertificateResponse struct {
	*DisassociateCertificateOutput

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// DisassociateCertificate request.
func (r *DisassociateCertificateResponse) SDKResponseMetdata() *aws.Response {
	return r.response
}
