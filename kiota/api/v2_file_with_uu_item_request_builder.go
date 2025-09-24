package api

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588 "go.artefactual.dev/ssclient/kiota/models"
)

// V2FileWithUuItemRequestBuilder builds and executes requests for operations under \api\v2\file\{uuid}
type V2FileWithUuItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// V2FileWithUuItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type V2FileWithUuItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// Check_fixity the check_fixity property
// returns a *V2FileItemCheck_fixityRequestBuilder when successful
func (m *V2FileWithUuItemRequestBuilder) Check_fixity()(*V2FileItemCheck_fixityRequestBuilder) {
    return NewV2FileItemCheck_fixityRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// NewV2FileWithUuItemRequestBuilderInternal instantiates a new V2FileWithUuItemRequestBuilder and sets the default values.
func NewV2FileWithUuItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*V2FileWithUuItemRequestBuilder) {
    m := &V2FileWithUuItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/api/v2/file/{uuid}", pathParameters),
    }
    return m
}
// NewV2FileWithUuItemRequestBuilder instantiates a new V2FileWithUuItemRequestBuilder and sets the default values.
func NewV2FileWithUuItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*V2FileWithUuItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewV2FileWithUuItemRequestBuilderInternal(urlParams, requestAdapter)
}
// returns a PackageEscapedable when successful
// returns a ErrorEscaped error when the service returns a 400 status code
func (m *V2FileWithUuItemRequestBuilder) Get(ctx context.Context, requestConfiguration *V2FileWithUuItemRequestBuilderGetRequestConfiguration)(ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588.PackageEscapedable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "400": ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588.CreateErrorEscapedFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588.CreatePackageEscapedFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588.PackageEscapedable), nil
}
// Move the move property
// returns a *V2FileItemMoveRequestBuilder when successful
func (m *V2FileWithUuItemRequestBuilder) Move()(*V2FileItemMoveRequestBuilder) {
    return NewV2FileItemMoveRequestBuilderInternal(m.BaseRequestBuilder.PathParameters, m.BaseRequestBuilder.RequestAdapter)
}
// returns a *RequestInformation when successful
func (m *V2FileWithUuItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *V2FileWithUuItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *V2FileWithUuItemRequestBuilder when successful
func (m *V2FileWithUuItemRequestBuilder) WithUrl(rawUrl string)(*V2FileWithUuItemRequestBuilder) {
    return NewV2FileWithUuItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
