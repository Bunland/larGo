import { get, serve, fetch } from "lar:http";

get("/", () => {
    return fetch("https://whenisthenextmcufilm.com/api", {
        method: "GET"
    })
});

serve({ port: 3000 });