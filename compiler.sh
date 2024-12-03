#!/bin/bash
go build -o lar
./lar run time.ts
./lar run readFile.ts
./lar run count.ts
