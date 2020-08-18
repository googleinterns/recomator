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
