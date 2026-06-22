GET http://localhost:1337/
HTTP 200
[Asserts]
xpath "string(/html/head/title)" contains "Workout"
xpath "string(/html/body/main/h1)" contains "Workout"
