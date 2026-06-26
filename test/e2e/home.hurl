POST {{base_url}}/login
[Form]
email: test@example.com
password: abc123
HTTP 302
[Captures]
session: cookie "session"

GET {{base_url}}/
[Cookies]
session: {{session}}
HTTP 200
[Asserts]
xpath "string(/html/head/title)" contains "Workout"
xpath "string(/html/body/main/h1)" contains "Workout"
