POST {{base_url}}/login
[Form]
email: test@example.com
password: abc123
HTTP 302
[Captures]
session: cookie "session"

GET {{base_url}}/workouts/
[Cookies]
session: {{session}}
HTTP 200
[Asserts]
xpath "string(/html/body/main/h1)" contains "Workouts"

GET {{base_url}}/workouts/new
[Cookies]
session: {{session}}
HTTP 200
[Asserts]
xpath "string(/html/body/main/h1)" contains "New Workout"

POST {{base_url}}/workouts/
[Cookies]
session: {{session}}
[Form]
name: test workout
HTTP 302
[Captures]
workout_url: header "Location"

GET {{base_url}}{{workout_url}}
[Cookies]
session: {{session}}
HTTP 200
[Asserts]
xpath "string(/html/body/main/h1)" contains "test workout"

POST {{base_url}}{{workout_url}}
[Cookies]
session: {{session}}
[Form]
name: updated test workout
HTTP 302

GET {{base_url}}{{workout_url}}
[Cookies]
session: {{session}}
HTTP 200
[Asserts]
xpath "string(/html/body/main/h1)" contains "updated test workout"

POST {{base_url}}{{workout_url}}/delete
[Cookies]
session: {{session}}
HTTP 302
