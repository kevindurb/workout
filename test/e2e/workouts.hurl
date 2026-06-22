GET http://localhost:1337/workouts/
HTTP 200
[Asserts]
xpath "string(/html/body/main/h1)" contains "Workouts"

GET http://localhost:1337/workouts/new
HTTP 200
[Asserts]
xpath "string(/html/body/main/h1)" contains "New Workout"
