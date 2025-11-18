let add = fn (x, y) { 
    return x + y;
}

let sub = fn (x, y) {
    return x - y;
}

let apply_function = fn (x, y, f) {
    f(x,y)
}

print(apply_function(3,2,add)); 
print(apply_function(3,2,sub));