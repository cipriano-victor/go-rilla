let i = 1;

for(; i <= 100; i++) {

    if((i == 50) || ((i ** 0) != 1)){
        break;
    }

    if(((i % 5) == 0) && ((i % 3) == 0)) {
        print("FizzBuzz");
        continue;
    } 
    
    if((i % 3) == 0) {
        print("Fizz");
        continue;
    } 
    
    if((i % 5) == 0) {
        print("Buzz");
        continue;
    }

    print(i);
}

print(i);