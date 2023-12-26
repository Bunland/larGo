import fs from "fs";
const text = fs.readFileSync("example.txt", { encoding: "utf-8" })
console.log(text)