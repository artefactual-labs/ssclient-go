package api

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588 "go.artefactual.dev/ssclient/kiota/models"
)

// V2LocationWithUuItemRequestBuilder builds and executes requests for operations under \api\v2\location\{uuid}
type V2LocationWithUuItemRequestBuilder struct {
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.BaseRequestBuilder
}
// V2LocationWithUuItemRequestBuilderGetRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type V2LocationWithUuItemRequestBuilderGetRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// V2LocationWithUuItemRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type V2LocationWithUuItemRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewV2LocationWithUuItemRequestBuilderInternal instantiates a new V2LocationWithUuItemRequestBuilder and sets the default values.
func NewV2LocationWithUuItemRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*V2LocationWithUuItemRequestBuilder) {
    m := &V2LocationWithUuItemRequestBuilder{
        BaseRequestBuilder: *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewBaseRequestBuilder(requestAdapter, "{+baseurl}/api/v2/location/{uuid}", pathParameters),
    }
    return m
}
// NewV2LocationWithUuItemRequestBuilder instantiates a new V2LocationWithUuItemRequestBuilder and sets the default values.
func NewV2LocationWithUuItemRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*V2LocationWithUuItemRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewV2LocationWithUuItemRequestBuilderInternal(urlParams, requestAdapter)
}
// returns a Locationable when successful
// returns a ErrorEscaped error when the service returns a 400 status code
func (m *V2LocationWithUuItemRequestBuilder) Get(ctx context.Context, requestConfiguration *V2LocationWithUuItemRequestBuilderGetRequestConfiguration)(ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588.Locationable, error) {
    requestInfo, err := m.ToGetRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "400": ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588.CreateErrorEscapedFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.Send(ctx, requestInfo, ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588.CreateLocationFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588.Locationable), nil
}
// Post move files to the specified location.
// returns a []byte when successful
// returns a ErrorEscaped error when the service returns a 400 status code
func (m *V2LocationWithUuItemRequestBuilder) Post(ctx context.Context, body ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588.MoveRequestable, requestConfiguration *V2LocationWithUuItemRequestBuilderPostRequestConfiguration)([]byte, error) {
    requestInfo, err := m.ToPostRequestInformation(ctx, body, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "400": ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588.CreateErrorEscapedFromDiscriminatorValue,
    }
    res, err := m.BaseRequestBuilder.RequestAdapter.SendPrimitive(ctx, requestInfo, "[]byte", errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.([]byte), nil
}
// returns a *RequestInformation when successful
func (m *V2LocationWithUuItemRequestBuilder) ToGetRequestInformation(ctx context.Context, requestConfiguration *V2LocationWithUuItemRequestBuilderGetRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.GET, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    return requestInfo, nil
}
// ToPostRequestInformation move files to the specified location.
// returns a *RequestInformation when successful
func (m *V2LocationWithUuItemRequestBuilder) ToPostRequestInformation(ctx context.Context, body ia31f303b98dc4e7292d1559872ed38681eda57e78e48a431654df5b787bc8588.MoveRequestable, requestConfiguration *V2LocationWithUuItemRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST, m.BaseRequestBuilder.UrlTemplate, m.BaseRequestBuilder.PathParameters)
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    requestInfo.Headers.TryAdd("Accept", "application/json")
    err := requestInfo.SetContentFromParsable(ctx, m.BaseRequestBuilder.RequestAdapter, "application/json", body)
    if err != nil {
        return nil, err
    }
    return requestInfo, nil
}
// WithUrl returns a request builder with the provided arbitrary URL. Using this method means any other path or query parameters are ignored.
// returns a *V2LocationWithUuItemRequestBuilder when successful
func (m *V2LocationWithUuItemRequestBuilder) WithUrl(rawUrl string)(*V2LocationWithUuItemRequestBuilder) {
    return NewV2LocationWithUuItemRequestBuilder(rawUrl, m.BaseRequestBuilder.RequestAdapter);
}
