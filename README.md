# Valet an API Gateway

I had the idea that building an api gateway would be a fun project.

## Design concepts
api requests get made to the gateway and it makes requests to services behind the gateway. Collects all those responses and then returns the compiled response to the client. Design restraint self imposed to try and not use a database.


## Ideas
- handle authorization
- handle rate limiting
- concurrently hit the backing services using a form of the go load balancer pattern
- configuration
    - need to have some way to configure so it knows where to look for things.
- url patterms
    -  mygateway.foo.com/api1/* -> api1.someserver.com/*
        -  passes through what ever http action was used.
    -  the first  / argument points to which app api to use.

- ability to mock end points
    - you make a call and the api sends back test/dummy data instead of reaching out backing services

- be able to run as a redundant pair by. local server will need config to define the other node(s).
nodes will just send increment messages for rate limiting. 

- ability to load new configs via its own api, and then reload those configs via its own REST api

- administrative rest endpoint that can be used to query for the current / configured / available apis this gateway is fronting

- rate limit, send back limit exceeded response.
    - X-Ratelimit-Used: Approximate number of requests used in this period
    - X-Ratelimit-Remaining: Approximate number of requests left to use
    - X-Ratelimit-Reset: Approximate number of seconds to end of period

- gateway daemon currently runs on port 8000. This will be configurable at runtime later. 

- app config is a json file that looks like the below example. They are currently just being read out of the conf dir in the app root dir. This will change

to be clear is this the config for the applications that will be sitting behind the gateway. Not the gateway application it self. 
```
{
    "name": "userauth",
    "description": "service to authorize users of system",
    "authorize": true,
    "authKey": "thisismyauth",
    "authHeader": "X-Valet-Auth",
    "rateLimit": true,
    "limitValue": 2,
    "endpoints": [
        { "host": "auth.test.com", "path":"", "port":9050 },
        { "host": "auth.test.com", "path":"", "port":9051 }
    ]
}
```

## TODO
- ~ authorize api request ~
- ~ administrative rest endpoint that can be used to query for the current / configured / available apis this gateway is fronting ~
- ~ rate limit per end point ~
- ability to mock end points
- communication with pair
- decide what to communincate with pairs, master slave? or master master?
- ~ Fix config reload bug ~
- Figure model for storing stats and rate limiting info for public consumers
    Currently its built assuming one main consumer being a single enterprise. 
