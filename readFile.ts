import { readFile } from "node:fs";
import data from "./data.json";
import { name, age } from "./data.json";

readFile("./data.json", "utf-8", (err, data) => {
  if (err) console.error(err);

  console.log(data);
});

console.log(data);
console.log(name, age);
