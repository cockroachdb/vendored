// +build go1.9

// Copyright 2018 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This code was auto-generated by:
// github.com/Azure/azure-sdk-for-go/tools/profileBuilder

package eventgrid

import original "github.com/Azure/azure-sdk-for-go/services/eventgrid/mgmt/2018-01-01/eventgrid"

const (
	DefaultBaseURI = original.DefaultBaseURI
)

type BaseClient = original.BaseClient

func New(subscriptionID string) BaseClient {
	return original.New(subscriptionID)
}
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return original.NewWithBaseURI(baseURI, subscriptionID)
}

type EventSubscriptionsClient = original.EventSubscriptionsClient

func NewEventSubscriptionsClient(subscriptionID string) EventSubscriptionsClient {
	return original.NewEventSubscriptionsClient(subscriptionID)
}
func NewEventSubscriptionsClientWithBaseURI(baseURI string, subscriptionID string) EventSubscriptionsClient {
	return original.NewEventSubscriptionsClientWithBaseURI(baseURI, subscriptionID)
}

type EndpointType = original.EndpointType

const (
	EndpointTypeEventHub                     EndpointType = original.EndpointTypeEventHub
	EndpointTypeEventSubscriptionDestination EndpointType = original.EndpointTypeEventSubscriptionDestination
	EndpointTypeWebHook                      EndpointType = original.EndpointTypeWebHook
)

type EventSubscriptionProvisioningState = original.EventSubscriptionProvisioningState

const (
	Canceled  EventSubscriptionProvisioningState = original.Canceled
	Creating  EventSubscriptionProvisioningState = original.Creating
	Deleting  EventSubscriptionProvisioningState = original.Deleting
	Failed    EventSubscriptionProvisioningState = original.Failed
	Succeeded EventSubscriptionProvisioningState = original.Succeeded
	Updating  EventSubscriptionProvisioningState = original.Updating
)

type ResourceRegionType = original.ResourceRegionType

const (
	GlobalResource   ResourceRegionType = original.GlobalResource
	RegionalResource ResourceRegionType = original.RegionalResource
)

type TopicProvisioningState = original.TopicProvisioningState

const (
	TopicProvisioningStateCanceled  TopicProvisioningState = original.TopicProvisioningStateCanceled
	TopicProvisioningStateCreating  TopicProvisioningState = original.TopicProvisioningStateCreating
	TopicProvisioningStateDeleting  TopicProvisioningState = original.TopicProvisioningStateDeleting
	TopicProvisioningStateFailed    TopicProvisioningState = original.TopicProvisioningStateFailed
	TopicProvisioningStateSucceeded TopicProvisioningState = original.TopicProvisioningStateSucceeded
	TopicProvisioningStateUpdating  TopicProvisioningState = original.TopicProvisioningStateUpdating
)

type TopicTypeProvisioningState = original.TopicTypeProvisioningState

const (
	TopicTypeProvisioningStateCanceled  TopicTypeProvisioningState = original.TopicTypeProvisioningStateCanceled
	TopicTypeProvisioningStateCreating  TopicTypeProvisioningState = original.TopicTypeProvisioningStateCreating
	TopicTypeProvisioningStateDeleting  TopicTypeProvisioningState = original.TopicTypeProvisioningStateDeleting
	TopicTypeProvisioningStateFailed    TopicTypeProvisioningState = original.TopicTypeProvisioningStateFailed
	TopicTypeProvisioningStateSucceeded TopicTypeProvisioningState = original.TopicTypeProvisioningStateSucceeded
	TopicTypeProvisioningStateUpdating  TopicTypeProvisioningState = original.TopicTypeProvisioningStateUpdating
)

type EventHubEventSubscriptionDestination = original.EventHubEventSubscriptionDestination
type EventHubEventSubscriptionDestinationProperties = original.EventHubEventSubscriptionDestinationProperties
type EventSubscription = original.EventSubscription
type BasicEventSubscriptionDestination = original.BasicEventSubscriptionDestination
type EventSubscriptionDestination = original.EventSubscriptionDestination
type EventSubscriptionFilter = original.EventSubscriptionFilter
type EventSubscriptionFullURL = original.EventSubscriptionFullURL
type EventSubscriptionProperties = original.EventSubscriptionProperties
type EventSubscriptionsCreateOrUpdateFuture = original.EventSubscriptionsCreateOrUpdateFuture
type EventSubscriptionsDeleteFuture = original.EventSubscriptionsDeleteFuture
type EventSubscriptionsListResult = original.EventSubscriptionsListResult
type EventSubscriptionsUpdateFuture = original.EventSubscriptionsUpdateFuture
type EventSubscriptionUpdateParameters = original.EventSubscriptionUpdateParameters
type EventType = original.EventType
type EventTypeProperties = original.EventTypeProperties
type EventTypesListResult = original.EventTypesListResult
type Operation = original.Operation
type OperationInfo = original.OperationInfo
type OperationsListResult = original.OperationsListResult
type Resource = original.Resource
type Topic = original.Topic
type TopicProperties = original.TopicProperties
type TopicRegenerateKeyRequest = original.TopicRegenerateKeyRequest
type TopicsCreateOrUpdateFuture = original.TopicsCreateOrUpdateFuture
type TopicsDeleteFuture = original.TopicsDeleteFuture
type TopicSharedAccessKeys = original.TopicSharedAccessKeys
type TopicsListResult = original.TopicsListResult
type TopicsUpdateFuture = original.TopicsUpdateFuture
type TopicTypeInfo = original.TopicTypeInfo
type TopicTypeProperties = original.TopicTypeProperties
type TopicTypesListResult = original.TopicTypesListResult
type TopicUpdateParameters = original.TopicUpdateParameters
type TrackedResource = original.TrackedResource
type WebHookEventSubscriptionDestination = original.WebHookEventSubscriptionDestination
type WebHookEventSubscriptionDestinationProperties = original.WebHookEventSubscriptionDestinationProperties
type OperationsClient = original.OperationsClient

func NewOperationsClient(subscriptionID string) OperationsClient {
	return original.NewOperationsClient(subscriptionID)
}
func NewOperationsClientWithBaseURI(baseURI string, subscriptionID string) OperationsClient {
	return original.NewOperationsClientWithBaseURI(baseURI, subscriptionID)
}

type TopicsClient = original.TopicsClient

func NewTopicsClient(subscriptionID string) TopicsClient {
	return original.NewTopicsClient(subscriptionID)
}
func NewTopicsClientWithBaseURI(baseURI string, subscriptionID string) TopicsClient {
	return original.NewTopicsClientWithBaseURI(baseURI, subscriptionID)
}

type TopicTypesClient = original.TopicTypesClient

func NewTopicTypesClient(subscriptionID string) TopicTypesClient {
	return original.NewTopicTypesClient(subscriptionID)
}
func NewTopicTypesClientWithBaseURI(baseURI string, subscriptionID string) TopicTypesClient {
	return original.NewTopicTypesClientWithBaseURI(baseURI, subscriptionID)
}
func UserAgent() string {
	return original.UserAgent() + " profiles/latest"
}
func Version() string {
	return original.Version()
}
