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

### Address race conditions

+ If  

## A bank

