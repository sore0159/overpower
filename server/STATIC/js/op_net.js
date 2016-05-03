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

net.getFullView = function() {
    var url = "/overpower/json/fullviews/"+overpower.GID+"/"+overpower.FID;
    var callbacks = {};
    callbacks.success = successFV;
    ajax.getJSEND(url, callbacks);
};

function successFV(data) {
    overpower.data.parseFullView(data);
}

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


net.putMapView = function() {
    var jMV = { 
        gid: overpower.GID, 
        fid: overpower.FID, 
        center: [overpower.map.center.x, overpower.map.center.y],
        scale: overpower.map.getScale(),
        theta: overpower.map.screen.transform.theta,
        frame: [overpower.map.screen.canvas.width, overpower.map.screen.canvas.height ],
        stars: !overpower.stars.stopAnimation,
    };
    
    var url = "/overpower/json/mapviews";
    console.log("SENDING", JSON.stringify(jMV));

    var callbacks = {
        error: function(err, data) {
            console.log("Error syncing mapview data with server:", err, data);
        },
    };
    ajax.putJSEND(url, jMV, callbacks);
};




})(muleObj);


