POST http://localhost:1337/login
[Form]
email: test@example.com
password: abc123
HTTP 302
[Captures]
session: cookie "session"

GET http://localhost:1337/
[Cookies]
session: {{session}}
HTTP 200
[Asserts]
xpath "string(/html/head/title)" contains "Workout"
xpath "string(/html/body/main/h1)" contains "Workout"
