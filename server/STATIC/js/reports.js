(function() {

var canvas = document.getElementById('mainscreen');
var reportButton = document.getElementById('reportbutton');
var reportChangeButton = document.getElementById('reportchangeturn');
var closeButton = document.getElementById('reportclose');
var reportScreen0 = document.getElementById('reportscreen0');
var reportScreen1 = document.getElementById('reportscreen1');
var reportDisplayLA = document.getElementById('launchbody');
var reportDisplayLD = document.getElementById('landingbody');
var filterAll = document.getElementById('reportfilterall');
var filterLA = document.getElementById('reportfilterlaunch');
var filterLD = document.getElementById('reportfilterlanding');

reportButton.addEventListener("mouseup", openReports);
reportScreen0.addEventListener("mouseup", closeReports);
closeButton.addEventListener("mouseup", closeReports);
reportChangeButton.addEventListener("mouseup", changeClick);
reportChangeButton.addEventListener("DOMMouseScroll", turnScroll);
reportChangeButton.onmousewheel = turnScroll;
filterAll.addEventListener("mouseup", setFilterAll);
filterLA.addEventListener("mouseup", setFilterLA);
filterLD.addEventListener("mouseup", setFilterLD);

reportChangeButton.setTurn = function(turn) {
    this.turn = turn;
    this.clickTurn = turn;
    this.redraw();
};
function setFilterAll() {
    var config = canvas.overpowerData.reportConfig;
    config.filter = {};
    parseRecords();
}
function setFilterLA() {
    var config = canvas.overpowerData.reportConfig;
    config.filter = {landing: true};
    parseRecords();
}
function setFilterLD() {
    var config = canvas.overpowerData.reportConfig;
    config.filter = {launch: true};
    parseRecords();
}
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
            filter: {}
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
    while (reportDisplayLA.firstChild) {
        reportDisplayLA.removeChild(reportDisplayLA.firstChild);
    }
    while (reportDisplayLD.firstChild) {
        reportDisplayLD.removeChild(reportDisplayLD.firstChild);
    }
    var turnText = document.getElementById('reportturntext');
    turnText.textContent = data.turn;
    var elem, uList, lItem, button, htmlStr;
    if (!data.filter.launch && !data.filter.landing) {
        filterAll.style.fontWeight = "normal";
        filterAll.style.background = "#ffffff";
        filterAll.style.borderWidth = "2px";
        filterAll.textContent = "Viewing all reports";
        filterLA.style.fontWeight = "bold";
        filterLA.style.background = "#aaccff";
        filterLA.style.borderWidth = "0px";
        filterLA.textContent = "Click to view only launch reports";
        filterLD.style.fontWeight = "bold";
        filterLD.style.background = "#aaccff";
        filterLD.style.borderWidth = "0px";
        filterLD.textContent = "Click to view only landing reports";
    } else if (data.filter.launch) {
        filterAll.style.fontWeight = "bold";
        filterAll.style.background = "#aaccff";
        filterAll.style.borderWidth = "0px";
        filterAll.textContent = "Click to view all reports";
        filterLA.style.fontWeight = "bold";
        filterLA.style.background = "#aaccff";
        filterLA.style.borderWidth = "0px";
        filterLA.textContent = "Click to view only launch reports";
        filterLD.style.fontWeight = "normal";
        filterLD.style.background = "#ffffff";
        filterLD.style.borderWidth = "2px";
        filterLD.textContent = "Viewing only landing reports";
    } else {
        filterAll.style.fontWeight = "bold";
        filterAll.style.background = "#aaccff";
        filterAll.style.borderWidth = "0px";
        filterAll.textContent = "Click to view all reports";
        filterLA.style.fontWeight = "normal";
        filterLA.style.background = "#ffffff";
        filterLA.style.borderWidth = "2px";
        filterLA.textContent = "Viewing only launch reports";
        filterLD.style.fontWeight = "bold";
        filterLD.style.background = "#aaccff";
        filterLD.style.borderWidth = "0px";
        filterLD.textContent = "Click to view only landing reports";
    }
    if (data.filter.launch) {
    } else if (data.launch.length) {
        elem = document.createElement("hr");
        reportDisplayLA.appendChild(elem);
        elem = document.createElement("b");
        elem.textContent = "Launch reports:";
        reportDisplayLA.appendChild(elem);
        uList = document.createElement("ul");
        reportDisplayLA.appendChild(uList);
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
        elem = document.createElement("hr");
        reportDisplayLA.appendChild(elem);
        elem = document.createElement("b");
        elem.textContent = "No launch reports";
        reportDisplayLA.appendChild(elem);
    }
    if (data.filter.landing) {
    } else if (data.landing.length) {
        elem = document.createElement("hr");
        reportDisplayLD.appendChild(elem);
        elem = document.createElement("b");
        elem.textContent = "Landing reports:";
        reportDisplayLD.appendChild(elem);
        uList = document.createElement("ul");
        reportDisplayLD.appendChild(uList);
        data.landing.forEach(function(landing) {
            lItem = document.createElement("li");
            uList.appendChild(lItem);
            button = planetButton(landing.target);
            lItem.appendChild(button);
            htmlStr = ", controlled by ";
            if (landing.firstcontroller === 0)  {
                htmlStr += "hostile natives";
            } else if (landing.firstcontroller === opData.faction.fid) {
                htmlStr += "you";
            } else {
                htmlStr += opData.fidMap.get(landing.firstcontroller).name;
            }
            htmlStr += ", was landed on by a size "+landing.size+" ship of ";
            if (landing.shipcontroller === opData.faction.fid) {
                htmlStr += "yours";
            } else {
                htmlStr += "faction "+opData.fidMap.get(landing.shipcontroller).name;
            }
            htmlStr += ", ending with "+landing.resultinhabitants+" ";
            if (landing.resultcontroller === 0)  {
                htmlStr += "hostile native";
                if (landing.resultinhabitants !== 1) {
                    htmlStr += "s";
                }
            } else if (landing.resultcontroller === opData.faction.fid) {
                htmlStr += "of your colonists";
            } else {
                htmlStr += "of faction "+opData.fidMap.get(landing.resultcontroller).name+" colonists";
            }
            addText(lItem, htmlStr);
        });
    } else {
        elem = document.createElement("hr");
        reportDisplayLD.appendChild(elem);
        elem = document.createElement("b");
        elem.textContent = "No landing reports";
        reportDisplayLD.appendChild(elem);
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
