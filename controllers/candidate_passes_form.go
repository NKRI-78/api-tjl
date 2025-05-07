package controllers

import (
	"encoding/json"
	"net/http"
	"superapps/entities"
	helper "superapps/helpers"
	service "superapps/services"
)

func CandidatePassesForm(w http.ResponseWriter, r *http.Request) {
	data := &entities.DepartureForm{}

	err := json.NewDecoder(r.Body).Decode(data)

	if err != nil {
		helper.Logger("error", "In Server: "+err.Error())
		helper.Response(w, 400, true, "Internal server error ("+err.Error()+")", map[string]any{})
		return
	}

	// DateDeparture := data.DateDeparture
	// TimeDeparture := data.TimeDeparture
	// Airplane := data.Airplane
	// Location := data.Location
	// Destination := data.Destination
	UserCandidateId := data.UserCandidateId
	ApplyJobId := data.ApplyJobId

	// if DateDeparture == "" {
	// 	helper.Logger("error", "In Server: date_departure required")
	// 	helper.Response(w, 400, true, "date_departure is required", map[string]any{})
	// 	return
	// }

	// if TimeDeparture == "" {
	// 	helper.Logger("error", "In Server: time_departure required")
	// 	helper.Response(w, 400, true, "time_departure is required", map[string]any{})
	// 	return
	// }

	// if Airplane == "" {
	// 	helper.Logger("error", "In Server: airplane required")
	// 	helper.Response(w, 400, true, "airplane is required", map[string]any{})
	// 	return
	// }

	// if Location == "" {
	// 	helper.Logger("error", "In Server: location required")
	// 	helper.Response(w, 400, true, "location is required", map[string]any{})
	// 	return
	// }

	// if Destination == "" {
	// 	helper.Logger("error", "In Server: destination required")
	// 	helper.Response(w, 400, true, "destination is required", map[string]any{})
	// 	return
	// }

	if UserCandidateId == "" {
		helper.Logger("error", "In Server: user_candidate_id required")
		helper.Response(w, 400, true, "user_candidate_id is required", map[string]any{})
		return
	}

	if ApplyJobId == "" {
		helper.Logger("error", "In Server: apply_job_id required")
		helper.Response(w, 400, true, "apply_job_id is required", map[string]any{})
		return
	}

	result, err := service.CandidatePassesForm(data)

	if err != nil {
		helper.Response(w, 400, true, err.Error(), map[string]any{})
		return
	}

	helper.Logger("info", "Candidate Passes Form success")
	helper.Response(w, http.StatusOK, false, "Successfully", map[string]any{
		// "date_departure":    result["date_departure"],
		// "time_departure":    result["time_departure"],
		// "departure_id":      result["departure_id"],
		// "destination":       result["destination"],
		// "airplane":          result["airplane"],
		// "location":          result["location"],
		"content":           result["content"],
		"departure_id":      result["departure_id"],
		"apply_job_id":      result["apply_job_id"],
		"user_candidate_id": result["user_candidate_id"],
	})
}
