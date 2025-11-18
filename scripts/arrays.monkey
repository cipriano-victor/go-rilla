let zero = 0;
let two = 2;
let three = 3;

let a = [zero += 1, three -= 1, ++two, 2**2];

print(a[0]);
print(a[4 - 3]);
print(a[6 / 3]);
print(a[99]); 

print(first(a));
print(last(a));
print(rest(a));
print(rest(rest(a)));

let b = push(a, 5);
print(b);
len(b);
