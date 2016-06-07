(function(muleObj) {

var html = muleObj.html;
var geometry = muleObj.geometry;
var overpower = muleObj.overpower;

if (!overpower.map) {
    overpower.map = {};
}

var map = overpower.map;

var canvas = document.getElementById('mainscreen');
map.screen = new html.ScreenTransform(canvas, 0,0, 0.25, 14.5, 0.55);
map.orderscreen = { canvas: document.getElementById('orderscreen'), clear: true };
map.hexGrid = new geometry.HexGrid(map.screen.transform);
map.getScale = function() {
    return this.screen.transform.scale;
};

map.setFrame = function(width, height) {
    width = (width < 120) ? 120: width;
    height = (height < 90) ? 90: height;
    var spacer = document.getElementById('canvasspacer');
    spacer.style.width = width+'px';
    var infobox = document.getElementById('infobox');
    infobox.style.marginLeft = ''+(width + 20)+'px';
    document.body.style.minWidth = ""+(width+540)+'px';
    map.screen.resize(width, height);
    var stars = overpower.stars;
    stars.screen.resize(width, height);
    stars.moreStars();
    map.orderscreen.canvas.width = width;
    map.orderscreen.canvas.height = height;
};
map.scaleStep = function(zoomIn) {
    var curScale = map.getScale();
    if ((curScale >= 40 && zoomIn) || (curScale <= 0.25 && !zoomIn)) {
        return false;
    }
    var dir = (zoomIn) ? 1 : -1;
    var delta;
    if (curScale > 20 || (zoomIn && curScale === 20)) {
        delta = 2*dir;
    } else if (curScale > 15 || (zoomIn && curScale == 15)) {
        delta = dir;
    } else if (curScale > 5 || (zoomIn && curScale == 5)) {
        delta = 0.5*dir;
    } else {
        delta = 0.25*dir;
    }
    var dScale;
    if (curScale + delta <= 0.25) {
        dScale = 0.25 - curScale;
    } else if (curScale+delta >= 40) {
        dScale = 40 - curScale;
    } else {
        dScale = delta;
    }
    map.screen.setScale(dScale);
    return true;
};
map.snapTo = function(hex) {
    map.center = hex;
    map.screen.setCenter(hex.centerPt());
};
map.distToCenter = function() {
    if (!map.center) {
        return 0;
    }
    var centerAt = map.hexGrid.centerPt(map.center);
    return centerAt.dist(map.screen.center());
};
map.moveTowardCenter = function(distLeft) {
    var endpt;
    if (distLeft <= 10) {
        endpt = map.center.centerPt();
    } else {
        var speed;
        if (distLeft <= 20) {
            speed = 10;
        } else {
            //speed = 10 + Math.pow((distLeft-20), 3);
            speed = 10 + Math.pow((distLeft - 20)/25, 1.25);
            if (speed > distLeft / 2) {
                speed = distLeft /2;
            }
        }
        var startPt = map.screen.center();
        var towardPt = map.hexGrid.centerPt(map.center);
        var polarDat = startPt.polarTo(towardPt);
        endpt = startPt.addPolar(speed, polarDat[1]);
        endpt = map.screen.out2in(endpt);
    }
    map.screen.setCenter(endpt);
};

function mapClick(inPoint, button, shift, ctrl) {
    var hex = inPoint.hexAt();
    if (button === 0 && ctrl && shift) { 
        overpower.commands.confirmLaunchOrder();
        return;
    }
    overpower.commands.clickHex(hex, button, shift);
}
map.screen.setInClick(mapClick);

function mapWheel(up, shift, ctrl) {
    var delta;
    if (shift && ctrl) {
        overpower.commands.modMapFrame(up > 0);
        return;
    }
    if (ctrl) {
        overpower.commands.modLaunchOrder(up > 0);
        return;
    }
    if (shift) {
        overpower.commands.rotateMap(up > 0);
        return;
    }
    overpower.commands.zoomMap(up > 0);
}

map.screen.setWheel(mapWheel);

map.hexPath = function(hex, path) {
    if (!path) {
        path = new Path2D();
    }
    var corners = map.hexGrid.cornerPts(hex);
    path.moveTo(corners[5].x, corners[5].y);
    corners.forEach(function(pt) {
        path.lineTo(pt.x, pt.y);
    });
    return path;
};
map.visibleHexes = function() {
    //return map.hexGrid.hexesIn(map.screen.canvas.width * 0.25, map.screen.canvas.height * 0.25 , map.screen.canvas.width * 0.75, map.screen.canvas.height * 0.75);
    return map.hexGrid.hexesIn(0,0, map.screen.canvas.width, map.screen.canvas.height);
};

map.renderGrid = function(ctx, list) {
    map.setLineWidth(ctx);
    path = new Path2D();
    var drawHex = function(hex) {
        var corners = map.hexGrid.cornerPts(hex);
        path.moveTo(corners[0].x, corners[0].y);
        path.lineTo(corners[1].x, corners[1].y);
        path.lineTo(corners[2].x, corners[2].y);
        path.lineTo(corners[3].x, corners[3].y);
    };
    list.forEach(drawHex);
    ctx.strokeStyle = "#2f2f8f";
    ctx.stroke(path);
};

map.renderTargets = function(ctx, visHexes, targets) {
    ctx.lineWidth += 1;
    if (targets.secondary && (!targets.primary || !targets.primary.hex.eq(targets.secondary.hex)) && visHexes.hasHex(targets.secondary.hex)) {
        ctx.strokeStyle = "#999900";
        ctx.stroke(map.hexPath(targets.secondary.hex));
    }
    if (targets.primary && visHexes.hasHex(targets.primary.hex)) {
        ctx.strokeStyle = "#ffff00";
        ctx.stroke(map.hexPath(targets.primary.hex));
    }
    ctx.lineWidth -= 1;
};

map.renderPlanet = function(ctx, pv, planetRad, showGrid) {
    var scale = map.getScale();
    if (!showGrid && overpower.data.targets.isT1(pv.hex)) {
        ctx.fillStyle = "#ffff00";
    } else if (!showGrid && overpower.data.targets.isT2(pv.hex)) {
        ctx.fillStyle = "#999900";
    } else if (pv.primaryfaction === overpower.FID) {
        if (pv.secondaryfaction) {
            ctx.fillStyle = "#af00af";  // contested
        } else {
            ctx.fillStyle = "#0fff0f";
        }
    } else if (pv.secondaryfaction === overpower.FID) {
        ctx.fillStyle = "#af00af";  // contested
    } else if (pv.primaryfaction) {
        ctx.fillStyle = "#ff0f0f";
    //} else if (pv.hex.x === 0) {
    } else {
        ctx.fillStyle = "#ffffff";
    }
    var path = new Path2D();
    var center = map.hexGrid.centerPt(pv.hex);
    center.y += map.yBoost();
    path.moveTo(center.x+planetRad, center.y);
    path.arc(center.x, center.y, planetRad, 0, 2*Math.PI);
    ctx.fill(path);
    ctx.stroke(path);
    if (scale < 1.5) {
        return;
    }
    // TOP //
    var nameStr = "";
    if (pv.avail) {
        nameStr = "("+pv.avail+")";
    }
    if (scale >= 3.75) {
        nameStr += pv.name;
    }
    if (nameStr) {
        ctx.fillText(nameStr, center.x+planetRad*0.25, center.y-1.15*planetRad);
    }
    // BOTTOM //
};

map.setLineWidth = function(ctx) {
   var scale = map.getScale();
    if (scale < 25) {
        ctx.lineWidth = 0.5;
    } else if (scale < 35) {
        ctx.lineWidth = 1;
    } else {
        ctx.lineWidth = 1.5;
    }
};
map.isShowGrid = function() {
    return map.getScale() >= 15;
};
map.yBoost = function() {
    if (map.screen.transform.scale < 15) {
        return 0;
    }
    return map.screen.transform.scale * -0.25;
};
map.getPlanetRad = function() {
    var scale = map.screen.transform.scale;
    if (scale > 10) {
        return scale * 0.45;
    } else if (scale > 4) {
        return 5;
    }  else if (scale > 2) {
        return 4;
    } else {
        return 3;
    }
};

map.setFont = function(ctx) {
    var scale = map.screen.transform.scale;
    if (scale >= 20) {
        ctx.font = "12pt Courier New";
        //ctx.font = "12pt Courier New";
        return 14;
    } else if (scale >= 15) {
        ctx.font = "11pt Courier New";
        //ctx.font = "11pt Courier New";
        return 13;
    } else if (scale >= 10) {
        ctx.font = "10pt Courier New";
        //ctx.font = "10pt Courier New";
        return 12;
    } else if (scale >= 7.5) {
        ctx.font = "9pt Courier New";
        //ctx.font = "9pt Courier New";
        return 11;
    } else { 
        ctx.font = "8pt Courier New";
        //ctx.font = "8pt Courier New";
        return 10;
    }
};

map.renderLaunchOrders = function(ctx, pv, visHexes) {
    var boost = map.yBoost();
    var f = function(ord) {
        if (ord.isTarget) {
            return;
        }
        var startPt = map.hexGrid.centerPt(ord.sourceHex);
        var endPt = map.hexGrid.centerPt(ord.targetHex);
        var path = new Path2D();
        path.moveTo(startPt.x, startPt.y + boost);
        path.lineTo(endPt.x, endPt.y + boost);
        ctx.stroke(path);
    };
    pv.sourceLaunchOrders.forEach(f);
    if (!visHexes) {
        return;
    }
    pv.targetLaunchOrders.forEach(function(ord) {
        if (visHexes.hasHex(ord.sourceHex)) {
            return;
        }
        f(ord);
    });
};

//map.render = function() {
function renderMap(timestamp) {
    var data = overpower.data;
    var canvas = map.screen.canvas;
    var ctx = canvas.getContext('2d');
    // CLEAR //
    ctx.clearRect(0,0,canvas.width, canvas.height);
    //
    var scale = map.getScale();
    var showGrid = map.isShowGrid();
    var visHexes;
    if (scale > 5) {
        visHexes = map.visibleHexes();
    }
    var path;
    // GRID //
    if (showGrid) {
        map.renderGrid(ctx, visHexes.list);
        if (data.targets) {
            map.renderTargets(ctx, visHexes, data.targets);
        }
    }
    var planetRad = map.getPlanetRad();
    var fontHeight = map.setFont(ctx);
    // ORDERS/TRAILS/DESTINATIONS  //
    // SHIPS //
    // PLANETS //
    if (data.planetGrid) {
        map.setLineWidth(ctx);
        var targetPLs = {};
        plList = [];
        ctx.lineWidth += 0.5;
        ctx.strokeStyle = "#0FAFAF";
        data.planetGrid.forEach(function(pv, hex) {
            if (!visHexes || visHexes.hasHex(hex)) {
                map.renderLaunchOrders(ctx, pv, visHexes);
                if (data.targets.isT1(hex)) {
                    targetPLs.t1 = pv;
                } else if (data.targets.isT2(hex)) {
                    targetPLs.t2 = pv;
                } else {
                    plList.push(pv);
                }
            }
        });
        ctx.lineWidth -= 0.5;
        ctx.strokeStyle = "#000000";
        plList.forEach(function(pv) {
            map.renderPlanet(ctx, pv, planetRad, showGrid);
        });
        if (targetPLs.t2) {
            map.renderPlanet(ctx, targetPLs.t2, planetRad, showGrid);
        }
        if (targetPLs.t1) {
            map.renderPlanet(ctx, targetPLs.t1, planetRad, showGrid);
        }
    }
}

function renderTargetOrder(timestamp) {
    var canvas = map.orderscreen.canvas;
    var ctx = canvas.getContext('2d');
    ctx.clearRect(0,0,canvas.width, canvas.height);
    map.orderscreen.clear = true;
    if (!overpower.data.targets.order) {
        return;
    }
    var startHex = overpower.data.targets.order.sourcePL.hex;
    var endHex = overpower.data.targets.order.targetPL.hex;
    map.orderscreen.clear = false;
    map.setLineWidth(ctx);
    ctx.lineWidth += 1;
    ctx.strokeStyle = "#0FAFAF";
    var scale = map.getScale();
    if (scale > 5) {
        ctx.setLineDash([8,8]);
        ctx.lineDashOffset = -(timestamp*0.01)%16;
    } else {
        ctx.setLineDash([3,3]);
        ctx.lineDashOffset = -(timestamp*0.01)%6;
    }
    var startPt = map.hexGrid.centerPt(startHex);
    var endPt = map.hexGrid.centerPt(endHex);
    var boost = map.yBoost();
    startPt.y += boost;
    endPt.y += boost;
    var path = new Path2D();
    path.moveTo(startPt.x, startPt.y);
    path.lineTo(endPt.x, endPt.y);
    ctx.stroke(path);
    var pRad = map.getPlanetRad() ;
    var numStr = overpower.data.targets.order.size; 
    var numWidth = ctx.measureText(numStr).width;
    ctx.font = "14pt Courier New";
    var textPt = endPt;
    var fHeight = 15;
    ctx.setLineDash([]);
    ctx.lineWidth = 1;
    ctx.fillStyle = "#000000";
    ctx.fillRect(textPt.x-(1.25*pRad+numWidth+1), textPt.y-fHeight-1, numWidth+2, fHeight+5);
    ctx.strokeRect(textPt.x-(1.25*pRad+numWidth+1), textPt.y-fHeight-1, numWidth+2, fHeight+5);
    ctx.fillStyle = "#0FAFAF";
    if (!overpower.data.targets.order.modified || ((timestamp*0.01)%8 > 3)) {
        ctx.fillText(numStr, textPt.x-(1.25*pRad+numWidth), textPt.y);
    }
    if (scale < 15) {
        return;
    }
    path = new Path2D();
    var cycle = 0.25* Math.sin(timestamp * 0.005);
    var rad = pRad * 0.45;
    path.moveTo(startPt.x+rad, startPt.y);
    path.arc(startPt.x, startPt.y, rad, 0, 2*Math.PI);
    ctx.fill(path);
    path.moveTo(endPt.x+rad, endPt.y);
    rad = pRad * 0.45 * (1 + cycle) ;
    path.arc(endPt.x, endPt.y, rad, 0, 2*Math.PI);
    ctx.fill(path);
}

map.animate = function(timestamp) {
    var redraw = map.redraw;
    if (redraw !== false) {
        redraw = true;
    }
    map.redraw = false;
    var dist = map.distToCenter();
    if (dist > 1) {
        map.moveTowardCenter(dist);
        redraw = true;
    }
    if (redraw) {
        renderMap(timestamp);
    }
    if (!map.orderscreen.clear || overpower.data.targets.order) {
        renderTargetOrder(timestamp);
    }
    window.requestAnimationFrame(map.animate);
};

map.animate();

})(muleObj);
