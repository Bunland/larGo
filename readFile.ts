import { readFile } from "node:fs";
import data from "./data.json";
import { name, age } from "./data.json";

readFile("./data.json", "utf-8", (err, dataFile) => {
  if (err) console.error(err);

  console.log(dataFile);
  console.log(JSON.parse(dataFile).name)
});

console.log(data);
console.log()
console.log(name, age);
