package automation

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"net/http"
)

// ActivityClient is the automation Client
type ActivityClient struct {
	BaseClient
}

// NewActivityClient creates an instance of the ActivityClient client.
func NewActivityClient(subscriptionID string, resourceGroupName string, clientRequestID string, automationAccountName string) ActivityClient {
	return NewActivityClientWithBaseURI(DefaultBaseURI, subscriptionID, resourceGroupName, clientRequestID, automationAccountName)
}

// NewActivityClientWithBaseURI creates an instance of the ActivityClient client.
func NewActivityClientWithBaseURI(baseURI string, subscriptionID string, resourceGroupName string, clientRequestID string, automationAccountName string) ActivityClient {
	return ActivityClient{NewWithBaseURI(baseURI, subscriptionID, resourceGroupName, clientRequestID, automationAccountName)}
}

// Get retrieve the activity in the module identified by module name and activity name.
//
// automationAccountName is the automation account name. moduleName is the name of module. activityName is the name
// of activity.
func (client ActivityClient) Get(ctx context.Context, automationAccountName string, moduleName string, activityName string) (result Activity, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.ResourceGroupName,
			Constraints: []validation.Constraint{{Target: "client.ResourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("automation.ActivityClient", "Get", err.Error())
	}

	req, err := client.GetPreparer(ctx, automationAccountName, moduleName, activityName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "automation.ActivityClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "automation.ActivityClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "automation.ActivityClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client ActivityClient) GetPreparer(ctx context.Context, automationAccountName string, moduleName string, activityName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"activityName":          autorest.Encode("path", activityName),
		"automationAccountName": autorest.Encode("path", automationAccountName),
		"moduleName":            autorest.Encode("path", moduleName),
		"resourceGroupName":     autorest.Encode("path", client.ResourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-10-31"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/modules/{moduleName}/activities/{activityName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client ActivityClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client ActivityClient) GetResponder(resp *http.Response) (result Activity, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// ListByModule retrieve a list of activities in the module identified by module name.
//
// automationAccountName is the automation account name. moduleName is the name of module.
func (client ActivityClient) ListByModule(ctx context.Context, automationAccountName string, moduleName string) (result ActivityListResultPage, err error) {
	if err := validation.Validate([]validation.Validation{
		{TargetValue: client.ResourceGroupName,
			Constraints: []validation.Constraint{{Target: "client.ResourceGroupName", Name: validation.Pattern, Rule: `^[-\w\._]+$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("automation.ActivityClient", "ListByModule", err.Error())
	}

	result.fn = client.listByModuleNextResults
	req, err := client.ListByModulePreparer(ctx, automationAccountName, moduleName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "automation.ActivityClient", "ListByModule", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListByModuleSender(req)
	if err != nil {
		result.alr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "automation.ActivityClient", "ListByModule", resp, "Failure sending request")
		return
	}

	result.alr, err = client.ListByModuleResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "automation.ActivityClient", "ListByModule", resp, "Failure responding to request")
	}

	return
}

// ListByModulePreparer prepares the ListByModule request.
func (client ActivityClient) ListByModulePreparer(ctx context.Context, automationAccountName string, moduleName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"automationAccountName": autorest.Encode("path", automationAccountName),
		"moduleName":            autorest.Encode("path", moduleName),
		"resourceGroupName":     autorest.Encode("path", client.ResourceGroupName),
		"subscriptionId":        autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2015-10-31"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/modules/{moduleName}/activities", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListByModuleSender sends the ListByModule request. The method will close the
// http.Response Body if it receives an error.
func (client ActivityClient) ListByModuleSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// ListByModuleResponder handles the response to the ListByModule request. The method always
// closes the http.Response Body.
func (client ActivityClient) ListByModuleResponder(resp *http.Response) (result ActivityListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listByModuleNextResults retrieves the next set of results, if any.
func (client ActivityClient) listByModuleNextResults(lastResults ActivityListResult) (result ActivityListResult, err error) {
	req, err := lastResults.activityListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "automation.ActivityClient", "listByModuleNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListByModuleSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "automation.ActivityClient", "listByModuleNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListByModuleResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "automation.ActivityClient", "listByModuleNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListByModuleComplete enumerates all values, automatically crossing page boundaries as required.
func (client ActivityClient) ListByModuleComplete(ctx context.Context, automationAccountName string, moduleName string) (result ActivityListResultIterator, err error) {
	result.page, err = client.ListByModule(ctx, automationAccountName, moduleName)
	return
}
