(function() {

var canvas = document.getElementById('mainscreen');

canvas.fetchData = function(dataSpecs, onPass, onFail) {
    var req = new XMLHttpRequest();
    function handleReq() {
        var err, data;
        if (req.status === 200) {
            resp = JSON.parse(req.responseText);
            if (resp.status === "success") {
                onPass(resp.data);
                return;
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
        onFail(err);
    }
    req.addEventListener("load", handleReq);
    req.addEventListener("error", function() {
        onFail("an error occured with the request");
    });
    req.open("GET", "/overpower/json/"+dataSpecs, true);
    req.send();
};

canvas.refetchFullView = function() {
    canvas.fetchFullView(this.overpowerData.game.gid, this.overpowerData.faction.fid, true);
};

canvas.fetchFullView = function(gid, fid, blockcycle) {
    canvas.blockScreen("Loading game data...");
    dataSpecs = "fullviews/"+gid+"/"+fid;
    function onPass(data) {
        canvas.stopBlink();
        canvas.overpowerData = data;
        canvas.parseOPData();
        canvas.redrawPage();
        canvas.unblockScreen();
        if (!blockcycle) {
            turnCheckCycle();
        }
    }
    function onFail(err) {
        console.log("ERROR FETCHING DATA:", err, "\n RETRYING...");
        canvas.blockScreen("Loading game data...<br>Server not responding, retrying...");

        window.setTimeout(function() {
            canvas.fetchFullView(gid, fid);
        }, 3000);
    }
    canvas.fetchData(dataSpecs, onPass, onFail);
};


var icon = document.querySelector('link[rel="shortcut icon"]');

function turnCheckCycle() {
    canvas.turnCheck(true);
}

canvas.turnCheck = function(cycle) {
    function onPass(game) {
        var turn = game.turn;
        if (game.turn > canvas.overpowerData.game.turn) {
            canvas.blockScreen("There is a new turn available", "Click here to fetch the new turn");
            stopblink = false;
            blinkIco();
        } else if (cycle){
            window.setTimeout(turnCheckCycle, 15000);
        }
    }
    function onFail(err) {
        console.log("ERROR FETCHING GAME TURN DATA:", err);
    }
    var dataSpecs = "games/"+canvas.overpowerData.game.gid;
    canvas.fetchData(dataSpecs, onPass, onFail);
};

var blink = false;
var stopblink = false;

canvas.stopBlink = function() {
    stopblink = true;
};

function blinkIco() {
    if (blink) {
        icon.href = "/static/img/yd32.ico";
    } else {
        icon.href = "/static/img/yd32blink.ico";
    }
    blink = !blink;
    if (stopblink) {
        icon.href = "/static/img/yd32.ico";
    } else {
        window.setTimeout(blinkIco, 1000);
    }
}

})();
