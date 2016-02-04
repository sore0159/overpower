(function() {

var canvas = document.getElementById('mainscreen');
var data = canvas.overpowerData;

data.setTargetOne = function(hex, help) {
    var grid = canvas.muleGrid;
    var targetPt;
    var targetInfo = {ships:[], trails:[]};
    var tolerance;
    if (grid.scale > 2.5) {
        tolerance = 2;
    } else if (grid.scale > 1.75) {
        tolerance = 3;
    } else {
        tolerance = 4;
    }
    if (!help) {
        targetPt = hex;
    }
    var closestPl, dist;
    data.planetviews.forEach(function(planet) {
        if (dist === 0) {
            return;
        }
        var steps = grid.stepsBetween(planet.loc, hex);
        if (steps === 0) {
            targetPt = hex;
            targetInfo.planet = planet;
        } else if (!dist || steps < dist) {
            closestPl = planet;
            dist = steps;
        }
    });

    var closestShLoc;
    data.shipviews.forEach(function(ship) {
        if (ship.loc) {
            if (targetPt && grid.ptsEq(targetPt, ship.loc)) {
                targetInfo.ships.push(ship);
            } else {
                var steps = grid.stepsBetween(ship.loc, hex);
                if (!dist || steps < dist) {
                    closestPl = null;
                    closestShLoc = hex;
                    dist = steps;
                }
            }
        }
        if (targetPt) {
            ship.trail.forEach(function(pt) {
                if (grid.ptsEq(pt, targetPt)) {
                    targetInfo.trails.push(ship);
                }
            });
        }
    });
    if (!targetPt) {
        if (closestPl && dist < tolerance) {
            targetPt = closestPl.loc;
            targetInfo.planet = closestPl;
        } else if (closestShLoc && dist < tolerance) {
            targetPt = closestShLoc;
        } else {
            targetPt = hex;
        }
        data.shipviews.forEach(function(ship) {
            if (ship.loc) {
                if (targetPt && grid.ptsEq(targetPt, ship.loc)) {
                    targetInfo.ships.push(ship);
                } 
            }
            ship.trail.forEach(function(pt) {
                if (grid.ptsEq(pt, targetPt)) {
                    targetInfo.trails.push(ship);
                }
            });
        });
    }
    if (targetInfo.planet && this.orders) {
        targetInfo.orders = [];
        this.orders.forEach(function(order) {
            if (order.source === targetInfo.planet.pid) {
                targetInfo.orders.push(order);
            }
        });
    }
    this.targetOne = targetPt;
    this.targetOneInfo = targetInfo;
};

data.setTargetTwo = function(hex, help, drop) {
    var grid = canvas.muleGrid;
    var tolerance;
    if (grid.scale > 2.5) {
        tolerance = 2;
    } else if (grid.scale > 1.75) {
        tolerance = 3;
    } else {
        tolerance = 4;
    }
    var closestPl, dist;
    data.planetviews.forEach(function(planet) {
        if (dist === 0) {
            return;
        }
        var steps = grid.stepsBetween(planet.loc, hex);
        if (!dist || steps < dist) {
            closestPl = planet;
            dist = steps;
        }
    });
    if (dist === 0 || help && dist < tolerance) {
        this.targetTwo = closestPl.loc;
        this.targetTwoInfo = {planet: closestPl};
        return;
    }
    if (drop) {
        this.targetTwo = null;
        this.targetTwoInfo = null;
        return;
    }
};


canvas.parseOPData = function() {
    var data = this.overpowerData;
    console.log("DATA:", data);
};





canvas.parseOPData();


})();

