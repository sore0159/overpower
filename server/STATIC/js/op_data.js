(function(muleObj) {

if (!muleObj.overpower) {
    muleObj.overpower = {};
}

if (!muleObj.overpower.data) {
    muleObj.overpower.data = {};
}

var overpower = muleObj.overpower;
var geometry = muleObj.geometry;
var data = muleObj.overpower.data;

data.factions = {};
data.targets = {};
data.power = {};
data.getName = function(fid) {
    if (fid === 0) {
        return "Hostile natives";
    }
    var fac = data.factions[fid];
    if (!fac || !fac.name ) {
        return "UNKNOWN FACTION";
    }
    return fac.name;
};

data.targets.setT1 = function(hex) {
    hex = data.targets.assist(hex);
    if (data.targets.primary && hex.eq(data.targets.primary.hex)) {
        return;
    }
    data.targets.dropOrder();
   
    data.targets.primary = { hex: hex };
    if (data.targets.secondary) {
        data.targets.secondary.copy = hex.eq(data.targets.secondary.hex);
        data.targets.secondary.dist = hex.stepsTo(data.targets.secondary.hex);
        data.targets.secondary.turns = Math.ceil(data.targets.secondary.dist / 10);
    }

    var planet = data.planetGrid.getHex(hex);
    if (!planet) {
        return;
    }
    data.targets.primary.planet = planet;
    if (!data.targets.secondary || !data.targets.secondary.planet) {
        return;
    }
    data.targets.makeOrder(planet, data.targets.secondary.planet);
};
data.targets.setT2 = function(hex, force) {
    hex = data.targets.assist(hex);
    var planet = data.planetGrid.getHex(hex);
    if (!planet) {
        if (force) {
            data.targets.dropOrder();
            delete data.targets.secondary;
        }
        return;
    }
    if (data.targets.secondary && hex.eq(data.targets.secondary.hex)) {
        return;
    }
    data.targets.dropOrder();
    data.targets.secondary = { hex: hex, planet: planet };
    if (data.targets.primary) {
        data.targets.secondary.copy = hex.eq(data.targets.primary.hex);
        data.targets.secondary.dist = hex.stepsTo(data.targets.primary.hex);
        data.targets.secondary.turns = Math.ceil(data.targets.secondary.dist / 10);
    } else {
        data.targets.secondary.copy = false;
    }
    if (!data.targets.primary || !data.targets.primary.planet) {
        return;
    }
    data.targets.makeOrder(data.targets.primary.planet, planet);
};
data.targets.dropOrder = function() {
    console.log("DROP ORDER");
    if (data.targets.order && data.targets.order.curO) {
        delete data.targets.order.curO.isTarget;
    }
    if (data.targets.order) {
        var planet = data.targets.order.sourcePL;
        planet.avail = planet.availBAK;
        delete planet.availBAK;
        planet = data.targets.order.targetPL;
        if (planet.landingBAK[1]) {
            planet.landing[planet.landingBAK[0]] = planet.landingBAK[1];
        } else {
            delete planet.landing[planet.landingBAK[0]];
        }
        delete planet.landingBAK;
    }
    delete data.targets.order;
};
data.targets.makeOrder = function(targetPL, sourcePL) {
    if (sourcePL.hex.eq(targetPL.hex)) {
        return;
    }
    var curO = sourcePL.sourceLaunchOrders.getHex(targetPL.hex);
    var size;
    if (curO) {
        size = curO.size;
        curO.isTarget = true;
    } else {
        size = 0;
    }
    if (!sourcePL.avail && !size) {
        return;
    }
    sourcePL.availBAK = sourcePL.avail;
    var dist = sourcePL.hex.stepsTo(targetPL.hex);
    var turns = Math.ceil(dist/10);
    targetPL.landingBAK = [turns, targetPL.landing[turns]];
    data.targets.order = { size: size,
        curO: curO,
        sourcePL: sourcePL,
        targetPL: targetPL,
        dist: dist,
        turns: turns,
    };
    console.log("SET ORDER", data.targets.order);
};
data.targets.modOrder = function(delta) {
    var dat = data.targets.order;
    if (!dat || !delta) {
        return false;
    }
    var newSize = dat.size + delta;
    if (delta + dat.size < 0) {
        delta = -1*dat.size;
    } else if (delta > dat.sourcePL.avail) {
        delta = dat.sourcePL.avail;
    }
    if (!delta) {
        return false;
    }
    dat.size += delta;
    dat.sourcePL.avail -= delta;
    dat.modified = dat.sourcePL.avail !== dat.sourcePL.availBAK;
    dat.targetPL.modLanding(dat.turns, delta);
    return true;
};

data.targets.assist = function(hex) {
    if (overpower.map.isShowGrid() || !data.planetGrid || data.planetGrid.getHex(hex)) {
        return hex;
    }
    var search;
    var scale = overpower.map.getScale();
    if (scale > 5) {
        search = 1;
    } else if (scale > 2.5) {
        search = 2;
    } else {
        search = 3;
    }
    for (var i = 1; i <= search; i += 1) {
        var list = hex.ring(i);
        for (var j = 0; j < list.length; j += 1) {
            if (data.planetGrid.getHex(list[j])) {
                return list[j];
            }
        }
    }
    return hex;
};

data.targets.isT1 = function(hex) {
    return (this.primary && hex.eq(this.primary.hex));
};
data.targets.isT2 = function(hex) {
    return (this.secondary && hex.eq(this.secondary.hex));
};
data.targets.setPowerOrder = function(planet, type) {
    if (!type || planet.myControl < 1 || planet.myPower === type) {
        return;
    }
    data.targets.power = {
        hex: planet.hex,
        planet: planet,
        type: type,
        changed: !data.power.hex.eq(planet.hex) || data.power.type !== type,
    };
};

data.parseFullView = function(fullView) {
    // GAME //
    data.game = fullView.game;
     // FACTIONS //
    data.factions = { list:[] };
    fullView.factions.forEach(function(fac) {
        data.factions.list.push(fac.fid);
        data.factions[fac.fid] = fac;
        if (fac.fid === overpower.FID) {
            fac.me = true;
            data.factions.myFaction = fac;
        }
    });
    var makeTr = function() {
        var tr = {};
        data.factions.list.forEach(function(fid) {
            tr[fid] = 0;
        });
        return tr;
    };

    // PLANETS //
    var setAvail = function() {
            var sum;
            var power = this.myPower;
            if (power === 1) {
                sum = this.antimatter;
            } else if (power === -1) {
                sum = this.tachyons;
            } else {
                this.avail = 0;
                return;
            }
            if (this.myPresence < sum) {
                sum = this.myPresence;
            }
            this.sourceLaunchOrders.forEach(function(ord) {
                sum -= ord.size;
            });
            if (sum > 0) {
                this.avail = sum;
            } else {
                this.avail = 0;
            }
            return;
    };
    var modLanding = function(turns, delta) {
        if (this.landing[turns]) {
            this.landing[turns] += delta;
        } else {
            this.landing[turns] = delta;
        }
        if (this.landing[turns] < 1) {
            delete this.landing[turns];
        }
    };
    var modTruce = function(fid, on) {
        if (fid === overpower.FID) {
            return;
        }
        var cur = this.truces[fid];
        if (!cur && cur !== 0) {
            return;
        }
        if (on) {
            if (cur === 0) {
                this.truces[fid] = 2;
            } else if (cur === -1) {
                this.truces[fid] = 1;
            }
        } else {
            if (cur === 1) {
                this.truces[fid] = -1;
            } else if (cur === 2) {
                this.truces[fid] = 0;
            }
        }
    };

       
    data.planetGrid = new geometry.HexMap();
    fullView.planetviews.forEach(function(pv) {
        pv.hex = new geometry.Hex(pv.loc[0], pv.loc[1]);
        data.planetGrid.setHex(pv.hex, pv);
        pv.setAvail = setAvail;
        pv.modLanding = modLanding;
        if (pv.primaryfaction === overpower.FID) {
            pv.myControl = 1;
            pv.myPower = pv.primarypower;
            pv.myPresence = pv.primarypresence;
        } else if (pv.secondaryfaction === overpower.FID) {
            pv.myControl = 2;
            pv.myPower = pv.secondarypower;
            pv.myPresence = pv.secondarypresence;
        } else {
            pv.myControl = 0;
            pv.myPower = 0;
            pv.myPresence = 0;
        }
        if (pv.myPower === 1) {
            pv.avail = pv.antimatter;
        } else if (pv.myPower == -1) {
            pv.avail = pv.tachyons;
        }
        if (pv.avail && pv.myPresence < pv.avail) {
            pv.avail = pv.myPresence;
        }
        pv.sourceLaunchOrders = new geometry.HexMap();
        pv.targetLaunchOrders = new geometry.HexMap();
        pv.truces = makeTr();
        pv.modTruce = modTruce;
        pv.landing = [];
    });
    // POWER ORDER //
    data.power = { 
        loc: fullView.powerorder.loc,
        hex: new geometry.Hex(fullView.powerorder.loc[0], fullView.powerorder.loc[1]),
        type: fullView.powerorder.uppower,
    };
    data.power.planet = data.planetGrid.getHex(data.power.hex);
    data.targets.setPowerOrder(data.power.planet, data.power.type);
    // LAUNCH ORDERS //
    fullView.launchorders.forEach(function(ord) {
        ord.sourceHex = new geometry.Hex(ord.source[0], ord.source[1]);
        ord.targetHex = new geometry.Hex(ord.target[0], ord.target[1]);
        ord.dist = ord.sourceHex.stepsTo(ord.targetHex);
        ord.turns = Math.ceil(ord.dist/10);
        var pv = data.planetGrid.getHex(ord.sourceHex);
        if (!pv) {
            console.log("ERROR: NO SOURCE PLANET FOUND FOR ORDER", ord);
            return;
        }
        pv.sourceLaunchOrders.setHex(ord.targetHex, ord);
        pv.avail -= ord.size;
        pv = data.planetGrid.getHex(ord.targetHex);
        if (!pv) {
            console.log("ERROR: NO TARGET PLANET FOUND FOR ORDER", ord);
            return;
        }
        pv.targetLaunchOrders.setHex(ord.sourceHex, ord);
        pv.modLanding(ord.turns, ord.size);
    });
    // TRUCES //
    fullView.truces.forEach(function(truce) {
        truce.hex = (new geometry.Hex()).addArray(truce.loc);
        var pl = data.planetGrid.getHex(truce.hex);
        if (!pl) {
            console.log("BAD TRUCE:", truce, "-- CAN'T FIND PLANET");
            return;
        }
        pl.truces[truce.trucee] = 1;
    });

    // MAP //
    data.mapView = fullView.mapview;
    overpower.map.snapTo(new geometry.Hex(fullView.mapview.center[0], fullView.mapview.center[1]));
    overpower.map.redraw = true;

};

data.launchOrderConfirmed = function(order) {
    var curOrder = order.sourcePL.sourceLaunchOrders.getHex(order.targetHex);
    var flag;
    if (curOrder) {
        var delta;
        if (curOrder.isTarget) {
            data.targets.dropOrder();
            flag = true;
        }
        if (order.size) {
            delta = order.size - curOrder.size;
            curOrder.size = order.size;
        } else {
            order.sourcePL.sourceLaunchOrders.deleteHex(order.targetHex);
            order.targetPL.targetLaunchOrders.deleteHex(order.sourceHex);
            delta = -curOrder.size;
        }
        order.sourcePL.avail -= delta;
        order.targetPL.modLanding(order.turns, delta);
        if (flag) {
            data.targets.makeOrder(order.targetPL, order.sourcePL);
        }
        return;
    }
    if (order.size < 1) {
        return;
    }
    var tOrd = data.targets.order;
    if (tOrd && tOrd.sourcePL.hex.eq(order.sourceHex) && tOrd.targetPL.hex.eq(order.targetHex)) {
        data.targets.dropOrder();
        flag = true;
        order.isTarget = true;
    }
    order.dist = order.sourceHex.stepsTo(order.targetHex);
    order.sourcePL.sourceLaunchOrders.setHex(order.targetHex, order);
    order.sourcePL.avail -= order.size;
    order.targetPL.targetLaunchOrders.setHex(order.sourceHex, order);
    order.targetPL.modLanding(order.turns, order.size);
    if (flag) {
        data.targets.makeOrder(order.targetPL, order.sourcePL);
    }
};

data.powerOrderConfirmed = function(order) {
    data.power = {
        type: order.uppower,
        hex: (new geometry.Hex()).addArray(order.loc),
    };
    data.power.planet = data.planetGrid.getHex(data.power.hex);
    if (data.targets.power.hex.eq(data.power.hex)) {
        data.targets.changed = !data.power.hex.eq(data.targets.power.planet.hex) || data.power.type !== data.targets.type;
    }
};

data.trucesConfirmed = function(truces) {
    var tr = {};
    truces.trucees.forEach(function(fid) {
        tr[fid] = true;
    });
    data.factions.list.forEach(function(fid) {
        var cur = truces.planet.truces[fid];
        if (tr[fid] && cur !== -1) {
            truces.planet.truces[fid] = 1;
        } else if (!tr[fid] && cur !== 2) {
            truces.planet.truces[fid] = 0;
        }
    });
};

data.turnBufferConfirmed = function(buff) {
    data.factions.myFaction.donebuffer = buff;
};

})(muleObj);
