// QUery logic
function getCountdown(callback) {
    var req = new XMLHttpRequest();
    const method = "GET";
    const url = "localhost:8080/countdown";
    req.open(method, url, false)
    req.onreadystatechange = () => {
        if (req.readyState = 4 && req.status < 200) {
            callback(req.response)
        }
        console.log("Failed to fetch data from backend")
    }
    req.send();
}

function renderPage(response) {
    obj = JSON.parse(response);
    console.log(obj);
    let header = document.getElementById("header");
    let countdown = document.getElementById("countdown");
    header.value = "Time to Myrtle: " + String(obj["Year"]);
    countdown.value = obj["Body"]
}

setInterval(getCountdown, 1000)