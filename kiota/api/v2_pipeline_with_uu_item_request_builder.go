package api

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588 "go.artefactual.dev/ssclient/kiota/models"
)

// V2PipelineWithUuItemRequestBuilder builds and executes requests for operations under \api\v2\pipeline\{uuid}
type V2PipelineWithUuItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// V2PipelineWithUuItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type V2PipelineWithUuItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewV2PipelineWithUuItemRequestBuilderInternal instantiates a new V2PipelineWithUuItemRequestBuilder and sets the default values.
func NewV2PipelineWithUuItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*V2PipelineWithUuItemRequestBuilder) {
    m := &V2PipelineWithUuItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/api/v2/pipeline/{uuid}", pathParameters),
    }
    return m
}
// NewV2PipelineWithUuItemRequestBuilder instantiates a new V2PipelineWithUuItemRequestBuilder and sets the default values.
func NewV2PipelineWithUuItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*V2PipelineWithUuItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewV2PipelineWithUuItemRequestBuilderInternal(urlParams, requestAdapter)
}
// returns a Pipelineable when successful
// returns a ErrorEscaped error when the service returns a 400 status code
func (m *V2PipelineWithUuItemRequestBuilder) Get(ctx context.Context, requestConfiguration *V2PipelineWithUuItemRequestBuilderGetRequestConfiguration)(ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588.Pipelineable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "400": ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588.CreateErrorEscapedFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588.CreatePipelineFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588.Pipelineable), nil
}
// returns a *RequestInformation when successful
func (m *V2PipelineWithUuItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *V2PipelineWithUuItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *V2PipelineWithUuItemRequestBuilder when successful
func (m *V2PipelineWithUuItemRequestBuilder) WithUrl(rawUrl string)(*V2PipelineWithUuItemRequestBuilder) {
    return NewV2PipelineWithUuItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
