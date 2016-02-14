(function() {

var canvas = document.getElementById('mainscreen');
var reportButton = document.getElementById('reportbutton');
var reportChangeButton = document.getElementById('reportchangeturn');
var closeButton = document.getElementById('reportclose');
var reportScreen0 = document.getElementById('reportscreen0');
var reportScreen1 = document.getElementById('reportscreen1');
var reportDisplay = document.getElementById('reportdisplaybody');

reportButton.addEventListener("mouseup", openReports);
reportScreen0.addEventListener("mouseup", closeReports);
closeButton.addEventListener("mouseup", closeReports);
reportChangeButton.addEventListener("mouseup", changeClick);
reportChangeButton.addEventListener("DOMMouseScroll", turnScroll);
reportChangeButton.onmousewheel = turnScroll;

reportChangeButton.setTurn = function(turn) {
    this.turn = turn;
    this.clickTurn = turn;
    this.redraw();
};
reportChangeButton.redraw = function() {
    if (this.turn === this.clickTurn) {
        this.textContent = "Mousescroll over this button to select new turn";
        this.style.fontWeight = "normal";
    } else {
        this.textContent = "Click to view reports for turn "+this.clickTurn;
        this.style.fontWeight = "bold";
    }
};
function changeClick(event) {
    if (!(this.clickTurn && this.turn && this.turn != this.clickTurn)) {
        return;
    }
    var config = canvas.overpowerData.reportConfig;
    config.turn = this.clickTurn;
    getRecords();
}

function turnScroll(event) {
    event.preventDefault();
    var gameturn = canvas.overpowerData.game.turn;
    var up;
    if (event.detail) {
        up = -1*event.detail/3;
    } else {
        up = (event.wheelDelta)/120;
    }
    if (up > 0 && up < 1) {
        up = 1;
    } else if (up < 0 && up > -1) {
        up = -1;
    } else if (up === 0) {
        return true;
    }
    if (up === 1 && this.clickTurn < gameturn - 1) {
        this.clickTurn += 1;
    } else if (up === -1 && this.clickTurn > 1) {
        this.clickTurn -= 1;
    } else {
        return false;
    }
    this.redraw();
    return false;
}

function openReports() {
    // ---------- SET UP --------------- //
    var data = canvas.overpowerData;
    if (!data.reportConfig) {
        data.reportConfig = {
            turn:data.game.turn-1,
            filters: {}
        };
        getRecords();
    } else {
        gotRecords();
    }
}

function getRecords() {
    var data = canvas.overpowerData;
    var config = data.reportConfig;
    if (config.turn === data.game.turn-1) {
        config.launch = data.launchrecords;
        config.landing = data.landingrecords;
        parseRecords();
        gotRecords();
    } else {
        canvas.fetchData("launchrecords/"+data.game.gid+"/"+data.faction.fid+"/"+config.turn, getLAPass, getFail);
    }
}

function getFail(err) {
    console.log("ERROR FETCHING GAME RECORDS DATA:", err);
    canvas.blockScreen("Error retrieving data from server", "Click to reconnect");
}
function getLAPass(data) {
    var opData = canvas.overpowerData;
    var config = canvas.overpowerData.reportConfig;
    config.launch = data;
    canvas.fetchData("landingrecords/"+opData.game.gid+"/"+opData.faction.fid+"/"+config.turn, function(data) {
        var config = canvas.overpowerData.reportConfig;
        config.landing = data;
        parseRecords();
        gotRecords();
    }, getFail);
}

function parseRecords() {
    var opData = canvas.overpowerData;
    var data = opData.reportConfig;
    reportChangeButton.setTurn(data.turn);
    while (reportDisplay.firstChild) {
        reportDisplay.removeChild(reportDisplay.firstChild);
    }
    var turnText = document.getElementById('reportturntext');
    turnText.textContent = data.turn;
    var elem, uList, lItem, button;
    if (data.launch.length) {
        elem = document.createElement("b");
        elem.textContent = "Launch reports:";
        reportDisplay.appendChild(elem);
        uList = document.createElement("ul");
        reportDisplay.appendChild(uList);
        data.launch.forEach(function(launch) {
            lItem = document.createElement("li");
            uList.appendChild(lItem);
            button = planetButton(launch.source);
            lItem.appendChild(button);
            elem = document.createTextNode(" launched a size "+launch.size+" ship toward ");
            lItem.appendChild(elem);
            button = planetButton(launch.target);
            lItem.appendChild(button);
        });
    } else {
        elem = document.createElement("b");
        elem.textContent = "No launch reports";
        reportDisplay.appendChild(elem);
        elem = document.createElement("br");
        reportDisplay.appendChild(elem);
    }
    elem = document.createElement("hr");
    reportDisplay.appendChild(elem);
    if (data.landing) {
        elem = document.createElement("b");
        elem.textContent = "Landing reports:";
        reportDisplay.appendChild(elem);
        uList = document.createElement("ul");
        reportDisplay.appendChild(uList);
        data.landing.forEach(function(landing) {
            lItem = document.createElement("li");
            uList.appendChild(lItem);
            button = planetButton(landing.target);
            lItem.appendChild(button);
            addText(lItem, " \u2022 Controlled by ");
            if (landing.firstcontroller === 0)  {
            } else if (landing.firstcontroller === opData.faction.fid) {
            } else {
                var fname = opData.fidMap.get(landing.firstcontroller).name;
            }
            addText(lItem, JSON.stringify(landing));
        });
    } else {
        elem = document.createElement("b");
        elem.textContent = "No landing reports";
        reportDisplay.appendChild(elem);
        elem = document.createElement("br");
        reportDisplay.appendChild(elem);
    }

}

function gotRecords() {
    // ---------- TURN ON --------------- //
    reportScreen0.style.display="inline";
    reportScreen1.style.display="inline";
    window.setTimeout(function() {
        reportScreen0.style.opacity = 0.75;
        reportScreen1.style.opacity = 1;
    }, 100);
}

function closeReports() {
    reportScreen0.style.opacity = 0;
    reportScreen1.style.opacity = 0;
    window.setTimeout(function() {
        reportScreen0.style.display="none";
        reportScreen1.style.display="none";
    }, 750);
}

function planetButton(hex) {
    var planet = canvas.overpowerData.gridMap.getHex(hex);
    var button = canvas.planetButton(planet);
    button.addEventListener("mouseup", closeReports, false);
    return button;
}

function addText(element, text) {
    element.appendChild(document.createTextNode(text));
}


})();
