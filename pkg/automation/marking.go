/*
Copyright 2020 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package automation

import (
	"google.golang.org/api/recommender/v1"
)

type gcloudClaimedRequest = recommender.GoogleCloudRecommenderV1MarkRecommendationClaimedRequest
type gcloudFailedRequest = recommender.GoogleCloudRecommenderV1MarkRecommendationFailedRequest
type gcloudSucceededRequest = recommender.GoogleCloudRecommenderV1MarkRecommendationSucceededRequest

// Marks the recommendation defined by the given name and etag as claimed
func (s *googleService) MarkRecommendationClaimed(name, etag string) (*gcloudRecommendation, error) {
	r := recommender.NewProjectsLocationsRecommendersRecommendationsService(s.recommenderService)
	request := gcloudClaimedRequest{
		Etag: etag,
	}

	markClaimedCall := r.MarkClaimed(name, &request)
	return markClaimedCall.Do()
}

// Marks the recommendation defined by the given name and etag as failed
func (s *googleService) MarkRecommendationFailed(name, etag string) (*gcloudRecommendation, error) {
	r := recommender.NewProjectsLocationsRecommendersRecommendationsService(s.recommenderService)
	request := gcloudFailedRequest{
		Etag: etag,
	}

	markFailedCall := r.MarkFailed(name, &request)
	return markFailedCall.Do()
}

// Marks the recommendation defined by the given name and etag as succeeded
func (s *googleService) MarkRecommendationSucceeded(name, etag string) (*gcloudRecommendation, error) {
	r := recommender.NewProjectsLocationsRecommendersRecommendationsService(s.recommenderService)
	request := gcloudSucceededRequest{
		Etag: etag,
	}

	markSucceededCall := r.MarkSucceeded(name, &request)
	return markSucceededCall.Do()
}
