console.time();
let n = 0;
for (let i = 0; i < 100000; i++) {
  n += i;
}
console.log(n);
console.timeLog("default", new Error("test"));
let x = 0;
for (let i = 0; i < 100000; i++) {
  x += i;
}
console.log(x);
console.timeEnd();