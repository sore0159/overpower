(function() {

var canvas = document.getElementById('mainscreen');

canvas.gridShowing = function() {
    return this.muleGrid.scale >= 15;
};

canvas.drawMap = function() {
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
        if (scale > 34) {
            ctx.lineWidth = 2;
        }
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
    }
    if (scale >= 15) {
        ctx.font = "12pt monospace";
    } else if (scale >= 10) {
        ctx.font = "11pt monospace";
    } else if (scale >= 7.5) {
        ctx.font = "10pt monospace";
    } else if (scale >= 5) {
        ctx.font = "9pt monospace";
    } else { 
        ctx.font = "8pt monospace";
    }
    actualShips = [];

    var radiusX = scale *0.45;
    var radiusY = scale *0.25;
    function drawDot(hex) {
        if (visHexes && !visHexes.hasHex(hex)) {
            return;
        }
        var center = grid.h2Center(hex);
        var path = new Path2D();
        path.moveTo(center[0]+radiusX, center[1]);
        if (path.ellipse) {
            path.ellipse(center[0], center[1], radiusX, radiusY, 0, 0, 2*Math.PI);
        } else {
            path.rect(center[0]-radiusX, center[1]-radiusY, 2*radiusX, 2*radiusY);
        }
        ctx.fill(path);
        ctx.stroke(path);
    }
    function drawTrailDots(ship) {
        if (ship.loc.valid && (!visHexes || visHexes.hasHex(ship.loc.coord))) {
            actualShips.push(ship);
        }
        if (ship.controller === data.faction.fid)  {
            ctx.fillStyle = "#0fff0f";
        } else {
            ctx.fillStyle = "#ff0f0f";
        }
        ship.trail.forEach(drawDot);
    }
    function drawTrail(ship) {
        if (ship.loc.valid && (!visHexes || visHexes.hasHex(ship.loc.coord))) {
            actualShips.push(ship);
        }
        if (ship.trail.length < 2) {
            if (!ship.trail.length) {
                return;
            }
            drawTrailDots(ship);
            return;
        }
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
        ctx.strokeStyle = "#000000";
        ctx.setLineDash([]);
    }
    actualShips.forEach(function(ship) {
        if (ship.controller === data.faction.fid)  {
            ctx.fillStyle = "#0fff0f";
        } else {
            ctx.fillStyle = "#ff0f0f";
        }
        drawDot(ship.loc.coord);
    });

    var tarOne, tarTwo;
    var drawTars = false;
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
        var radius;
        if (scale > 10) {
            radius = scale * 0.45;
        } else if (scale > 4) {
            radius = 5;
        }  else if (scale > 2) {
            radius = 4;
        } else {
            radius = 3;
        }
        path.moveTo(center[0]+radius, center[1]);
        path.arc(center[0], center[1], radius, 0, 2*Math.PI);
        ctx.fill(path);
        ctx.stroke(path);
        if (scale < 1.5) {
            return;
        }
        ctx.strokeText(planet.name, center[0]+radius*0.25, center[1]-1.15*radius);
        ctx.fillText(planet.name, center[0]+radius*0.25, center[1]-1.15*radius);
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






canvas.drawMap();

})();
