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
    var hexList;
    if (scale > 5) {
        hexList = grid.visibleHexList(0,0,this.width,this.height);
    }
    var ctx = this.getContext('2d');
    ctx.clearRect(0,0,this.width, this.height);
    var path;
    var gridShowing = this.gridShowing();
    if (gridShowing) {
        this.drawGrid(hexList);
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
    var tarOne, tarTwo;
    var drawTars = false;
    function drawPlanet(planet) {
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
        ctx.strokeStyle = "#000000";
        var path = new Path2D();
        var center = grid.h2Center(planet.loc);
        var radius = scale *0.45;
        if (scale >= 15) {
            center[1] = center[1]-scale*0.25;
        }
        if (scale < 11 && scale > 4) {
            radius = 5;
        }  else if (scale < 5) {
            radius = 3;
        }
        path.moveTo(center[0]+radius, center[1]);
        path.arc(center[0], center[1], radius, 0, 2*Math.PI);
        ctx.fill(path);
        ctx.stroke(path);
        if (scale < 1.5) {
            return;
        }
        if (scale >= 15) {
            ctx.font = "12pt mono";
        } else if (scale >= 10) {
            ctx.font = "11pt mono";
        } else if (scale >= 7.5) {
            ctx.font = "10pt mono";
        } else if (scale >= 5) {
            ctx.font = "9pt mono";
        } else { 
            ctx.font = "8pt mono";
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
