


func forumDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	helper.Logger("info", "Forum Detail success")
}