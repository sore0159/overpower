(function() {
    var canvas = document.getElementById('mainscreen');
    var grid = canvas.muleGrid;
    var boxTargetOrder = document.getElementById('targetorderbox');
    var boxTargetOne = document.getElementById('targetonebox');
    var boxTargetTwo = document.getElementById('targettwobox');
    var boxReport = document.getElementById('reportbox');
    var boxScore = document.getElementById('scorebox');
    var boxTurn = document.getElementById('turnbox');
    var boxConfirm = document.getElementById('orderconfirm');

    canvas.refreshGameBoxes = function() {
    };

    canvas.refreshTargetBoxes = function() {
        var htmlStr;
        var dist;
        // ------------- TARGET ORDER ------------ //
        var order = this.overpowerData.targetOrder;
        if (!order) {
            boxTargetOrder.textContent = "";
            boxConfirm.style.display = 'none';
        } else {
            var changedOrder = ((order.brandnew && order.size !== 0) || (!order.brandnew && order.originSize != order.size)) ;
            dist = grid.stepsBetween(order.sourcePl.loc, order.targetPl.loc);
            dist = Math.ceil(dist/10);
            if (order.brandnew) {
                htmlStr =  "<hr class='target'> Launch from "+order.sourcePl.name + " to "+ order.targetPl.name+" ("+dist+" turns)?<br> &nbsp;&nbsp;&nbsp;";
                if (changedOrder) {
                    htmlStr += "<b>Launch ship sized: "+order.size+"</b>";
                } else {
                    htmlStr += "Launch ship sized: "+order.size;
                }
            } else {
                htmlStr =  "<hr class='target'> Modifiy launch order from "+order.sourcePl.name + " to "+ order.targetPl.name+" (currently sized "+order.originSize+") ("+dist+" turns)?<br> &nbsp;&nbsp;&nbsp;";
                if (changedOrder) {
                    htmlStr += "<b>New ship size: "+order.size+"</b>";
                } else {
                    htmlStr += "New ship size: "+order.size;
                }
            }
            if (changedOrder) {
                boxConfirm.style.display = 'inline-block';
            } else {
                boxConfirm.style.display = 'none';
            }
            boxTargetOrder.innerHTML = htmlStr;
        }
        // ------------- TARGET TWO ------------ //
        var info = this.overpowerData.targetTwoInfo;
        var infoT1 = this.overpowerData.targetOneInfo;
        var hex = this.overpowerData.targetTwo;
        if (info && info.planet) {
            var targetOne = this.overpowerData.targetOne;
            if (targetOne && hex[0] === targetOne[0] && hex[1] == targetOne[1]) {
                htmlStr = "<b>Secondary Target:</b> Main";
            } else {
                htmlStr = "<b>Secondary Target:</b> ("+hex[0]+","+hex[1]+") <hr class='target'><b>"+info.planet.name+"</b>";
                if (targetOne) {
                    if (infoT1 && infoT1.planet) {
                        //htmlStr += "<br><b>[ Swap main/secondary targets ]</b>";
                    }
                    dist = grid.stepsBetween(targetOne, hex);
                    htmlStr += "<br>"+dist+" sectors from main target";
                }
            }
        } else {
            htmlStr = "<b>Secondary Target:</b> None";
        }
        boxTargetTwo.innerHTML  = htmlStr; 
        // ------------- TARGET ONE ------------ //
        data = this.overpowerData;
        hex = this.overpowerData.targetOne;
        if (!hex) {
            htmlStr = "<div class='center'> <b>Main Target:</b> None</div>";
        } else {
            htmlStr = "<div class='center'><b>Main Target:</b> ("+hex[0]+","+hex[1]+")</div>";
            if (!infoT1.planet && infoT1.ships.length < 1 && infoT1.trails.length < 1) {
                htmlStr += "<hr class=\"target\"><b>The Great Void</b>";
            } else {
                if (infoT1.planet) {
                    htmlStr += planetStr(infoT1.planet, infoT1, data);
                }
                if (infoT1.ships.length > 0) {
                    htmlStr += "<hr class='target'><b>Ship";
                    if (infoT1.ships.length > 1) {
                        htmlStr += "s";
                    }
                    htmlStr += " detected:</b><br><ul>";
                    infoT1.ships.forEach(function(ship) {
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
                }
                if (infoT1.trails.length > 0) {
                    htmlStr += "<hr class='target'><b>Ship passage";
                    if (infoT1.trails.length > 1) {
                        htmlStr += "s";
                    }
                    htmlStr += " detected:</b><br><ul>";
                    infoT1.trails.forEach(function(ship) {
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
                }
            }
        }
        boxTargetOne.innerHTML  = htmlStr; 
    };
   
    function planetStr(planet, info, data) {
        var str = "<hr class='target'><b>"+planet.name+"</b> &bull; ";
        if (!planet.turn) {
            str += "No planetary information available";
        } else {
            if (planet.controller === planet.fid) {
                str += "Your planet";
            } else {
                if (planet.controller) {
                    var name = data.fidMap.get(planet.controller).name;
                    str += name+" planet";
                } else {
                    str += "Neutral planet";
                }
                str += " &bull; Last seen on turn "+planet.turn;
            }
            str += "<br><b>Inhabitants:</b> "+planet.inhabitants+
                " &bull; <b>Resources:</b> "+planet.resources+
                " &bull; <b>Parts:</b> "+planet.parts;
            if (planet.avail) {
                str += " ("+planet.avail+" available)";
            }
        }
        if (info.orders && info.orders.length > 0 && (!data.targetOrder || info.orders.length > 1 || data.targetOrder !== info.orders[0])) {
            str += "<br><b>Launch orders:</b><ul>";
            for (var i = 0; i < info.orders.length; i++) {
                if (info.orders[i] === data.targetOrder) {
                    continue;
                }
                str += "<li><b>Target:</b> "+info.orders[i].targetPl.name+
                    " <b>Size:</b> "+info.orders[i].size+"</li>";
            }
            str += "</ul>";
        }
        return str;
    }

})();
