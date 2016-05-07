(function(muleObj) {

if (!muleObj.ajax) {
    console.log("OP AJAX FAILED: REQUIRES MULEOBJ AJAX LIB");
    return;
}
if (!muleObj.overpower) {
    muleObj.overpower = {};
}

// callbacks can contain:
//         .success .error .netError .serverError .fail
//ajax.getJSEND = function(url, callbacks) {
//ajax.putJSEND = function(url, obj, callbacks) {
var ajax = muleObj.ajax;
var overpower = muleObj.overpower;
if (!overpower.net) {
    overpower.net = {};
}
var net = overpower.net;
net.timers = {};


net.pingTurn = function(repeat) {
    var url = "/overpower/json/games/"+overpower.GID;
    var callbacks = {
        error: function(err, data) {
            console.log("Error checking turn data with server:", err, data);
            if (repeat) {
                window.setTimeout(net.pingTurn, 300000, true);
            }
        },
        success: function(gameDat) {
            if (gameDat[0] && gameDat[0].turn > overpower.data.game.turn) {
                overpower.data.game.newTurn = true;
                overpower.html.icon.animate();
            } else if (repeat) {
                var setTime = (document.hidden === false) ? 60000: 300000;
                window.setTimeout(net.pingTurn, setTime, true);
            }
        },
    };
    ajax.getJSEND(url, callbacks);
};

window.setTimeout(net.pingTurn, 60000, true);



net.getFullView = function() {
    var url = "/overpower/json/fullviews/"+overpower.GID+"/"+overpower.FID;
    var callbacks = {};
    callbacks.success = successFV;
    ajax.getJSEND(url, callbacks);
};

function successFV(data) {
    overpower.data.parseFullView(data);
    overpower.map.redraw = true;
    overpower.html.infobox.render();
}

net.putTurnBuffer = function(buff) {
    var curTime = Date.now();
    if (net.timers.turnBuff && curTime - net.timers.turnBuff < 1000) {
        console.log("CANCELLING TURN BUFFER UPDATE: TOO SOON, EXECUTUS (1 second)");
        return;
    }
    net.timers.turnBuff = curTime;
    var jF = { gid: overpower.GID,
        fid: overpower.FID,
        donebuffer: buff,
    };
    var url = "/overpower/json/factions";
 
    var callbacks = {
        error: function(err, data) {
            console.log("Error syncing donebuffer data with server:", err, data);
        },
        success: function(jDat) {
            overpower.data.turnBufferConfirmed(buff);
            overpower.html.infobox.overview.renderTurn();
            if (buff) {
                net.pingTurn(true);
            }
        },
    };
    console.log("SENDING", JSON.stringify(jF));
    ajax.putJSEND(url, jF, callbacks);
};

net.putPowerOrder = function(order) {
    var curTime = Date.now();
    if (net.timers.powerOrder && curTime - net.timers.powerOrder < 1000) {
        console.log("CANCELLING POWER ORDER UPDATE: TOO SOON, EXECUTUS (1 second)");
        return;
    }
    net.timers.powerOrder = curTime;
    var jPO = { gid: overpower.GID,
        fid: overpower.FID,
        loc: [order.hex.x, order.hex.y],
        uppower: order.type,
    };
    var url = "/overpower/json/powerorders";
 
    var callbacks = {
        error: function(err, data) {
            console.log("Error syncing powerorder data with server:", err, data);
        },
        success: function(jDat) {
            overpower.data.powerOrderConfirmed(jPO);
            overpower.html.infobox.overview.renderPowerOrder();
            overpower.map.redraw = true;
        },
    };
    ajax.putJSEND(url, jPO, callbacks);
};

net.putLaunchOrder = function(order) {
    var curTime = Date.now();
    if (net.timers.launchOrder && curTime - net.timers.launchOrder < 1000) {
        console.log("CANCELLING LAUNCH ORDER UPDATE: TOO SOON, EXECUTUS (1 second)");
        return;
    }
    net.timers.launchOrder = curTime;
    var sourceHex = order.sourcePL.hex;
    var targetHex = order.targetPL.hex;
    var targetPL = order.targetPL;
    var sourcePL = order.sourcePL;
    var turns = order.turns;
    var jLO = { gid: overpower.GID,
        fid: overpower.FID,
        source: [sourceHex.x, sourceHex.y],
        target: [targetHex.x, targetHex.y],
        size: order.size,
    };
    var url = "/overpower/json/launchorders";
 
    var callbacks = {
        error: function(err, data) {
            console.log("Error syncing launchorder data with server:", err, data);
        },
        success: function(jDat) {
            jLO.sourceHex = sourceHex;
            jLO.targetHex = targetHex;
            jLO.sourcePL = sourcePL;
            jLO.targetPL = targetPL;
            jLO.turns = turns;
            overpower.data.launchOrderConfirmed(jLO);
            overpower.map.redraw = true;
        },
    };
    ajax.putJSEND(url, jLO, callbacks);
};


net.putTruces = function(planet) {
    var curTime = Date.now();
    if (net.timers.truces && curTime - net.timers.truces < 1000) {
        console.log("CANCELLING TRUCES UPDATE: TOO SOON, EXECUTUS (1 second)");
        return;
    }
    net.timers.truces = curTime;
    var jTR = { gid: overpower.GID,
        fid: overpower.FID,
        loc: [planet.hex.x, planet.hex.y],
        trucees: [],
    };
    var trList = Object.keys(planet.truces);
    trList.forEach(function(key) {
        if (planet.truces[key]) {
            jTR.trucees.push(parseInt(key));
        }
    });
    var url = "/overpower/json/truces";
 
    var callbacks = {
        error: function(err, data) {
            console.log("Error syncing truce data with server:", err, data);
        },
        success: function(jDat) {
            jTR.planet = planet;
            overpower.data.trucesConfirmed(jTR);
            if ( overpower.data.targets.isT1(planet.hex) ||overpower.data.targets.isT2(planet.hex) ) {
                console.log("REDRAW TARGET BOX");
            }
        },
    };
    ajax.putJSEND(url, jTR, callbacks);
};




})(muleObj);


