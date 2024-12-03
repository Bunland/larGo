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
type ServeObject = {
  port: string
}
type ConsoleObject = {
  time: (label?: string) => undefined;
  timeEnd: (label?: string) => undefined;
  timeLog: (label?: string, ...data: any[]) => undefined;
  log: (...data: any[]) => undefined;
  clear: () => undefined;
  error: (...data: string[]) => undefined;
  assert: (expression: boolean, ...obj: any[]) => undefined;
  warn: (...data: string[]) => undefined;
  count: (label?: string) => undefined;
  countReset: (label?: string) => undefined;
};
type ModuleObject = {
  exports: any
}
type ProcessObject = { version: string };
type Awaitable<T> = Promise<T> | T
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
declare const serve: (obj: ServeObject) => unknown;
declare const get: (routeName: string, callback: () => Awaitable<unknown>) => unknown;

declare module "fs" {
  function readFileSync(fileName: string, obj?: ReadFileObject): string;
  function writeFileSync(fileName: string, content: string): undefined;
  function readdirSync(directoryName: string): Array<string>;
  function mkdirSync(directoryName: string): undefined;
  function readFile(
    fileName: string,
    encoding: "utf8" | "utf-8" | "base64" | "hex",
    callback: (err: string, data: string) => void
  ): undefined;
}

declare module "node:fs" {
  function readFileSync(fileName: string, obj?: ReadFileObject): string;
  function writeFileSync(fileName: string, content: string): undefined;
  function readdirSync(directoryName: string): Array<string>;
  function mkdirSync(directoryName: string): undefined;
  function readFile(
    fileName: string,
    encoding: "utf-8" | "utf8" | "base64" | "hex",
    callback: (err: string, data: string) => void
  ): undefined;
}
