import { get, fetch, serve, post } from "lar:http"

get("/", () => {
    const myFetch = fetch("https://whenisthenextmcufilm.com/api")
    console.log(myFetch)
    return myFetch
})

post("/examplePost", () => {
    console.log("POST Request")
    return "Example POST Request"
})

serve({ port: 3000 })