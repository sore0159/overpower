(function() {

var canvas = document.getElementById('mainscreen');

canvas.parseOPData = function() {
    var data = this.overpowerData;
    var grid = this.muleGrid;
    var fidMap = new Map();
    data.factions.forEach(function(faction) {
        fidMap.set(faction.fid, faction);
    });
    data.fidMap = fidMap;


    var gridMap = new Map();
    gridMap.getHex = function(hex) {
        var yMap = this.get(hex[0]);
        if (!yMap) {
            return null;
        }
        return yMap.get(hex[1]);
    };
    data.gridMap = gridMap;
  
    data.planetviews.forEach(function(planet) {
        var yMap = gridMap.get(planet.loc[0]);
        if (!yMap) {
            yMap = new Map();
            gridMap.set(planet.loc[0], yMap);
        }
        yMap.set(planet.loc[1], planet);
        //
        if (planet.controller === data.faction.fid) {
            planet.avail = planet.parts;
        }
        data.shipviews.forEach(function(ship) {
            if (ship.dest.valid && grid.ptsEq(ship.dest.coord, planet.loc)) {
                ship.dest.planet = planet;
            }
        });
    });
 

    data.orders.forEach(function(order) {
        order.sourcePl = gridMap.getHex(order.source);
        order.sourcePl.avail -= order.size;
        order.targetPl = gridMap.getHex(order.target);
    });
    data.map = {"center":data.mapview.center};
    canvas.centerHex(canvas.overpowerData.map.center);

    data.swapTargets = function() {
        var t1 = this.targetOne;
        var t2 = this.targetTwo;
        if (this.targetOneInfo && this.targetOneInfo.planet) {
            this.targetTwo = t1;
            this.targetTwoInfo = {planet: this.targetOneInfo.planet};
        } else {
            this.targetTwo = null;
            this.targetTwoInfo = null;
        }
        this.setTargetOne(t2);
    };

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
        if (hex) {
            var closestPl, dist;
            this.planetviews.forEach(function(planet) {
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
            this.shipviews.forEach(function(ship) {
                if (ship.loc.valid) {
                    if (targetPt && grid.ptsEq(targetPt, ship.loc.coord)) {
                        targetInfo.ships.push(ship);
                    } else {
                        var steps = grid.stepsBetween(ship.loc.coord, hex);
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
                this.shipviews.forEach(function(ship) {
                    if (ship.loc.valid) {
                        if (targetPt && grid.ptsEq(targetPt, ship.loc.coord)) {
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
                    var loc = targetInfo.planet.loc;
                    if (grid.ptsEq(order.source, targetInfo.planet.loc)) {
                        targetInfo.orders.push(order);
                    }
                });
            }
        }
        this.targetOne = targetPt;
        this.targetOneInfo = targetInfo;
        if (targetInfo.planet && this.targetTwoInfo && this.targetTwoInfo.planet) {
            this.setTargetOrder(targetInfo.planet, this.targetTwoInfo.planet);
        } else {
            this.setTargetOrder();
        }
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
        this.planetviews.forEach(function(planet) {
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
        } else if (drop) {
            this.targetTwo = null;
            this.targetTwoInfo = null;
        }
        if (this.targetOneInfo && this.targetOneInfo.planet && this.targetTwoInfo && this.targetTwoInfo.planet) {
            this.setTargetOrder(this.targetOneInfo.planet, this.targetTwoInfo.planet);
        } else {
            this.setTargetOrder();
        }
    };

    data.setTargetOrder = function(pl1, pl2) {
        if (this.targetOrder) {
            if (this.targetOrder.sourcePl === pl1 && this.targetOrder.targetPl === pl2) {
                return;
            }
            if (this.targetOrder.originSize || this.targetOrder.originSize === 0) {
                var diff = this.targetOrder.originSize - this.targetOrder.size;
                this.targetOrder.size = this.targetOrder.originSize;
                this.targetOrder.originSize = null;
                this.targetOrder.sourcePl.avail -= diff;
            }
        }
        if (!pl1 || !pl2 || pl1.name === pl2.name) {
            this.targetOrder = null;
        } else {
            var order = null;
            for (var i = 0; i< this.orders.length; i++) {
                var o = this.orders[i];
                if (grid.ptsEq(o.source, pl1.loc) && grid.ptsEq(o.target, pl2.loc)) {
                    order = o;
                    if (!o.originSize) {
                        o.originSize = o.size;
                    }
                    break;
                }
            }
            if (!order && pl1.avail) {
                order = {"brandnew": true, "gid": this.game.gid, "fid": this.faction.fid, "source":pl1.loc, "target": pl2.loc, "size":0, "targetPl": pl2, "sourcePl": pl1};
            }
            this.targetOrder = order;
        }
    };

};


})();

