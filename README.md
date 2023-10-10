# Testing Keepalives

Example usage: 

```
go build
./keepalive-test -address "localhost:80" -period 59 -interval 2
```

This will open a TCP connection, send a request with to localhost at port 80, wait 59 seconds, and then send a second request. Assuming this succeeds, this script will repeat this with 61 seconds, and then keep incrementing the wait duration by 2 seconds until the request fails. Right now, this is configured to have a hard stop at 305 seconds, since it's uncommon to set keepalive idle timeouts longer than 5 minutes (aka 300 seconds). 

## Notes

This code was written to verify the `keepalive_timeout` behavior of nginx. A few notes and common mistakes about testing this type of behavior if you don't understand Go (like me, yay) and just try to stackoverflow this: 

1. This code sends a `HTTP/1.0` GET request with the `Connection: keep-alive` through the TCP connection. This is necessary if you want to tell the server to keep the connection alive. HTTP 1.0 defaults to `Connection: close`, which means that without specifying this, the server will just close the connection after the first request. 
2. I initially made the mistake of assuming that I could verify the connection was terminated by checking for `syscall.ECONNRESET`. This is incorrect. Not all servers will send an RST or terminate the connection. Sometimes it will be possible to send a request and read EOF from an invalid connection which will not respond to new requests. Basically you can't just try to read from the connection, get EOF, and assume the connection is valid. You have to send a request and verify that you get a valid response. 
3. "Reading until the end" in Go is usually reading until EOF. However, in TCP connections, HTTP requests you get back are not guaranteed to send an EOF. Sometimes you will only get an EOF when the connection is terminated on the other side. So "trying to read until an EOF" will sometimes just block and hang until the connection is terminated. 

There are probably more edgecases that this code does not cover. I mostly just wanted to jot down the common set of pitfalls I ran into to save folks some time, in case anyone comes across the need to verify keepalive_timeouts in the future. 



