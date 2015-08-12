TODO:

* The current implementation will be really slow because data is fetched as it
is needed. Instead prefetch some user configured amount - maybe a block size.
* The current implementation is computationally intense because for every read
it readjusts the buffer to include only the range from min(subreader_offset) 
to max(subreader_offset) - fix this by only adjusting array size when you have
to fetch the next block