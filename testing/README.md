# Acknowledgement
The test data originates from [PuerkitoBio/gocrawl](https://github.com/PuerkitoBio/gocrawl)

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
