(function() {
var newTurn = false;
var blink = false;
var icon = document.querySelector('link[rel="shortcut icon"]');
var canvas = document.getElementById('mainscreen');

function log(msg) {
    console.log("Overpower: "+msg);
}

function curGame() {
    return canvas.overpowerData.game.gid;
}
function curTurn() {
    return canvas.overpowerData.game.turn;
}

function getTurn(cb) {
    var req = new XMLHttpRequest();
    function handleTurn() {
        var err, turn;
        if (req.status === 200) {
            resp = JSON.parse(req.responseText);
            if (resp.status === "success") {
                turn = resp.data.turn;
            } else if (resp.status === "fail") {
                err = "json fail response from server:"+'\n'+resp.data;
            } else if (resp.status === "error") {
                err = "json error response from server:"+'\n'+resp.message;
            } else {
                err = "unknown json response from server";
            }
        } else {
            err = req.status +'\n'+req.statusText;
        }
        cb(err, turn);
    }
    req.addEventListener("load", handleTurn);
    req.addEventListener("error", function() {
        cb("an error occured with the request");
    });
    req.open("GET", "/overpower/json/games/"+curGame(), true);
    req.send();
}

function turnCheck(isNew) {
    if (newTurn) {
        return isNew(true);
    }
    function gotTurn(err, turn) {
        if (err) {
            log("error fetching turn:", err);
            return isNew(false);
        }
        if (turn > curTurn()) {
            newTurn = true;
        }
        return isNew(turn > curTurn());
    }
    getTurn(gotTurn);
}

function whenNew(newCheck) {
    if (newCheck) {
        var blocker = document.querySelector("div.blocker");
        if (blocker) {
            blocker.style.display = 'block';
        }
        blinkIco();
    } else {
        window.setTimeout(cycle, 15000);
    }
}

function blinkIco() {
    if (blink) {
        icon.href = "/static/img/yd32.ico";
    } else {
        icon.href = "/static/img/yd32blink.ico";
    }
    blink = !blink;
    window.setTimeout(blinkIco, 1000);
}

function cycle() {
    turnCheck(whenNew);
}

cycle();
})();


