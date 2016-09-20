# gocrawl
A crawler written in golang. Initally this is just my attempt at learning golang, hopefully I can make it into something useful.

## Components
### URL consumer
0. Take a URL off the URL queue
1. Download the DOM
2. Parse links and put them on the URL queue
3. Put the DOM on the data queue

### Data consumer
0. Take the DOM off the data queue
1. Parse data from the DOM and store it in a data store (TODO: figure out how to make this configurable)