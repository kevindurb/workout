POST {{base_url}}/login
[Form]
email: test@example.com
password: abc123
HTTP 302
[Captures]
session: cookie "session"

GET {{base_url}}/exercises
[Cookies]
session: {{session}}
HTTP 200
[Asserts]
xpath "string(/html/body/main/h1)" contains "Exercises"

GET {{base_url}}/exercises/new
[Cookies]
session: {{session}}
HTTP 200
[Asserts]
xpath "string(/html/body/main/h1)" contains "New Exercise"

POST {{base_url}}/exercises
[Cookies]
session: {{session}}
[Form]
name: test exercise
HTTP 302
[Captures]
exercise_url: header "Location"

GET {{base_url}}{{exercise_url}}
[Cookies]
session: {{session}}
HTTP 200
[Asserts]
xpath "string(/html/body/main/h1)" contains "test exercise"

POST {{base_url}}{{exercise_url}}
[Cookies]
session: {{session}}
[Form]
name: updated test exercise
HTTP 302

GET {{base_url}}{{exercise_url}}
[Cookies]
session: {{session}}
HTTP 200
[Asserts]
xpath "string(/html/body/main/h1)" contains "updated test exercise"

POST {{base_url}}{{exercise_url}}/delete
[Cookies]
session: {{session}}
HTTP 302
