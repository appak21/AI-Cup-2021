# Alem AI Cup 2021
## Pacman game  
*Description:* https://github.com/alem-io/cup-2021-rules/tree/main/about  
*Data format:* https://github.com/alem-io/cup-2021-rules/tree/main/format <br>

### Short Story
This competition gave me a lot of experience. It's my second time entering an AI cup and first time using Golang (last year I used C++ for the bomberman game). <br>
I used **A star search algorithm** for several purposes just by changing formulas, making it dependent on coins, bonus and monster values. <br>
After changing the algorithms as I needed, I tested different values of the coin and bonus prices. At 14 and less values of coins, my bot wasn't smart enough – always went to get the closest coins. 
That wasn't optimal in many cases. Setting its price to 30 and more made my bot choose more valuable places to go, 
where many coins are collected together. My bot started to look for more coins first, 
not for the closest ones – the decision was made according to the `total price of coins and minus total cost of distance`. 
I kept the first place and sometimes second, however, 3-2 hours before the last submission, 
I slightly changed my decision-maker algorithm 😄. As a result, I got the 4th place. <br>
Big thanks to Alem School for organizing the championship.
