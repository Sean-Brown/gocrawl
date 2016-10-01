# Integration Testing

## How to test

### Setup
To test serving these pages through localhost, add to the hosts file
* Windows - C:\Windows\system32\drivers\etc\hosts
* Linux - /etc/hosts

The testing directory has hosta, hostb, hostc, hostd with files to serve.
Edit the hosts file, for example, on Windows add these lines to indicate to Windows
that we want requests to localhost (127.0.0.1) to go to the host (e.g. hosta):

>127.0.0.1	hosta  
>127.0.0.1	hostb  
>127.0.0.1	hostc  
>127.0.0.1	hostd  

NOTE: this will not work behind a corporate proxy (or at least I need to figure out
how to configure it for a corporate proxy)

### Testing
To run the integration tests, open a terminal and type:
<pre>go test -v github.com/Sean-Brown/gocrawl/tests</pre>

## Acknowledgement
The html test data originates from [PuerkitoBio/gocrawl](https://github.com/PuerkitoBio/gocrawl)
