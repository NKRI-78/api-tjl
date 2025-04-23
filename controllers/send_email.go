package controllers

import (
	"net/http"
	helper "superapps/helpers"
)

func SendEmail(w http.ResponseWriter, r *http.Request) {

	data := `<div style='font-family: Helvetica, Arial, sans-serif; min-width: 1000px; overflow: auto; line-height: 2;'> 
            <div style='margin: 50px auto; width: 70%; padding: 20px 0;'> 
                <p style='font-size: 1.1em;'> Hi, </p> 
                <p>Use the following OTP to complete your Sign Up procedures. OTP is valid for 2 minutes</p>
                <h2 style='background: #00466a; margin: 0 auto; width: max-content; padding: 0 10px; color: #fff; border-radius: 4px;'>tset</h2>
                <p style='font-size: 0.9em;'>Regards, <br/>${name}</p>
                <hr style='border: none; border-top: 1px solid #eee;' />
            </div>
        </div>`

	helper.SendEmail("reihanagam7@gmail.com", "TJL", "Verification Account", data, "")

	helper.Response(w, http.StatusOK, false, "Successfully", map[string]any{})
}
