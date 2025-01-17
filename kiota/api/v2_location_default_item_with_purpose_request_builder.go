package api

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588 "go.artefactual.dev/ssclient/kiota/models"
)

// V2LocationDefaultItemWithPurposeRequestBuilder builds and executes requests for operations under \api\v2\location\default\{purpose}\
type V2LocationDefaultItemWithPurposeRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// V2LocationDefaultItemWithPurposeRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type V2LocationDefaultItemWithPurposeRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewV2LocationDefaultItemWithPurposeRequestBuilderInternal instantiates a new V2LocationDefaultItemWithPurposeRequestBuilder and sets the default values.
func NewV2LocationDefaultItemWithPurposeRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter, purpose *string)(*V2LocationDefaultItemWithPurposeRequestBuilder) {
    m := &V2LocationDefaultItemWithPurposeRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/api/v2/location/default/{purpose}/", pathParameters),
    }
    if purpose != nil {
        m.BaseRequestBuilder.PathParameters["purpose"] = *purpose
    }
    return m
}
// NewV2LocationDefaultItemWithPurposeRequestBuilder instantiates a new V2LocationDefaultItemWithPurposeRequestBuilder and sets the default values.
func NewV2LocationDefaultItemWithPurposeRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*V2LocationDefaultItemWithPurposeRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewV2LocationDefaultItemWithPurposeRequestBuilderInternal(urlParams, requestAdapter, nil)
}
// returns a ErrorEscaped error when the service returns a 400 status code
func (m *V2LocationDefaultItemWithPurposeRequestBuilder) Get(ctx context.Context, requestConfiguration *V2LocationDefaultItemWithPurposeRequestBuilderGetRequestConfiguration)(error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "400": ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588.CreateErrorEscapedFromDiscriminatorValue,
    }
    err = m.BaseRequestBuilder.RequestAdapter.SendNoContent(ctx, requestInfo, errorMapping)
    if err != nil {
        return err
    }
    return nil
}
// returns a *RequestInformation when successful
func (m *V2LocationDefaultItemWithPurposeRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *V2LocationDefaultItemWithPurposeRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *V2LocationDefaultItemWithPurposeRequestBuilder when successful
func (m *V2LocationDefaultItemWithPurposeRequestBuilder) WithUrl(rawUrl string)(*V2LocationDefaultItemWithPurposeRequestBuilder) {
    return NewV2LocationDefaultItemWithPurposeRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
