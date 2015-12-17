# Valet an API Gateway

I had the idea that building an api gateway would be a fun project.

## Design concepts
api requests get made to the gateway and it makes requests to services behind the gateway. Collects all those responses and then returns the compiled response to the client.

## Ideas
- handle authentication
-  handle rate limiting
-  concurrently hit the backing services using a form of the go load balancer pattern
- configuration
    - need to have some way to configure so it knows where to look for things.
- url patterms
    -  mygateway.foo.com/api1/* -> api1.someserver.com/*
        -  passes through what ever http action was used.
    -  the first  / argument points to 
