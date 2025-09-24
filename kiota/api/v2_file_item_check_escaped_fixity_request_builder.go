package api

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588 "go.artefactual.dev/ssclient/kiota/models"
)

// V2FileItemCheck_fixityRequestBuilder builds and executes requests for operations under \api\v2\file\{uuid}\check_fixity\
type V2FileItemCheck_fixityRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// V2FileItemCheck_fixityRequestBuilderGetQueryParameters run a fixity check for the specified package.
type V2FileItemCheck_fixityRequestBuilderGetQueryParameters struct {
    Force_local *bool `uriparametername:"force_local"`
}
// V2FileItemCheck_fixityRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type V2FileItemCheck_fixityRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
    // Request query parameters
    QueryParameters *V2FileItemCheck_fixityRequestBuilderGetQueryParameters
}
// NewV2FileItemCheck_fixityRequestBuilderInternal instantiates a new V2FileItemCheck_fixityRequestBuilder and sets the default values.
func NewV2FileItemCheck_fixityRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*V2FileItemCheck_fixityRequestBuilder) {
    m := &V2FileItemCheck_fixityRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/api/v2/file/{uuid}/check_fixity/{?force_local}", pathParameters),
    }
    return m
}
// NewV2FileItemCheck_fixityRequestBuilder instantiates a new V2FileItemCheck_fixityRequestBuilder and sets the default values.
func NewV2FileItemCheck_fixityRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*V2FileItemCheck_fixityRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewV2FileItemCheck_fixityRequestBuilderInternal(urlParams, requestAdapter)
}
// Get run a fixity check for the specified package.
// returns a FixityResponseable when successful
// returns a ErrorEscaped error when the service returns a 400 status code
func (m *V2FileItemCheck_fixityRequestBuilder) Get(ctx context.Context, requestConfiguration *V2FileItemCheck_fixityRequestBuilderGetRequestConfiguration)(ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588.FixityResponseable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "400": ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588.CreateErrorEscapedFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588.CreateFixityResponseFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588.FixityResponseable), nil
}
// ToGetRequestInformation run a fixity check for the specified package.
// returns a *RequestInformation when successful
func (m *V2FileItemCheck_fixityRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *V2FileItemCheck_fixityRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        if requestConfiguration.QueryParameters != nil {
            requestInfo.AddQueryParameters(*(requestConfiguration.QueryParameters))
        }
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *V2FileItemCheck_fixityRequestBuilder when successful
func (m *V2FileItemCheck_fixityRequestBuilder) WithUrl(rawUrl string)(*V2FileItemCheck_fixityRequestBuilder) {
    return NewV2FileItemCheck_fixityRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
