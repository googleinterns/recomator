package automation

import (
	"google.golang.org/api/recommender/v1"
)

type gcloudClaimedRequest = recommender.GoogleCloudRecommenderV1MarkRecommendationClaimedRequest
type gcloudFailedRequest = recommender.GoogleCloudRecommenderV1MarkRecommendationFailedRequest
type gcloudSucceededRequest = recommender.GoogleCloudRecommenderV1MarkRecommendationSucceededRequest

// Marks the recommendation defined by the given name and etag as claimed
func (s *googleService) MarkRecommendationClaimed(name, etag string) error {
	r := recommender.NewProjectsLocationsRecommendersRecommendationsService(s.recommenderService)
	request := gcloudClaimedRequest{
		Etag: etag,
	}

	markClaimedCall := r.MarkClaimed(name, &request)
	_, err := markClaimedCall.Do()

	return err
}

// Marks the recommendation defined by the given name and etag as failed
func (s *googleService) MarkRecommendationFailed(name, etag string) error {
	r := recommender.NewProjectsLocationsRecommendersRecommendationsService(s.recommenderService)
	request := gcloudFailedRequest{
		Etag: etag,
	}

	markFailedCall := r.MarkFailed(name, &request)
	_, err := markFailedCall.Do()

	return err
}

// Marks the recommendation defined by the given name and etag as succeeded
func (s *googleService) MarkRecommendationSucceeded(name, etag string) error {
	r := recommender.NewProjectsLocationsRecommendersRecommendationsService(s.recommenderService)
	request := gcloudSucceededRequest{
		Etag: etag,
	}

	markSucceededCall := r.MarkSucceeded(name, &request)
	_, err := markSucceededCall.Do()

	return err
}
