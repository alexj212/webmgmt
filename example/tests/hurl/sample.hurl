# ------------------------------------
POST http://localhost:1100/auth/signup
{
  "name": "Alex J",
  "email": "alexj",
  "password": "gjpwl0v3",
  "role": "admin"
}
HTTP/1.1 200
[Asserts]
jsonpath "$['access_token']" exists
jsonpath "$['refresh_token']" exists
[Captures]
access_token: jsonpath "$['access_token']"
refresh_token: jsonpath "$['refresh_token']"

