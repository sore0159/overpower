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
        var boxTurnBuffer = document.getElementById('bufferbox');
        // ---------- TURN INFO STUFF ------------ //
        boxTurnBuffer.redraw();
        // ---------- REPORT INFO STUFF ------------ //
        var reportText = document.getElementById('reporttext');
        var htmlStr;
        if (data.launchrecords.length === 0) {
            htmlStr = "No launch reports";
        } else if (data.launchrecords.length === 1) {
            htmlStr = "One launch report";
        } else {
            htmlStr = ""+data.launchrecords.length+" launch reports";
        }
        htmlStr += " \u2022 ";
        if (data.landingrecords.length === 0) {
            htmlStr += "No landing reports";
        } else if (data.landingrecords.length === 1) {
            htmlStr += "One landing report";
        } else {
            htmlStr += ""+data.landingrecords.length+" landing reports";
        }
      
        reportText.textContent = htmlStr;
        canvas.setupTargets();
    };

    // setupText performs a refresh of the target area of the infobox
    canvas.setupTargets = function() {
        var data = canvas.overpowerData;

        var secTargetDiv = document.getElementById('secondtargetbox');
        var planetInfoBox = document.getElementById('planetinfobox');
        var orderInfoBox = document.getElementById('orderinfobox');
        var shipInfoBox = document.getElementById('shipinfobox');
        var trailInfoBox = document.getElementById('trailinfobox');
        var targetOrderDiv = document.getElementById('targetorderbox');

        var divList = [secTargetDiv, planetInfoBox,
           orderInfoBox, shipInfoBox, trailInfoBox, targetOrderDiv];
        for (var i = 0; i < divList.length; i++) {
            while (divList[i].firstChild) {
                divList[i].removeChild(divList[i].firstChild);
            }
        }
       
        var secTargetText = document.getElementById('secondtargettext');
        var mainTargetText = document.getElementById('maintargettext');
        var mainTargetBox = document.getElementById('maintargetbox');
        var swapButton = document.getElementById('swapbutton');
        var dist, elem, button;
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
            if (data.targetOne) {
                elem = document.createElement("br");
                secTargetDiv.appendChild(elem);
                dist = grid.stepsBetween(data.targetOne, data.targetTwo);
                elem = document.createTextNode(dist+" sectors from main target");
                secTargetDiv.appendChild(elem);
                if (data.targetOneInfo.planet) {
                    swapButton.style.display = "inline";
                }
            } else {
                swapButton.style.display = "none";
            }
        }
        var htmlStr;
        if (data.targetOne) {
            mainTargetText.textContent = "("+data.targetOne[0]+","+data.targetOne[1]+")";
            mainTargetBox.style.display="block";
            var planet = data.targetOneInfo.planet;
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
                            htmlStr = "Neutral planet";
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
            var uList, lItem;
            if (data.targetOneInfo.orders && data.targetOneInfo.orders.length > 0 && (!data.targetOrder || data.targetOneInfo.orders.length > 1 || data.targetOrder !== data.targetOneInfo.orders[0])) {
                elem = document.createElement("hr");
                elem.className = "target";
                orderInfoBox.appendChild(elem);
                elem = document.createElement("b");
                elem.textContent = "Launch orders";
                orderInfoBox.appendChild(elem);
                elem = document.createTextNode(" from "+planet.name+":");
                orderInfoBox.appendChild(elem);
                elem = document.createElement("br");
                orderInfoBox.appendChild(elem);
                uList = document.createElement("ul");
                orderInfoBox.appendChild(uList);
                for (i = 0; i < data.targetOneInfo.orders.length; i++) {
                    if (data.targetOneInfo.orders[i] === data.targetOrder) {
                        continue;
                    }
                    lItem = document.createElement("li");
                    uList.appendChild(lItem);
                    elem = document.createElement("b");
                    elem.textContent = "Target: ";
                    lItem.appendChild(elem);
                    button = planetButton(data.targetOneInfo.orders[i].targetPl);
                    lItem.appendChild(button);
                    elem = document.createElement("b");
                    elem.textContent = " \u2022 Size: ";
                    lItem.appendChild(elem);
                    elem = document.createTextNode(data.targetOneInfo.orders[i].size);
                    lItem.appendChild(elem);
                }
            }
            if (data.targetOneInfo.ships.length > 0 ) {
                elem = document.createElement("hr");
                elem.className = "target";
                shipInfoBox.appendChild(elem);
                htmlStr = "Ship";
                if (data.targetOneInfo.ships.length > 1) {
                    htmlStr += "s";
                }
                htmlStr += " detected:";
                elem = document.createElement("b");
                elem.textContent = htmlStr;
                shipInfoBox.appendChild(elem);
                elem = document.createElement("br");
                shipInfoBox.appendChild(elem);
                uList = document.createElement("ul");
                shipInfoBox.appendChild(uList);
                data.targetOneInfo.ships.forEach(function(ship) {
                    lItem = document.createElement("li");
                    uList.appendChild(lItem);
                    if (ship.controller == ship.fid) {
                        elem = document.createTextNode("Your ship \u2022 Sized "+ship.size+ " \u2022 Destination: ");
                        lItem.appendChild(elem);
                        button = planetButton(ship.dest.planet);
                        lItem.appendChild(button);
                    } else {
                        elem = document.createTextNode(data.fidMap.get(ship.controller).name+" ship \u2022 Sized "+ship.size);
                        lItem.appendChild(elem);
                    }
                    if (ship.trail.length > 0) {
                        elem = document.createElement("br");
                        lItem.appendChild(elem);
                        elem = document.createTextNode("Travelled through:");
                        lItem.appendChild(elem);
                        ship.trail.forEach(function(trailHex) {
                            button = coordButton(trailHex);
                            lItem.appendChild(button);
                        });
                    }
              
                });
            }
            if (data.targetOneInfo.trails.length > 0 ) {
                elem = document.createElement("hr");
                elem.className = "target";
                trailInfoBox.appendChild(elem);
                htmlStr = "Ship passage";
                if (data.targetOneInfo.ships.length > 1) {
                    htmlStr += "s";
                }
                htmlStr += " detected:";
                elem = document.createElement("b");
                elem.textContent = htmlStr;
                trailInfoBox.appendChild(elem);
                elem = document.createElement("br");
                trailInfoBox.appendChild(elem);
                uList = document.createElement("ul");
                trailInfoBox.appendChild(uList);
                data.targetOneInfo.trails.forEach(function(ship) {
                    lItem = document.createElement("li");
                    uList.appendChild(lItem);
                    if (ship.controller == ship.fid) {
                        htmlStr = "Your ship \u2022 Sized "+ship.size+" \u2022 ";
                        if (ship.loc.valid) {
                            htmlStr += "Currently at: ";
                            elem = document.createTextNode(htmlStr);
                            lItem.appendChild(elem);
                            button = coordButton(ship.loc.coord);
                            lItem.appendChild(button);
                        } else {
                            htmlStr += "Landed on: ";
                            elem = document.createTextNode(htmlStr);
                            lItem.appendChild(elem);
                            button = planetButton(ship.dest.planet);
                            lItem.appendChild(button);
                        }
                    } else {
                        htmlStr = data.fidMap.get(ship.controller).name+" ship \u2022 Sized "+ship.size+ " \u2022 ";
                        if (ship.loc.valid) {
                            htmlStr += "Currently at: ";
                            elem = document.createTextNode(htmlStr);
                            lItem.appendChild(elem);
                            button = coordButton(ship.loc.coord);
                            lItem.appendChild(button);
                        } else {
                            htmlStr += "Location unknown";
                            elem = document.createTextNode(htmlStr);
                            lItem.appendChild(elem);
                        }
                    }
                    elem = document.createElement("br");
                    lItem.appendChild(elem);
                    elem = document.createTextNode("Travelled through:");
                    lItem.appendChild(elem);
                    ship.trail.forEach(function(trailHex) {
                        button = coordButton(trailHex);
                        lItem.appendChild(button);
                    });
                });
            }
            if (!planet && data.targetOneInfo.ships.length === 0 && data.targetOneInfo.trails.length === 0 ) {
                elem = document.createElement("hr");
                elem.className = "target";
                planetInfoBox.appendChild(elem);
                elem = document.createElement("b");
                elem.textContent = "The Great Void";
                planetInfoBox.appendChild(elem);
            }
        } else {
            mainTargetText.textContent = "None";
            mainTargetBox.style.display="none";
        }
        var orderConfirmButton = document.getElementById('orderconfirm');
        if (data.targetOrder) {
            elem = document.createElement("hr");
            elem.className = "target";
            targetOrderDiv.appendChild(elem);
            var order = data.targetOrder;
            var changedOrder = ((data.targetOrder.brandnew && data.targetOrder.size !== 0) || (!data.targetOrder.brandnew && data.targetOrder.originSize != data.targetOrder.size)) ;
            dist = grid.stepsBetween(order.sourcePl.loc, order.targetPl.loc);
            dist = Math.ceil(dist/10);
            if (order.brandnew) {
                elem = document.createTextNode("Launch from ");
                targetOrderDiv.appendChild(elem);
                htmlStr = "New ship sized: "+order.size;
            } else {
                elem = document.createTextNode("Modify launch order from ");
                targetOrderDiv.appendChild(elem);
                htmlStr = "Launch ship sized: "+order.size;
            }
            elem = document.createElement("b");
            elem.textContent = order.sourcePl.name;
            targetOrderDiv.appendChild(elem);
            elem = document.createTextNode(" to ");
            targetOrderDiv.appendChild(elem);
            elem = document.createElement("b");
            elem.textContent = order.targetPl.name;
            targetOrderDiv.appendChild(elem);
            elem = document.createTextNode(" ("+dist+" turns)?");
            targetOrderDiv.appendChild(elem);
            elem = document.createElement("br");
            targetOrderDiv.appendChild(elem);
            elem = document.createTextNode("\u00A0\u00A0\u00A0\u00A0\u2022 ");
            targetOrderDiv.appendChild(elem);
            if (changedOrder) {
                elem = document.createElement("b");
                elem.textContent = htmlStr;
            } else {
                elem = document.createTextNode(htmlStr);
            }
            targetOrderDiv.appendChild(elem);
            if (changedOrder) {
                orderConfirmButton.style.display = "inline";
            } else {
                if (order.sourcePl.avail || order.size > 0) {
                    elem = document.createElement("br");
                    targetOrderDiv.appendChild(elem);
                    elem = document.createTextNode("(hold shift and scroll the mousewheel to change)");
                    targetOrderDiv.appendChild(elem);
                }
                orderConfirmButton.style.display = "none";
            }
        } else if (data.targetOneInfo && data.targetOneInfo.planet && data.targetOneInfo.planet.avail > 0)  {
            elem = document.createElement("hr");
            elem.className = "target";
            targetOrderDiv.appendChild(elem);
            elem = document.createTextNode("To launch a spaceship from this planet, right click another planet to set it as your secondary target");
            targetOrderDiv.appendChild(elem);
            orderConfirmButton.style.display = "none";
        } else {
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
        screen.style.opacity = 0;
        screen.style.display = 'block';
        window.setTimeout(function() {
            screen.style.opacity = 0.75;
        }, 100);

    };

    canvas.unblockScreen = function() {
        var screen = document.getElementById('fullscreenblock');
        screen.style.display = 'none';
    };

    function coordButton(hex) {
        return hexButton(hex, "("+hex[0]+","+hex[1]+")");
    }
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
    canvas.planetButton = planetButton;
    canvas.coordButton = planetButton;

})();
