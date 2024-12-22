import { get, fetch, serve } from "lar:http"

get("/", () => {
    const myFetch = fetch("https://whenisthenextmcufilm.com/api")
    console.log(myFetch)
    return myFetch
})

serve({ port: 3000 })