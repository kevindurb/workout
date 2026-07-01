# vim: foldmethod=marker

# SETUP {{{
POST {{base_url}}/login
[Form]
email: test@example.com
password: abc123
HTTP 302
[Captures]
session: cookie "session"

POST {{base_url}}/workouts
[Cookies]
session: {{session}}
[Form]
name: test workout
HTTP 302
[Captures]
workout_url: header "Location"
workout_id: variable "workout_url" split "/" nth 2

POST {{base_url}}/exercises
[Cookies]
session: {{session}}
[Form]
name: test exercise
HTTP 302
[Captures]
exercise_url: header "Location"
exercise_id: variable "exercise_url" split "/" nth 2
# }}}

# TEST {{{
POST {{base_url}}/workouts_exercises
[Cookies]
session: {{session}}
[Form]
workout_id: {{workout_id}}
exercise_id: {{exercise_id}}
HTTP 302

GET {{base_url}}{{workout_url}}
[Cookies]
session: {{session}}
HTTP 200
[Asserts]
xpath "string(/html/body/main/h1)" contains "test workout"
xpath "string(/html/body/main/ul/li)" contains "test exercise"
# }}}

# CLEANUP {{{
POST {{base_url}}{{exercise_url}}/delete
[Cookies]
session: {{session}}
HTTP 302

POST {{base_url}}{{workout_url}}/delete
[Cookies]
session: {{session}}
HTTP 302
# }}}
