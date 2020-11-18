#Mandatory 3

## A supermarket

### The idea of the application

This is a simple simulation of a supermarket with tills, and an opening time and a closing time. _(Time is scaled down so that 1 hour is one second and so on)_

+ Because of corona restrictions the supermarket can only hold 50 customers shopping and the lines can each only hold 10 customers. 
+ When the shop opens, customers will arrive to the door constantly with a random interval, and they will be let in if there is less than 50 in the shop already.
+ Each customer has an individual duration they will use to shop, and they will each collect different amounts of groceries with different prices.
+ When a customer finishes his shopping, he will try to go to a line. If the line is full he will check another line until he finds a line which has space for him.
+ The till can each deal with one customer at a time, and has an execution time on 2 milliseconds 
+ When the time is 22, the front door will be closed, and no more customers are allowed in.
+ The tills will then finish the last customers and when there is no more customers in the shop, the tills will close and the shop will close totally.
+ The revenue of each till will then be printed.

### Why is concurrency needed

+ The clock counting from 7 to 22 where the door should be closed needs its own routine
+ each customer doing his shopping on his own time
+ the tills handling the customers (here sharing the load)

### What could be the potential issues specifically

+ If everything was running on the main thread, then the clock would just run without anything happening, and then when the shop would close, then everything would happen. 
+ Only one customer as the time would come through the shop at a time then the shop would not make as much money.
+ If the tills would not be able to deal share the load of customers equally then the flow would be in balanced and one till might be empty while another would be over full, and it would slow the execution time down, and lead to less money ( and potentially angry customers )

### Address race conditions / solutions

+ The tills are used in many different routines, so when a routine needs to read or write to the till then it is a good idea to use Mutual Exclusion to avoid race conditions.

+ I am using channels a lot in this program and the channels in go are avoiding race conditions when sending customers from one place to another, by copying the customer instead of pointing to it, when they are enqueuing/dequeuing the channel, and for the channel itself they are using mutual exclusion to avoid race conditions.

* Later when sharing the data from the tills to the shop global account then a race conditions could also appear. 

### Address deadlocks and starvation. (and solutions)

* The shop is only open for a specific duration and therefore I am using something in go called a "wait group" which is a thread safe conditional variable which can get _n_ workers assigned to them. I am assigning the tills and when the tills are closed then they will let the "wait group" know that they are done and when all tills are done then the wait group will let the program continue. (If a till is never marked as done then the wait group will create a deadlock)

* Also there is an aspect where customers are flowing constantly into the shop, during the opening time, So I made my own conditional variable, which is either 0 or 1. 
    + When the shop has been open for the full interval between the opening time and the closing time, a channel will be marked as closed. 
    + This channel will be looked at at the customer spawner and when the channel is closed then the spawner will stop spawning customers.
    + In the tills they will each check if they have customers in the line, and 
        + if there is no customers in the line (starvation) then they will check if this channel is closed, and 
        + if the channel is closed then they will close their line and return their functions. 
        + _(Comments can be found in the code)_

