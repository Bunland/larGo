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
  port: number
}
type ConsoleObject = {
  time: (label?: string) => undefined;
  timeEnd: (label?: string) => undefined;
  timeLog: (label?: string, ...data: any[]) => undefined;
  /**
   * Logs anything to the console
   * @param data The data to log
   * @returns undefined
   * @example console.log("Hello, world!")
   */
  log: (...data: any[]) => undefined;
  /**
   * Clears the console
   * @returns undefined
   * @example console.clear()
   */
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
declare module "lar:http" {
  const serve: (obj: ServeObject) => unknown;
  const get: (routeName: string, callback: () => Awaitable<any>) => any;
  const fetch: (url: string) => Awaitable<any>;
}

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
