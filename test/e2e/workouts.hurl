POST http://localhost:1337/login
[Form]
email: test@example.com
password: abc123
HTTP 302
[Captures]
session: cookie "session"

GET http://localhost:1337/workouts/
[Cookies]
session: {{session}}
HTTP 200
[Asserts]
xpath "string(/html/body/main/h1)" contains "Workouts"

GET http://localhost:1337/workouts/new
[Cookies]
session: {{session}}
HTTP 200
[Asserts]
xpath "string(/html/body/main/h1)" contains "New Workout"
