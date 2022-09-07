# mid request

### /add
- url : callback address
- txt : callback value, json format
- req : timestamp to visit this url

every 5 second, it will fetch all req, and post them out

every 10 second, it will post a heart bit pack '{"bit":"1"}' to url find in conf.ini file