# Todo list

## Backlog
* Add timers
* Display timers
* Flags to add
  * -f flood 
  * -l preload (but with safety rails)
  * -p pattern Let the user specify up to 16 bytes to be used as the data packet.
  * -q Quiet
  * -M hint. set the MTU Discovery strategy
    * TODO does apple even have this? 
    * do (prohibit fragmentation, even local one)
    * want (do PMTU discovery, fragment locally when packet size is large)
    * dont (do not set DF flag)
* Tests
* Concurrency y'all
* Summary/Statistics line