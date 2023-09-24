// QUery logic
function getCountdown(callback) {
    var req = new XMLHttpRequest();
    const method = "GET";
    const url = window.location.href + "countdown";
    console.log(url)
    req.open(method, url, false);
    req.onreadystatechange = () => {
        if (req.readyState == 4 && req.status == 200) {
            callback(req.response);
        } else {
            console.log("Failed to fetch data from backend");
        }
    }
    req.send();
}

function renderPage(response) {
    obj = JSON.parse(response);
    console.log(obj);
    let countdown = document.getElementById("countdown");
    countdown.innerHTML = obj["Body"];
}

setInterval(() => { getCountdown(renderPage) }, 1000);