declare const require: (module: string) => any;
declare const exports: any;
interface ReadFileObject {
    encoding: 'utf8' | 'utf-8' | 'base64'
}
type ColorType = "red" | "green" | "blue" | "black" | "cyan" | "yellow" | "magenta" | "white"
type LarObject = { color: (colorType: ColorType, ...data: any[]) => undefined }
type ConsoleObject = { time: (label?: string) => undefined, timeEnd: (label?: string) => undefined, log: (...data: any[]) => undefined, clear: () => undefined, error: (message: string) => undefined }
type FSObject = { readFile: (fileName: string, obj?: ReadFileObject) => string }
declare const Lar: LarObject
declare const console: ConsoleObject
declare const fs: FSObject
declare const Add: (a: number, b: number) => number
declare const Mult: (a: number, b: number) => number
declare const Div: (a: number, b: number) => number
declare const print: (...data: any[]) => undefined
declare const prompt: (question: string) => string
