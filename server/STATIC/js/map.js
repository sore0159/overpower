(function() {

var canvas = document.getElementById('mainscreen');

canvas.gridShowing = function() {
    return this.muleGrid.scale >= 15;
};

canvas.drawMap = function(timestamp) {
    var data = this.overpowerData;
    if (!data) {
        console.log("NO OVERPOWER DATA TO DRAW MAP");
        return;
    }
    var grid = this.muleGrid;
    var scale = grid.scale;
    var visHexes;
    if (scale > 5) {
        visHexes = grid.visibleHexes(0,0,this.width,this.height);
    }
    var ctx = this.getContext('2d');
    ctx.clearRect(0,0,this.width, this.height);
    var path;
    var gridShowing = this.gridShowing();
    if (gridShowing) {
        this.drawGrid(visHexes.list);
        ctx.lineWidth += 1;
        if (data.targetTwo && data.targetOne != data.targetTwo) {
            path = this.hexPath(data.targetTwo);
            ctx.strokeStyle = "#999900";
            ctx.stroke(path);
        }
        if (data.targetOne) {
            path = this.hexPath(data.targetOne);
            ctx.strokeStyle = "#ffff00";
            ctx.stroke(path);
        }
        ctx.lineWidth -= 1;
    }
    var fontHeight;
    if (scale >= 15) {
        ctx.font = "12pt monospace";
        fontHeight = 14;
    } else if (scale >= 10) {
        ctx.font = "11pt monospace";
        fontHeight = 13;
    } else if (scale >= 7.5) {
        ctx.font = "10pt monospace";
        fontHeight = 12;
    } else if (scale >= 5) {
        ctx.font = "9pt monospace";
        fontHeight = 11;
    } else { 
        ctx.font = "8pt monospace";
        fontHeight = 10;
    }

    var planetRad;
    if (scale > 10) {
        planetRad = scale * 0.45;
    } else if (scale > 4) {
        planetRad = 5;
    }  else if (scale > 2) {
        planetRad = 4;
    } else {
        planetRad = 3;
    }
    actualShips = new Map();
    destShips = [];
    function shipCheck(ship) {
        if (ship.loc.valid && (!visHexes || visHexes.hasHex(ship.loc.coord))) {
            var hex = ship.loc.coord;
            var shipList = actualShips.get(hex);
            if (shipList) {
                shipList.push(ship);
            } else {
                actualShips.set(hex, [ship]);
            }
            if (ship.dest.valid) {
                destShips.push(ship);
            }
        } else if (ship.dest.valid && ship.loc.valid && (!visHexes ||visHexes.hasHex(ship.dest.coord))) {
            destShips.push(ship);
        }
    }
 
    function drawDot(hex) {
        if (visHexes && !visHexes.hasHex(hex)) {
            return;
        }
        var center = grid.h2Center(hex);
        var path = new Path2D();
        path.moveTo(center[0]+planetRad, center[1]);
        if (path.ellipse) {
            path.ellipse(center[0], center[1], planetRad, planetRad*0.55, 0, 0, 2*Math.PI);
        } else {
            path.rect(center[0]-planetRad, center[1]-planetRad*0.55, 2*planetRad, 1.1*planetRad);
        }
        ctx.fill(path);
        ctx.stroke(path);
    }
    function drawTrailDots(ship) {
        shipCheck(ship);
        if (ship.controller === data.faction.fid)  {
            ctx.fillStyle = "#0fff0f";
        } else {
            ctx.fillStyle = "#ff0f0f";
        }
        ship.trail.forEach(drawDot);
    }
    function drawTrail(ship) {
        if (!ship.trail.length) {
            shipCheck(ship);
            return;
        }
        if (ship.trail.length < 2 && !ship.loc.valid) {
            ctx.strokeStyle = "#000000";
            drawTrailDots(ship);
            return;
        }
        shipCheck(ship);
        if (ship.controller === data.faction.fid)  {
            ctx.strokeStyle = "#0fff0f";
        } else {
            ctx.strokeStyle = "#ff0f0f";
        }
        var startHex = ship.trail[0], endHex;
        if (ship.loc.valid) {
            endHex = ship.loc.coord;
        } else {
            endHex = ship.trail[ship.trail.length-1];
        }
        if (visHexes && !visHexes.hasHex(startHex) && !visHexes.hasHex(endHex)) {
            return;
        }
        var path = new Path2D();
        var start = grid.h2Center(startHex), end = grid.h2Center(endHex);
        path.moveTo(start[0], start[1]);
        path.lineTo(end[0], end[1]);
        ctx.stroke(path);
    }
  
    if (gridShowing) {
        ctx.strokeStyle = "#000000";
        ctx.globalAlpha = 0.25;
        data.shipviews.forEach(drawTrailDots);
        ctx.globalAlpha = 1;
    } else {
        ctx.lineWidth = 1;
        if (scale > 5) {
            ctx.setLineDash([8,8]);
        } else if (scale > 2.5) {
            ctx.setLineDash([3,3]);
        }
        data.shipviews.forEach(drawTrail);
        ctx.setLineDash([]);
    }
    ctx.strokeStyle = "#0FAFAF";
    destShips.forEach(function(ship) {
        var path = new Path2D();
        var start = grid.h2Center(ship.dest.coord);
        var end = grid.h2Center(ship.loc.coord);
        if (gridShowing) {
            start[1] = start[1]-scale*0.25;
        }
        path.moveTo(start[0], start[1]);
        path.lineTo(end[0], end[1]);
        ctx.stroke(path);
    });
    ctx.strokeStyle = "#000000";
    actualShips.forEach(function(shipList, hex) {
        var shipNames = [];
        var shipSizes = [];
        var ship;
        var enemy = false;
        for (var i = 0; i < shipList.length; i++) {
            ship = shipList[i];
            var fac = data.fidMap.get(ship.controller);
            shipNames.push(data.fidMap.get(ship.controller).name);
            //shipNames.push(data.fidMap.get(ship.controller).name+"("+ship.size+")");
            shipSizes.push(ship.size);
            if (ship.controller !== data.faction.fid)  {
                enemy = true;
            }
        }
        if (enemy) {
            ctx.fillStyle = "#ff0f0f";
        } else {
            ctx.fillStyle = "#0fff0f";
        }
        drawDot(hex);
        if (scale < 2.5) {
            return;
        }
        center = grid.h2Center(hex);
        var numStr = shipSizes.join("|");
        var numStrWidth = ctx.measureText(numStr).width;
        ctx.strokeText(numStr, center[0]-numStrWidth/2, center[1]+0.60*planetRad+fontHeight);
        ctx.fillText(numStr, center[0]-numStrWidth/2, center[1]+0.60*planetRad+fontHeight);
        if (scale < 5) {
            return;
        }
        var nameStr = shipNames.join();
        ctx.strokeText(nameStr, center[0]+planetRad*0.25, center[1]-0.55*planetRad);
        ctx.fillText(nameStr, center[0]+planetRad*0.25, center[1]-0.55*planetRad);
    });
    var drawTarOrder = false;
    function drawOrder(order) {
        if (order === data.targetOrder && !drawTarOrder) {
                return;
        }
        if (visHexes && !visHexes.hasHex(order.sourcePl.loc) && !visHexes.hasHex(order.targetPl.loc)) {
            return;
        }
        var path = new Path2D();
        var start = grid.h2Center(order.sourcePl.loc);
        var end = grid.h2Center(order.targetPl.loc);
        if (gridShowing) {
        start[1] = start[1]-scale*0.25;
        end[1] = end[1]-scale*0.25;
        }
        path.moveTo(start[0], start[1]);
        path.lineTo(end[0], end[1]);
        ctx.stroke(path);
    }
    if (data.orders) {
        ctx.lineWidth += 0.5;
        ctx.strokeStyle = "#0FAFAF";
        data.orders.forEach(drawOrder);
        ctx.lineWidth -= 0.5;
    }
    if (data.targetOrder) {
        drawTarOrder = true;
        ctx.strokeStyle = "#0FAFAF";
        ctx.lineWidth += 1;
        if (scale > 5) {
            ctx.setLineDash([8,8]);
            ctx.lineDashOffset = -(timestamp*0.01)%16;
        } else {
            ctx.setLineDash([3,3]);
            ctx.lineDashOffset = -(timestamp*0.01)%6;
        }

        drawOrder(data.targetOrder);

        ctx.lineDashOffset = 0;
        ctx.lineWidth -= 1;
        ctx.setLineDash([]);
    }

    var tarOne, tarTwo;
    var drawTars = false;
    ctx.strokeStyle = "#000000";
    function drawPlanet(planet) {
        if (visHexes && !visHexes.hasHex(planet.loc)) {
            return;
        }
        if (!gridShowing && grid.ptsEq(planet.loc, data.targetOne)) {
            if (drawTars) {
                ctx.fillStyle = "#ffff00";
            } else {
                tarOne = planet;
                return;
            }
        } else if (!gridShowing && grid.ptsEq(planet.loc, data.targetTwo)) {
            if (drawTars) {
                ctx.fillStyle = "#999900";
            } else {
                tarTwo = planet;
                return;
            }
        } else if (planet.controller === data.faction.fid)  {
            ctx.fillStyle = "#0fff0f";
        } else if (planet.controller) {
            ctx.fillStyle = "#ff0f0f";
        } else {
            ctx.fillStyle = "#ffffff";
        }
        var path = new Path2D();
        var center = grid.h2Center(planet.loc);
        if (gridShowing) {
            center[1] = center[1]-scale*0.25;
        }
        path.moveTo(center[0]+planetRad, center[1]);
        path.arc(center[0], center[1], planetRad, 0, 2*Math.PI);
        ctx.fill(path);
        ctx.stroke(path);
        if (scale < 1.5) {
            return;
        }
        var nameStr;
        if (scale < 3.75) {
            if (planet.avail) {
                nameStr = "("+planet.avail+")";
                ctx.strokeText(nameStr, center[0]+planetRad*0.25, center[1]-1.15*planetRad);
                ctx.fillText(nameStr, center[0]+planetRad*0.25, center[1]-1.15*planetRad);
            }

        } else {
            if (planet.avail) {
                nameStr = "("+planet.avail+")"+planet.name;
            } else {
                nameStr = planet.name;
            }
            ctx.strokeText(nameStr, center[0]+planetRad*0.25, center[1]-1.15*planetRad);
            ctx.fillText(nameStr, center[0]+planetRad*0.25, center[1]-1.15*planetRad);
        }
        if (scale >= 2.5) {
            if (planet.turn !== 0) {
                var numStr;
                if (planet.controller === data.faction.fid) {
                    numStr = planet.presence+"|"+planet.resources+"|"+planet.parts;
                } else {
                    numStr = planet.presence+"?|"+planet.resources+"?";
                }
                var numStrWidth = ctx.measureText(numStr).width;
                ctx.strokeText(numStr, center[0]-numStrWidth/2, center[1]+planetRad+fontHeight);
                ctx.fillText(numStr, center[0]-numStrWidth/2, center[1]+planetRad+fontHeight);
            }
      
        }

    }
    data.planetviews.forEach(drawPlanet);
    drawTars = true;
    if (tarTwo) {
        drawPlanet(tarTwo);
    }
    if (tarOne) {
        drawPlanet(tarOne);
    }
};


function animateMap(timestamp) {
    if (canvas.overpowerData) {
        if (canvas.animating) {
            canvas.animating = false;
            canvas.drawMap(timestamp);
        }
        if (canvas.overpowerData.targetOrder) {
            canvas.animating = true;
        }
        var curLoc = canvas.muleGrid.h2Center(canvas.overpowerData.map.center);
        if (!canvas.mapCentered()) {
            canvas.moveTowardCenter(curLoc);
            canvas.animating = true;
        }
    }
    window.requestAnimationFrame(animateMap);
}

canvas.setCenterDest = function(hex) {
    this.aniInfoCenter = 0;
    this.overpowerData.map.center = hex;
};

canvas.moveTowardCenter = function(pt) {
    if (!this.aniInfoCenter) {
        this.aniInfoCenter = 0;
    }
    this.aniInfoCenter += 1;
    var shiftX = this.width/2 - pt[0];
    var shiftY = this.height/2 - pt[1];
    var scale = this.aniInfoCenter/15;
    this.muleGrid.shift(shiftX*scale, shiftY*scale);
};


canvas.mapCentered = function() {
    var curLoc = this.muleGrid.h2Center(this.overpowerData.map.center);
    return !(Math.abs(curLoc[0] - this.width/2) > 1 || Math.abs(curLoc[1] - this.height/2) > 1); 
};


window.requestAnimationFrame(animateMap);

})();
