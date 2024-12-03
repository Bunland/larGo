get("/hello", () => {
    console.log("X")
    return "Hello!"
})

serve({ port: "3000" })