package api

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i2b7a3625368152c59661ed1a63c26960f7e9cda05d0fbc8e5d79ac57ca250e0a "go.artefactual.dev/ssclient/kiota/models"
)

// V2LocationDefaultWithPurposeItemRequestBuilder builds and executes requests for operations under \api\v2\location\default\{purpose}
type V2LocationDefaultWithPurposeItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// V2LocationDefaultWithPurposeItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type V2LocationDefaultWithPurposeItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewV2LocationDefaultWithPurposeItemRequestBuilderInternal instantiates a new V2LocationDefaultWithPurposeItemRequestBuilder and sets the default values.
func NewV2LocationDefaultWithPurposeItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*V2LocationDefaultWithPurposeItemRequestBuilder) {
    m := &V2LocationDefaultWithPurposeItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/api/v2/location/default/{purpose}", pathParameters),
    }
    return m
}
// NewV2LocationDefaultWithPurposeItemRequestBuilder instantiates a new V2LocationDefaultWithPurposeItemRequestBuilder and sets the default values.
func NewV2LocationDefaultWithPurposeItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*V2LocationDefaultWithPurposeItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewV2LocationDefaultWithPurposeItemRequestBuilderInternal(urlParams, requestAdapter)
}
// returns a Error error when the service returns a 400 status code
func (m *V2LocationDefaultWithPurposeItemRequestBuilder) Get(ctx context.Context, requestConfiguration *V2LocationDefaultWithPurposeItemRequestBuilderGetRequestConfiguration)(error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "400": i2b7a3625368152c59661ed1a63c26960f7e9cda05d0fbc8e5d79ac57ca250e0a.CreateErrorFromDiscriminatorValue,
    }
    err = m.BaseRequestBuilder.RequestAdapter.SendNoContent(ctx, requestInfo, errorMapping)
    if err != nil {
        return err
    }
    return nil
}
// returns a *RequestInformation when successful
func (m *V2LocationDefaultWithPurposeItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *V2LocationDefaultWithPurposeItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *V2LocationDefaultWithPurposeItemRequestBuilder when successful
func (m *V2LocationDefaultWithPurposeItemRequestBuilder) WithUrl(rawUrl string)(*V2LocationDefaultWithPurposeItemRequestBuilder) {
    return NewV2LocationDefaultWithPurposeItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
