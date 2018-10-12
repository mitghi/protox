# Router

**Example** : Topic hierarchy:

```
|                                                  |<nil>|
|test                                             +|<nil>|(t)
|..ing                                            +|<nil>|(i)
|..../                                            +|<nil>|(/)
|......a                                          +|<nil>|(a)
|......../                                        +|<nil>|(/)
|..........sim                                     |<nil>|
|............ple                                  +|<nil>|(p)
|............../                                  +|<nil>|(/)
|................string                           +|&{map[client:0xc42001b5c0]}|(s)
|............ulation                              +|&{map[client3:0xc42001b740]}|(u)
|......*                                          +|&{map[client2:0xc42001b680 SPC :0xc42001b890]}|(*)
|../                                              +|<nil>|(/)
|....a                                            +|&{map[client4:0xc42001b800]}|(a)
-------------------------------
|                                                  |<nil>|
|test                                             +|<nil>|(t)
|..ing                                            +|<nil>|(i)
|..../                                            +|<nil>|(/)
|......a                                          +|<nil>|(a)
|......../                                        +|<nil>|(/)
|..........sim                                     |<nil>|
|............ple                                  +|<nil>|(p)
|............../                                  +|<nil>|(/)
|................string                           +|&{map[client:0xc42001ba10]}|(s)
|............ulation                              +|&{map[client3:0xc42001bb90]}|(u)
|......*                                          +|&{map[client2:0xc42001bad0 SPC :0xc42001bce0]}|(*)
|../                                              +|<nil>|(/)
|....a                                            +|&{map[client4:0xc42001bc50]}|(a)
-------------------------------
|                                                  |<nil>|
|test                                             +|<nil>|(t)
|..ing                                            +|<nil>|(i)
|..../                                            +|<nil>|(/)
|......a                                          +|<nil>|(a)
|......../                                        +|<nil>|(/)
|..........sim                                     |<nil>|
|............ple                                  +|<nil>|(p)
|............../                                  +|<nil>|(/)
|................string                           +|&{map[client:0xc42001ba10]}|(s)
|............ulation                              +|&{map[client3:0xc42001bb90]}|(u)
|../                                              +|<nil>|(/)
-------------------------------
|                                                  |<nil>|
|a                                                +|<nil>|(a)
|../                                              +|<nil>|(/)
|....simple                                       +|<nil>|(s)
|....../                                          +|<nil>|(/)
|........path                                     +|&{map[client:0xc42001be60]}|(p)
|........*                                        +|&{map[client:0xc42001bf20]}|(*)
|........../                                      +|<nil>|(/)
|............thing                                +|&{map[client:0xc42001bf80]}|(t)
|....another                                      +|<nil>|(a)
|....../                                          +|<nil>|(/)
|........sim                                      +|&{map[client:0xc420104360]}|(s)
|..........ple                                    +|<nil>|(p)
|............/                                    +|<nil>|(/)
|..............thing                              +|&{map[client:0xc420104060]}|(t)
|..........ul                                     +|&{map[client:0xc4201042a0]}|(u)
|............ating                                +|<nil>|(a)
|............../                                  +|<nil>|(/)
|................thing                            +|&{map[client:0xc420104120]}|(t)
|..........a                                      +|&{map[client:0xc420104420]}|(a)
|..a                                              +|<nil>|(a)
|..../                                            +|<nil>|(/)
|......branch                                     +|&{map[client:0xc4201041e0]}|(b)
-------------------------------
client 1
|                                                  |<nil>|
|a                                                +|<nil>|(a)
|../                                              +|<nil>|(/)
|....simple                                       +|<nil>|(s)
|....../                                          +|<nil>|(/)
|........path                                     +|&{map[client1:0xc420104660]}|(p)
|....*                                            +|&{map[client2:0xc420104720 client4:0xc420104870]}|(*)
|....../                                          +|<nil>|(/)
|........location                                 +|&{map[client4:0xc4201048d0]}|(l)
|....another                                      +|<nil>|(a)
|....../                                          +|<nil>|(/)
|........simple                                   +|<nil>|(s)
|........../                                      +|<nil>|(/)
|............thing                                +|&{map[client3:0xc4201047e0]}|(t)
-------------------------------
client3, &{topic:a/another/simple/thing uid:client3 qos:1 isLeaf:true}
true
client2, &{topic:a/* uid:client2 qos:1 isLeaf:true}
true
|                                                  |<nil>|
|a                                                +|<nil>|(a)
|../                                              +|<nil>|(/)
|....simple                                       +|<nil>|(s)
|....../                                          +|<nil>|(/)
|........path                                     +|&{map[client1:0xc420104b40]}|(p)
|....*                                            +|&{map[client2:0xc420104c00 client4:0xc420104d50]}|(*)
|....../                                          +|<nil>|(/)
|........location                                 +|&{map[client4:0xc420104db0]}|(l)
|....another                                      +|<nil>|(a)
|....../                                          +|<nil>|(/)
|........simple                                   +|<nil>|(s)
|........../                                      +|<nil>|(/)
|............thing                                +|&{map[client3:0xc420104cc0]}|(t)
-------------------------------
client2, &{topic:a/* uid:client2 qos:1 isLeaf:true}
true
client3, &{topic:a/another/simple/thing uid:client3 qos:1 isLeaf:true}
true
client2, &{topic:a/*/topic uid:client2 qos:2 isLeaf:false}
false
client2, &{topic:a/awesome/topic uid:client2 qos:2 isLeaf:true}
false
client3, &{topic:a/* uid:client3 qos:0 isLeaf:true}
true
client1, &{topic:a/awesome/topic uid:client1 qos:1 isLeaf:true}
false
PASS
ok  	github.com/mitghi/protox/server	0.022s
```
