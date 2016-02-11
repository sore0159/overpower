(function() {
    var canvas = document.getElementById('mainscreen');
    var grid = canvas.muleGrid;

    canvas.redrawPage = function() {
        canvas.setupText();
        canvas.drawMap();
    };

    // setupText performs a full refresh of the infobox
    canvas.setupText = function() {
        var data = canvas.overpowerData;
        var scoreText = document.getElementById('scoretext');
        scoreText.textContent = data.faction.score;
        var highScoreText = document.getElementById('highscoretext');
        highScoreText.textContent = data.game.highscore;
        var turnText = document.getElementById('turntext');
        turnText.textContent = data.game.turn;
        var turnCompText = document.getElementById('turncompletetext');
        var turnChangeButton = document.getElementById('turnchange');
        if (data.faction.done) {
            turnCompText.textContent = "Turn Complete";
            turnChangeButton.textContent = "Cancel Turn Complete";
        } else {
            turnCompText.textContent = "Turn In Progress";
            turnChangeButton.textContent = "Set Turn Complete";
        }
        canvas.setupTargets();
    };

    // setupText performs a refresh of the target area of the infobox
    canvas.setupTargets = function() {
        var data = canvas.overpowerData;
        var secTargetText = document.getElementById('secondtargettext');
        var secTargetDiv = document.getElementById('secondtargetbox');
        var swapButton = document.getElementById('swapbutton');
        var dist, elem, button;
        secTargetDiv.textContent = "";
        if (!data.targetTwo) {
            secTargetText.textContent = "None";
            swapButton.style.display = "none";
        } else if (data.targetOne && data.targetOne[0] === data.targetTwo[0] && data.targetOne[1] === data.targetTwo[1]) {
            secTargetText.textContent = "Main";
            swapButton.style.display = "none";
        } else {
            secTargetText.textContent = "("+data.targetTwo[0]+","+data.targetTwo[1]+") ";
            elem = document.createElement("hr");
            elem.className = "target";
            secTargetDiv.appendChild(elem);
            button = planetButton(data.targetTwoInfo.planet);
            secTargetDiv.appendChild(button);
            if (data.targetOne && data.targetOneInfo.planet) {
                elem = document.createElement("br");
                secTargetDiv.appendChild(elem);
                dist = grid.stepsBetween(data.targetOne, data.targetTwo);
                elem = document.createTextNode(dist+" sectors from main target");
                secTargetDiv.appendChild(elem);
                swapButton.style.display = "inline";
            } else {
                swapButton.style.display = "none";
            }
           
        }
        var mainTargetText = document.getElementById('maintargettext');
        var mainTargetBox = document.getElementById('maintargetbox');
        var planetInfoBox = document.getElementById('planetinfobox');
        var orderInfoBox = document.getElementById('orderinfobox');
        var shipInfoBox = document.getElementById('shipinfobox');
        var trailInfoBox = document.getElementById('trailinfobox');
        var htmlStr;
                function hello() {
                    console.log("hELLO");
                }
        if (data.targetOne) {
            mainTargetText.textContent = "("+data.targetOne[0]+","+data.targetOne[1]+")";
            mainTargetBox.style.display="block";
            var planet = data.targetOneInfo.planet;
            planetInfoBox.innerHTML = "";
            if (planet) {
                elem = document.createElement("hr");
                elem.className = "target";
                planetInfoBox.appendChild(elem);
                button = planetButton(planet);
                planetInfoBox.appendChild(button);
                elem = document.createElement("br");
                planetInfoBox.appendChild(elem);
                if (!planet.turn) {
                    elem = document.createTextNode("No planetary information available");
                    planetInfoBox.appendChild(elem);
                } else {
                    if (planet.controller === planet.fid) {
                        htmlStr = "Your planet";
                    } else {
                        if (planet.controller) {
                            var name = data.fidMap.get(planet.controller).name;
                            htmlStr = name + " planet";
                        } else {
                            htmlStr = "Nuetral planet";
                        }
                        htmlStr += " \u2022 Last seen on turn "+planet.turn;
                    }
                    elem = document.createTextNode(htmlStr);
                    planetInfoBox.appendChild(elem);
                    elem = document.createElement("br");
                    planetInfoBox.appendChild(elem);

                    elem = document.createElement("b");
                    elem.textContent = "Inhabitants: ";
                    planetInfoBox.appendChild(elem);
                    elem = document.createTextNode(planet.inhabitants + " \u2022 ");
                    planetInfoBox.appendChild(elem);
                    elem = document.createElement("b");
                    elem.textContent = "Resources: ";
                    planetInfoBox.appendChild(elem);
                    elem = document.createTextNode(planet.resources+ " \u2022 ");
                    planetInfoBox.appendChild(elem);

                    elem = document.createElement("b");
                    elem.textContent = "Parts: ";
                    planetInfoBox.appendChild(elem);
                    htmlStr = ""+planet.parts;
                    if (planet.avail) {
                        htmlStr += " ("+planet.avail+" available)";
                    }
                    elem = document.createTextNode(htmlStr);
                    planetInfoBox.appendChild(elem);
                }
            }
            orderInfoBox.innerHTML = "";
            if (data.targetOneInfo.orders && data.targetOneInfo.orders.length > 0 && (!data.targetOrder || data.targetOneInfo.orders.length > 1 || data.targetOrder !== data.targetOneInfo.orders[0])) {
                elem = document.createElement("hr");
                elem.className = "target";
                orderInfoBox.appendChild(elem);
                elem = document.createElement("b");
                elem.textContent = "Launch orders from "+planet.name+":";
                orderInfoBox.appendChild(elem);
                elem = document.createElement("br");
                orderInfoBox.appendChild(elem);
                var uList = document.createElement("ul");
                orderInfoBox.appendChild(uList);
                for (var i = 0; i < data.targetOneInfo.orders.length; i++) {
                    if (data.targetOneInfo.orders[i] === data.targetOrder) {
                        continue;
                    }
                    var lItem = document.createElement("li");
                    uList.appendChild(lItem);
                    elem = document.createElement("b");
                    elem.textContent = "Target: ";
                    lItem.appendChild(elem);
                    button = planetButton(data.targetOneInfo.orders[i].targetPl);
                    lItem.appendChild(button);
                    elem = document.createElement("b");
                    elem.innerHTML = " &bull; Size: ";
                    lItem.appendChild(elem);
                    elem = document.createTextNode(data.targetOneInfo.orders[i].size);
                    lItem.appendChild(elem);
                }
            }
            if (data.targetOneInfo.ships.length > 0 ) {
                htmlStr = "<hr class=\"target\">";
                htmlStr += "<b>Ship";
                if (data.targetOneInfo.ships.length > 1) {
                    htmlStr += "s";
                }
                htmlStr += " detected:</b><br><ul>";
                data.targetOneInfo.ships.forEach(function(ship) {
                    htmlStr += "<li>";
                    if (ship.controller == ship.fid) {
                        htmlStr += "Your ship &bull;";
                        htmlStr += " Sized "+ship.size;
                        htmlStr += " &bull; Destination: "+ship.dest.planet.name;
                    } else {
                        htmlStr += data.fidMap.get(ship.controller).name+" ship &bull;";
                        htmlStr += " Sized "+ship.size;
                    }
                });
                htmlStr += "</ul>";
                shipInfoBox.innerHTML = htmlStr;
            } else {
                shipInfoBox.innerHTML = "";
            }
            if (data.targetOneInfo.trails.length > 0 ) {
                htmlStr = "<hr class=\"target\">";
                htmlStr += "<b>Ship passage";
                if (data.targetOneInfo.trails.length > 1) {
                    htmlStr += "s";
                }
                htmlStr += " detected:</b><br><ul>";
                data.targetOneInfo.trails.forEach(function(ship) {
                    htmlStr += "<li>";
                    if (ship.controller == ship.fid) {
                        htmlStr += "Your ship &bull;";
                        htmlStr += " Sized "+ship.size;
                        if (ship.loc.valid) {
                            htmlStr += " &bull; Currently at: ("+ship.loc.coord[0]+", "+ship.loc.coord[1]+")";
                        } else {
                            htmlStr += " &bull; Landed on: "+ship.dest.planet.name;
                        }
                    } else {
                        htmlStr += data.fidMap.get(ship.controller).name+" ship &bull;";
                        htmlStr += " Sized "+ship.size;
                        if (ship.loc.valid) {
                            htmlStr += " &bull; Currently at: ("+ship.loc.coord[0]+", "+ship.loc.coord[1]+")";
                        } else {
                            htmlStr += " &bull; Location unknown";
                        }
                    }
                });
                htmlStr += "</ul>";
                trailInfoBox.innerHTML = htmlStr;
            } else {
                trailInfoBox.innerHTML = "";
            }
        } else {
            mainTargetText.textContent = "None";
            mainTargetBox.style.display="none";
        }
        var targetOrderDiv = document.getElementById('targetorderbox');
        var orderConfirmButton = document.getElementById('orderconfirm');
        if (data.targetOrder) {
            var order = data.targetOrder;
            var changedOrder = ((data.targetOrder.brandnew && data.targetOrder.size !== 0) || (!data.targetOrder.brandnew && data.targetOrder.originSize != data.targetOrder.size)) ;
            htmlStr = "<hr class=\"target\">";
            dist = grid.stepsBetween(order.sourcePl.loc, order.targetPl.loc);
            dist = Math.ceil(dist/10);
            if (order.brandnew) {
                htmlStr +=  "Launch from <b>"+order.sourcePl.name + "</b> to <b>"+ order.targetPl.name+"</b> ("+dist+" turns)?<br> &nbsp;&nbsp;&nbsp;&nbsp;&bull; ";
                if (changedOrder) {
                    htmlStr += "<b>Launch ship sized: "+order.size+"</b>";
                } else {
                    htmlStr += "Launch ship sized: "+order.size;
                }
            } else {
                htmlStr +=  "Modifiy launch order from <b>"+order.sourcePl.name + "</b> to <b>"+ order.targetPl.name+"</b> (currently sized "+order.originSize+") ("+dist+" turns)?<br> &nbsp;&nbsp;&nbsp;&nbsp;&bull; ";
                if (changedOrder) {
                    htmlStr += "<b>New ship size: "+order.size+"</b>";
                } else {
                    htmlStr += "New ship size: "+order.size;
                }
            }
            if (changedOrder) {
                orderConfirmButton.style.display = "inline";
            } else {
                if (order.sourcePl.avail || order.size > 0) {
                    htmlStr += "<br>(hold shift and scroll the mousewheel to change)";
                }
                orderConfirmButton.style.display = "none";
            }
            targetOrderDiv.innerHTML = htmlStr;
        } else if (data.targetOneInfo && data.targetOneInfo.planet && data.targetOneInfo.planet.avail > 0)  {
            targetOrderDiv.innerHTML = "<hr class=\"target\">"+
                "To launch a spaceship from this planet, right click another planet to set it as your secondary target";
            orderConfirmButton.style.display = "none";
        } else {
            targetOrderDiv.textContent = "";
            orderConfirmButton.style.display = "none";
        }
    };

    canvas.blockScreen = function(text, buttonT) {
        var screen = document.getElementById('fullscreenblock');
        var blocktext = document.getElementById('blockertext');
        var blockbutton = document.getElementById('blockbutton');
        blocktext.innerHTML = text;
        if (buttonT) {
            blockbutton.textContent = buttonT;
            blockbutton.style.display = "inline";
        } else {
            blockbutton.style.display = "none";
        }
        screen.style.display = 'block';

    };

    canvas.unblockScreen = function() {
        var screen = document.getElementById('fullscreenblock');
        screen.style.display = 'none';
    };

    function planetButton(planet) {
        return hexButton(planet.loc, planet.name);
    }
    function hexButton(hex, text) {
        var hexB = document.createElement("button");
        hexB.className = "hexbutton";
        function hexClick(event) {
            if (event.button === 0 && !event.shiftKey) {
                canvas.overpowerData.setTargetOne(hex);
                canvas.redrawPage();
            } else if (event.button === 2) {
                canvas.overpowerData.setTargetTwo(hex);
                canvas.redrawPage();
            } else if (event.button === 1 || event.button === 0) {
                canvas.setCenterDest(hex);
            }
            return false;
        }
        hexB.textContent = text;
        hexB.oncontextmenu= function() { return false; };
        hexB.addEventListener("mouseup", hexClick, false);
        return hexB;
    }

})();
