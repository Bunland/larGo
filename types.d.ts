interface ReadFileObject {
  encoding: "utf8" | "utf-8" | "base64" | "hex";
}
type ColorType =
  | "red"
  | "green"
  | "blue"
  | "black"
  | "cyan"
  | "yellow"
  | "magenta"
  | "white";
type LarObject = {
  color: (colorType: ColorType, ...data: any[]) => undefined;
  version: string;
};
type ConsoleObject = {
  time: (label: string = "default") => undefined;
  timeEnd: (label: string = "default") => undefined;
  timeLog: (label: string = "default", ...data: any[]) => undefined;
  log: (...data: any[]) => undefined;
  clear: () => undefined;
  error: (...data: string[]) => undefined;
  assert: (expression: boolean, ...obj: any[]) => undefined;
  warn: (...data: string[]) => undefined;
  count: (label: string = "default") => undefined;
  countReset: (label: string = "default") => undefined;
};
type ProcessObject = { version: string };
declare const Lar: LarObject;
declare const console: ConsoleObject;
declare const Add: (a: number, b: number) => number;
declare const Mult: (a: number, b: number) => number;
declare const Div: (a: number, b: number) => number;
declare const print: (...data: any[]) => undefined;
declare const module: ModuleObject;
declare const prompt: (question: string) => string;
declare const process: ProcessObject;
declare const __dirname: string;
/** @deprecated */
declare const __filename: string;

declare module "fs" {
  declare function readFileSync(fileName: string, obj?: ReadFileObject): string;
  declare function writeFileSync(fileName: string, content: string): undefined;
  declare function readdirSync(directoryName: string): Array<string>;
  declare function mkdirSync(directoryName: string): undefined;
  declare function readFile(
    fileName: string,
    encoding: "utf8" | "utf-8" | "base64" | "hex",
    callback: (err: string, data: string) => void
  ): undefined;
}

declare module "node:fs" {
  declare function readFileSync(fileName: string, obj?: ReadFileObject): string;
  declare function writeFileSync(fileName: string, content: string): undefined;
  declare function readdirSync(directoryName: string): Array<string>;
  declare function mkdirSync(directoryName: string): undefined;
  declare function readFile(
    fileName: string,
    encoding: "utf-8" | "utf8" | "base64" | "hex",
    callback: (err: string, data: string) => void
  ): undefined;
}
